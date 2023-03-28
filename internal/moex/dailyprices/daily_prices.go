package dailyprices

import (
	"encoding/xml"
	"io"
	"moexobserver/internal/models"
	"strconv"
	"time"
)

type documentDailyPriceList struct {
	Rows []rowDailyPrice `xml:"data>rows>row"`
}

type rowDailyPrice struct {
	Tradedate string `xml:"TRADEDATE,attr"`
	Open      string `xml:"OPEN,attr"`
	High      string `xml:"HIGH,attr"`
	Low       string `xml:"LOW,attr"`
	Close     string `xml:"CLOSE,attr"`
	Volume    string `xml:"VOLUME,attr"`
}

func GetDailyPrices(in io.Reader, tickerId int32) ([]models.DailyPrice, error) {
	var key documentDailyPriceList
	decodeXML := xml.NewDecoder(in)
	err := decodeXML.Decode(&key)
	if err != nil {
		return nil, err
	}
	var dailyPrices []models.DailyPrice
	for _, row := range key.Rows {

		tradeDate, err := time.Parse("2006-01-02", row.Tradedate)
		if err != nil {
			return nil, err
		}
		var vol int64
		if len(row.Volume) == 0 {
			vol = 0
		} else {
			vol, err = strconv.ParseInt(row.Volume, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		open, high, low, close, err := rowToOHLC(row)
		if err != nil {
			return nil, err
		}

		price := models.DailyPrice{
			TickerId:  tickerId,
			Tradedate: tradeDate,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    vol,
		}

		if open == 0 && close == 0 {
			continue
		}

		dailyPrices = append(dailyPrices, price)
	}

	return dailyPrices, nil
}

func rowToOHLC(row rowDailyPrice) (open, high, low, close float64, err error) {
	if len(row.Open) == 0 {
		open = 0
	} else {
		open, err = strconv.ParseFloat(row.Open, 64)
		if err != nil {
			return open, high, low, close, err
		}
	}

	if len(row.High) == 0 {
		high = 0
	} else {
		high, err = strconv.ParseFloat(row.High, 64)
		if err != nil {
			return open, high, low, close, err
		}
	}

	if len(row.Low) == 0 {
		low = 0
	} else {
		low, err = strconv.ParseFloat(row.Low, 64)
		if err != nil {
			return open, high, low, close, err
		}
	}

	if len(row.Close) == 0 {
		close = 0
	} else {
		close, err = strconv.ParseFloat(row.Close, 64)
		if err != nil {
			return open, high, low, close, err
		}
	}

	return open, high, low, close, err
}
