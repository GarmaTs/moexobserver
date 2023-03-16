package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"moexobserver/internal/config"
	"moexobserver/internal/moex/moexreader"
	"moexobserver/internal/moex/tickers"
	"time"
)

const TICKERS_URL = "http://iss.moex.com/iss/history/engines/stock/markets/shares/boards/TQBR/securities.xml?iss.meta=off&history.columns=BOARDID,TRADEDATE,SHORTNAME,SECID,NUMTRADES,VALUE,VOLUME"

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

func main() {
	cfg, err := config.ReadConf("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	return

	chanTickers := make(chan []tickers.TickerName)
	timerForTickers := time.NewTicker(1 * time.Second)

	for {
		select {
		// Получение всех тикеров реквестом и запись в канал
		case <-timerForTickers.C:
			go WriteAllTickersToChan(TICKERS_URL, 0, chanTickers)
		// Чтение из канала тикеров
		case tickers, ok := <-chanTickers:
			if ok {
				go ProceedAllTickers(tickers)
			}
		}
	}

	// tickers, err := GetAllTickers(TICKERS_URL, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, row := range tickers {
	// 	fmt.Println(row.Secid, row.Boardid, row.Tradedate, row.Volume)
	// }

	return

	//	XML FROM FILE
	// filename := "tickers.xml"
	// reader, err := moexreader.GetXMLTickersFromFile(filename)
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

func ProceedAllTickers(tickers []tickers.TickerName) {
	// for _, row := range tickers {
	// 	fmt.Println(row.Secid, row.Boardid, row.Tradedate, row.Volume)
	// }
	fmt.Println("len(tickers)=", len(tickers))
	fmt.Println(time.Now())
}

func WriteAllTickersToChan(srcUrl string, start int, c chan<- []tickers.TickerName) {
	var allTickers []tickers.TickerName

	bRun := true
	for bRun {
		url := fmt.Sprintf("%s&start=%d", srcUrl, start)
		start += 100

		reader, err := moexreader.GetXMLTickerByRequest(url)
		if err != nil {
			return
		}

		tickerSlice, err := tickers.GetTickerList(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		allTickers = append(allTickers, tickerSlice...)

		if len(tickerSlice) == 0 || start > 10000 {
			bRun = false
		}
	}
	c <- allTickers
}

func GetAllTickers(srcUrl string, start int) ([]tickers.TickerName, error) {
	var allTickers []tickers.TickerName

	bRun := true
	for bRun {
		url := fmt.Sprintf("%s&start=%d", srcUrl, start)
		start += 100

		reader, err := moexreader.GetXMLTickerByRequest(url)
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
