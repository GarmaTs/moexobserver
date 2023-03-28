create table public.tickers (
	id bigserial primary key,
	secid     text,
	shortname text,
	boardid   text,
	tradedate timestamp,
	volume    bigint,
	unique (secid, boardid)
);