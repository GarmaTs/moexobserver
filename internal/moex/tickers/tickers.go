package tickers

import (
	"encoding/xml"
	"io"
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

type TickerName struct {
	Secid     string
	Shortname string
	Boardid   string
	Tradedate time.Time
	Volume    int64
}

func GetTickerList(in io.Reader) (error, []TickerName) {
	var key documentTickerList
	decodeXML := xml.NewDecoder(in)
	err := decodeXML.Decode(&key)
	if err != nil {
		return err, nil
	}

	var tickers []TickerName
	for _, row := range key.Rows {
		if len(row.Secid) == 0 {
			continue
		}
		tradeDate, err := time.Parse("2006-01-02", row.Tradedate)
		if err != nil {
			return err, nil
		}
		vol, err := strconv.ParseInt(row.Volume, 10, 64)
		if err != nil {
			return err, nil
		}

		ticker := TickerName{
			Secid:     row.Secid,
			Shortname: row.Shortname,
			Boardid:   row.Boardid,
			Tradedate: tradeDate,
			Volume:    vol,
		}
		tickers = append(tickers, ticker)
	}

	return nil, tickers
}
