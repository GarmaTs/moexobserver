package models

import (
	"time"
)

type Ticker struct {
	Id        int64
	Secid     string
	Shortname string
	Boardid   string
	Tradedate time.Time
	Volume    int64
}

type DailyPrice struct {
	TickerId  int32
	Tradedate time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}
