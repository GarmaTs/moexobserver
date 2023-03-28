package tickers

import (
	"encoding/xml"
	"io"
	"moexobserver/internal/models"
	"strconv"
	"time"
)

type documentTickerList struct {
	Rows []rowTickerList `xml:"data>rows>row"`
}

type rowTickerList struct {
	Boardid   string `xml:"BOARDID,attr"`
	Tradedate string `xml:"TRADEDATE,attr"`
	Shortname string `xml:"SHORTNAME,attr"`
	Secid     string `xml:"SECID,attr"`
	Volume    string `xml:"VOLUME,attr"`
}

func GetTickerList(in io.Reader) ([]models.Ticker, error) {
	var key documentTickerList
	decodeXML := xml.NewDecoder(in)
	err := decodeXML.Decode(&key)
	if err != nil {
		return nil, err
	}

	var tickers []models.Ticker
	for _, row := range key.Rows {
		if len(row.Secid) == 0 {
			continue
		}
		tradeDate, err := time.Parse("2006-01-02", row.Tradedate)
		if err != nil {
			return nil, err
		}
		vol, err := strconv.ParseInt(row.Volume, 10, 64)
		if err != nil {
			return nil, err
		}

		ticker := models.Ticker{
			Secid:     row.Secid,
			Shortname: row.Shortname,
			Boardid:   row.Boardid,
			Tradedate: tradeDate,
			Volume:    vol,
		}
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}
