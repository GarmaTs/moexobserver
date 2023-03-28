package data

import (
	"database/sql"
	"fmt"
	"log"
	"moexobserver/internal/models"
)

type Store struct {
	DB      *sql.DB
	Tickers interface {
		Insert([]models.Ticker)
	}
	DailyPrices interface {
		GetLastTradeDates() []models.Ticker
		Insert(dailyPrices []models.DailyPrice)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{DB: db,
		Tickers:     tickersModel{DB: db},
		DailyPrices: dailyPricesModel{DB: db},
	}
}

type tickersModel struct {
	DB *sql.DB
}

func (m tickersModel) Insert(tickers []models.Ticker) {
	//truncate table public.tickers RESTART IDENTITY CASCADE;
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

	_, err := m.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("tickers updated")
}

type dailyPricesModel struct {
	DB *sql.DB
}

func (m dailyPricesModel) GetLastTradeDates() []models.Ticker {
	var tickers []models.Ticker

	query := `
select
	t.id as ticker_id, t.secid, t.boardid, COALESCE(p.tradedate, '1899-12-31') as tradedate
from (
	select max(tradedate) as tradedate, p.ticker_id
	from public.daily_prices as p
	group by p.ticker_id
) as p
right join public.tickers t on t.id = p.ticker_id
--where t.id in (4)
order by ticker_id;`

	rows, err := m.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var t models.Ticker
	for rows.Next() {
		err := rows.Scan(&t.Id, &t.Secid, &t.Boardid, &t.Tradedate)
		if err != nil {
			log.Fatal(err)
		}
		tickers = append(tickers, t)
	}

	if len(tickers) == 0 {
		log.Println("GetLastTradeDates: no tickers")
	}
	return tickers
}

func (m dailyPricesModel) Insert(prices []models.DailyPrice) {
	if prices == nil {
		return
	}

	var subQuery string
	var i int
	for i = 0; i < len(prices); i++ {
		subQuery += fmt.Sprintf("(%d, '%s', %.3f, %.3f, %.3f, %.3f, %d),",
			prices[i].TickerId, prices[i].Tradedate.Format("20060102"),
			prices[i].Open, prices[i].High, prices[i].Low, prices[i].Close, prices[i].Volume)

		if len(prices) == 1 || (i > 0 && i%1000 == 0) || i == len(prices)-1 {
			subQuery = subQuery[:len(subQuery)-1]

			query := `
insert into public.daily_prices (ticker_id, tradedate, open, high, low, close, volume)
values ` + subQuery + `
on conflict (ticker_id, tradedate) do update
set
open = excluded.open,
high = excluded.high,
low = excluded.low,
close = excluded.close,
volume = excluded.volume;`

			_, err := m.DB.Exec(query)
			if err != nil {
				fmt.Println(query)
				log.Fatal(err)
			}

			subQuery = ""
		}
	}

	fmt.Println("Inserted ticker_id:", prices[0].TickerId, "len(prices):", len(prices))

	// ДОБАВИТЬ ТАБЛИЦУ С ДАТАМИ ПОСЛЕДНИХ ПРОВЕРОК
}
