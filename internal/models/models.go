package models

import (
	"time"
)

type Ticker struct {
	Secid     string
	Shortname string
	Boardid   string
	Tradedate time.Time
	Volume    int64
}
