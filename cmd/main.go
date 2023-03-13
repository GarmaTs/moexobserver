package main

import (
	"fmt"
	"moexobserver/internal/moex/moexreader"
	"moexobserver/internal/moex/tickers"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a filename!")
		return
	}

	filename := arguments[1]
	reader, err := moexreader.GetXMLTickersFromFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer reader.Close()

	err, tickers := tickers.GetTickerList(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, row := range tickers {
		fmt.Println(row.Secid, row.Boardid, row.Tradedate, row.Volume)
	}

}
