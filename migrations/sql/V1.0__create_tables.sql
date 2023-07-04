-- Creating persongdb table
create table persondb (
	id uuid,
	salary INTEGER,
	married BOOLEAN,
	profession VARCHAR(30),
	primary key (id)
);

-- Creating users table
create table users (
	id uuid,
	username VARCHAR(30),
	password VARCHAR,
	refreshToken VARCHAR,
	primary key (id)
);
