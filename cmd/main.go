package main

import (
	"fmt"
	"os"
	"proj_moex/internal/tickers"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a filename!")
		return
	}

	filename := arguments[1]
	reader, err := os.Open(filename)
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
