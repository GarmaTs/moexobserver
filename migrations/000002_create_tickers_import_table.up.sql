create table public.tickers_import (
	secid     text,
	shortname text,
	boardid   text,
	Tradedate timestamp,
	Volume    bigint,
	primary key (secid, boardid)
);