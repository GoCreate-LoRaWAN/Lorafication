package db

// schema is the constant that contains the postgres database schema for
// the lorafication daemon.
const schema = `
CREATE TABLE IF NOT EXISTS nodes(
	id serial PRIMARY KEY,
	name varchar(255) NOT NULL,
	description text,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entities(
	id serial PRIMARY KEY,
	name varchar(255) NOT NULL,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notification_contracts(
	id serial PRIMARY KEY,
	node_id integer NOT NULL,
	entity_id integer NOT NULL,
	email varchar(255),
	sms integer,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW(),
	FOREIGN KEY(node_id) REFERENCES nodes(id),
	FOREIGN KEY(entity_id) REFERENCES entities(id)
);

ALTER TABLE notification_contracts
	ADD CONSTRAINT verify_notification_channels CHECK (num_nonnulls(email, sms) = 1);
`
