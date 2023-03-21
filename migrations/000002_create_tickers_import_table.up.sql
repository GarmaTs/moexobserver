create table public.tickers_import (
	secid     text,
	shortname text,
	boardid   text,
	tradedate timestamp,
	volume    bigint,
	primary key (secid, boardid)
);