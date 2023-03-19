create table public.tickers_import (
	id bigserial primary key,
	secid     text,
	shortname text,
	boardid   text,
	Tradedate timestamp,
	Volume    bigint,
	unique (secid, boardid)
);