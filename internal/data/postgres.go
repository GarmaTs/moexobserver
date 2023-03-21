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
	var tmpStr, subQuery string
	for i, row := range tickers {
		if i == len(tickers)-1 {
			tmpStr = fmt.Sprintf("('%s', '%s', '%s', '%s', %d)\n",
				row.Secid, row.Shortname, row.Boardid, row.Tradedate.Format("20060102"), row.Volume)
		} else {
			tmpStr = fmt.Sprintf("('%s', '%s', '%s', '%s', %d),\n",
				row.Secid, row.Shortname, row.Boardid, row.Tradedate.Format("20060102"), row.Volume)
		}
		subQuery += tmpStr
	}

	query := `
delete from public.tickers_import;
insert into public.tickers_import (secid, shortname, boardid, tradedate, volume)
values ` + subQuery + `;

update public.tickers a
	set shortname = i.shortname,
		tradedate = i.tradedate,
		volume = i.volume
from public.tickers t	
inner join public.tickers_import i on i.secid = t.secid and i.boardid = t.boardid
where
	a.id = t.id;
	
insert into public.tickers
	(secid, shortname, boardid, tradedate, volume)
select 
	i.secid, i.shortname, i.boardid, i.tradedate, i.volume
from public.tickers_import as i
where
	not exists (
		select 1 from public.tickers as t
		where t.secid = i.secid and t.boardid = i.boardid
	);

delete from public.tickers_import;`

	m.DB.Exec(query)

	fmt.Println("tickers updated")
}
