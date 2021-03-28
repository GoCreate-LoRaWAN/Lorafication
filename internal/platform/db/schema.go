package db

// schema is the constant that contains the postgres database schema for
// the lorafication daemon.
const schema = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS node(
	public_key UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	secret UUID NOT NULL DEFAULT uuid_generate_v4(),
	name varchar(255) NOT NULL,
	description text,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entity(
	id serial PRIMARY KEY,
	name varchar(255) NOT NULL,
	email varchar(255),
	sms integer,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW(),
	CONSTRAINT notify_channels_check CHECK (email IS NOT NULL OR sms IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS contract(
	id serial PRIMARY KEY,
	node_public_key UUID NOT NULL,
	entity_id integer NOT NULL,
	created timestamp NOT NULL DEFAULT NOW(),
	modified timestamp NOT NULL DEFAULT NOW(),
	FOREIGN KEY(node_public_key) REFERENCES node(public_key),
	FOREIGN KEY(entity_id) REFERENCES entity(id)
);`
