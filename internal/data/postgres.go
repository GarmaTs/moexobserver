package data

import (
	"database/sql"
	"fmt"
	"moexobserver/internal/models"
)

type store struct {
	DB      *sql.DB
	Tickers interface {
		Insert([]models.Ticker)
	}
}

func NewStore(db *sql.DB) store {
	return store{DB: db,
		Tickers: tickersModel{DB: db},
	}
}

type tickersModel struct {
	DB *sql.DB
}

func (m tickersModel) Insert(tickers []models.Ticker) {
	query := `delete from public.tickers_import;
		insert into public.tickers_import (secid, shortname, boardid, tradedate, volume)
		values `
	var tmpStr string
	for i, row := range tickers {
		//fmt.Println(row.Secid, row.Shortname, row.Boardid, row.Tradedate, row.Volume)
		if i == len(tickers)-1 {
			tmpStr = fmt.Sprintf("('%s', '%s', '%s', '%s', %d)\n", row.Secid, row.Shortname, row.Boardid, row.Tradedate.Format("20060102"), row.Volume)
		} else {
			tmpStr = fmt.Sprintf("('%s', '%s', '%s', '%s', %d),\n", row.Secid, row.Shortname, row.Boardid, row.Tradedate.Format("20060102"), row.Volume)
		}
		query += tmpStr
	}

	m.DB.Exec(query)

	fmt.Println(len(tickers))
}
