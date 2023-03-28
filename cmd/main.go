package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"moexobserver/internal/config"
	"moexobserver/internal/data"
	"moexobserver/internal/models"
	"moexobserver/internal/moex/dailyprices"
	"moexobserver/internal/moex/moexreader"
	"moexobserver/internal/moex/tickers"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const TICKERS_URL = "http://iss.moex.com/iss/history/engines/stock/markets/shares/boards/TQBR/securities.xml?iss.meta=off&history.columns=BOARDID,TRADEDATE,SHORTNAME,SECID,NUMTRADES,VALUE,VOLUME"

// Получение ohlc для одного тикера
func BatchFillDailyPricesChan(c chan<- []models.DailyPrice, ticker models.Ticker) {
	var dailyPrices []models.DailyPrice
	bRun := true
	start := 0
	for bRun {
		url := fmt.Sprintf("%s%s%s%s%s%d",
			"http://iss.moex.com/iss/history/engines/stock/markets/shares/boards/",
			ticker.Boardid, "/securities/", ticker.Secid,
			".xml?iss.meta=off&history.columns=TRADEDATE,OPEN,HIGH,LOW,CLOSE,VOLUME&start=",
			start)
		//fmt.Println("from BatchFillDailyPricesChan:\n", url)

		reader, err := moexreader.GetXMLByRequest(url)
		if err != nil {
			fmt.Println("err:", err)
		}
		tmpPrices, err := dailyprices.GetDailyPrices(reader, int32(ticker.Id))
		if err != nil {
			fmt.Println("err:", err)
		}

		if len(tmpPrices) == 0 || start > 20000 {
			bRun = false
		} else {
			dailyPrices = append(dailyPrices, tmpPrices...)
		}

		start += 100
	}
	c <- dailyPrices
}

func main() {
	minDate, _ := time.Parse("2006-01-02", "1900-01-01")
	cfg, err := config.ReadConf("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(db)

	// Канал и таймер для обновления списка тикеров
	chanTickers := make(chan []models.Ticker)
	timerForTickers := time.NewTicker(2 * time.Second)

	// Канал и таймер для получения последней даты каждого тикера
	chanLastTradeDates := make(chan models.Ticker)
	timerForLastTradeDates := time.NewTicker(5 * time.Second)

	// Канал для ohlc
	chanUpdateDailyPrices := make(chan []models.DailyPrice, 20)
	store := data.NewStore(db)

	for {
		select {
		// Получение всех тикеров реквестом и запись в канал
		case <-timerForTickers.C:
			go WriteAllTickersToChan(TICKERS_URL, 0, chanTickers)
		// Чтение из канала тикеров
		case tickers, ok := <-chanTickers:
			if ok {
				timerForTickers = time.NewTicker(100 * time.Second)
				go store.Tickers.Insert(tickers)
			}
		// Запись последних дат для всех тикеров в канал
		case <-timerForLastTradeDates.C:
			go GetLastTradeDates(chanLastTradeDates, &store)
		// Получение последней даты тикера из канала
		case tickerLastDate, ok := <-chanLastTradeDates:
			if ok {
				if tickerLastDate.Tradedate.Before(minDate) {
					// РЕАЛИЗОВАТЬ полную вставку
					go BatchFillDailyPricesChan(chanUpdateDailyPrices, tickerLastDate)
					//fmt.Println("Need full insert for", tickerLastDate.Id)
				} else {
					// Получение ohlc тикера реквестом и запись в канал
					//go FillDailyPricesChan(chanUpdateDailyPrices, tickerLastDate)
					fmt.Println("Need update for", tickerLastDate.Id)
				}
			}
		// Получение ohlc из канала и запись в базу
		case dailyPrices, ok := <-chanUpdateDailyPrices:
			if ok {
				go store.DailyPrices.Insert(dailyPrices)
			}
		}
	}

	//	XML FROM FILE
	// filename := "tickers.xml"
	// reader, err := moexreader.GetXMLFromFile(filename)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer reader.Close()

	// err, tickers := tickers.GetTickerList(reader)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// for _, row := range tickers {
	// 	fmt.Println(row.Secid, row.Boardid, row.Tradedate, row.Volume)
	// }
}

func GetLastTradeDates(c chan<- models.Ticker, store *data.Store) {
	lastTradeDates := store.DailyPrices.GetLastTradeDates()
	for _, lastDate := range lastTradeDates {
		c <- lastDate
	}
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}

func ProceedAllTickers(tickers []models.Ticker) {
	// for _, row := range tickers {
	// 	fmt.Println(row.Secid, row.Boardid, row.Tradedate, row.Volume)
	// }
	fmt.Println("len(tickers)=", len(tickers))
	fmt.Println(time.Now())

	//data.Store.Insert(tickers)
}

func WriteAllTickersToChan(srcUrl string, start int, c chan<- []models.Ticker) {
	var allTickers []models.Ticker

	bRun := true
	for bRun {
		url := fmt.Sprintf("%s&start=%d", srcUrl, start)
		reader, err := moexreader.GetXMLByRequest(url)
		if err != nil {
			return
		}

		tickerSlice, err := tickers.GetTickerList(reader)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(tickerSlice) == 0 || start > 10000 {
			bRun = false
		} else {
			allTickers = append(allTickers, tickerSlice...)
		}

		start += 100
	}
	c <- allTickers
}

func GetAllTickers(srcUrl string, start int) ([]models.Ticker, error) {
	var allTickers []models.Ticker

	bRun := true
	for bRun {
		url := fmt.Sprintf("%s&start=%d", srcUrl, start)
		start += 100

		reader, err := moexreader.GetXMLByRequest(url)
		if err != nil {
			return nil, err
		}

		tickerSlice, err := tickers.GetTickerList(reader)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		allTickers = append(allTickers, tickerSlice...)

		if len(tickerSlice) == 0 || start > 10000 {
			bRun = false
		}
	}

	return allTickers, nil
}
