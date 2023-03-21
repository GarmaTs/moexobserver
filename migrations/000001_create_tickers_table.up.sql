create table public.tickers (	
	secid     text,
	shortname text,
	boardid   text,
	Tradedate timestamp,
	Volume    bigint,
	primary key (secid, boardid)
);