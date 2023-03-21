create table public.daily_prices (
	ticker_id int not null,
	tradedate date not null,
	open decimal(18,3) not null,
	high decimal(18,3) not null,
	low decimal(18,3) not null,
	close  decimal(18,3) not null,
	volume bigint not null,
	primary key(ticker_id, tradedate)
);