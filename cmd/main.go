package main

import (
	"fmt"
	"moexobserver/internal/moex/moexreader"
	"moexobserver/internal/moex/tickers"
)

const TICKERS_URL = "http://iss.moex.com/iss/history/engines/stock/markets/shares/boards/TQBR/securities.xml?iss.meta=off&history.columns=BOARDID,TRADEDATE,SHORTNAME,SECID,NUMTRADES,VALUE,VOLUME"

func main() {
	tickers, err := GetAllTickers(TICKERS_URL, 0)
	if err != nil {
		fmt.Println(err)
	}
	for _, row := range tickers {
		fmt.Println(row.Secid, row.Boardid, row.Tradedate, row.Volume)
	}

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
