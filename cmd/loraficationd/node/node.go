// Package node interfaces between the node table in the database and the
// lorafication daemon.
package node

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// TODO: The secret should probably be stored using a hashing algorithm, similar to how a password is stored.

// Node is a struct representing the structure of a row in the node table
// of the database.
type Node struct {
	PublicKey   string    `db:"public_key"` // Primary key (it's a UUID).
	Secret      string    `db:"secret"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Created     time.Time `db:"created"`
	Modified    time.Time `db:"modified"`
}

// AuthenticateNode takes the key and secret of a node and finds the corresponding
// row in the node table.
func AuthenticateNode(ctx context.Context, dbc *sqlx.DB, key, secret string) (*Node, error) {
	stmt, err := dbc.Preparex("SELECT * FROM node WHERE public_key=$1 AND secret=$2;")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowxContext(ctx, key, secret)

	var node Node
	if err := row.StructScan(&node); err != nil {
		return nil, fmt.Errorf("retrieve record from table: %w", err)
	}

	return &node, nil
}

// CreateNode takes a name and a description and returns a created node with the key
// and secret filled out.
func CreateNode(ctx context.Context, dbc *sqlx.DB, name, description string) (*Node, error) {
	stmt, err := dbc.Preparex("INSERT INTO node (\"name\", description) VALUES ($1, $2) RETURNING public_key, secret;")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowxContext(ctx, name, description)

	var node Node
	if err := row.StructScan(&node); err != nil {
		return nil, fmt.Errorf("retrieve created node: %w", err)
	}

	// Set name and description on created node.
	node.Name = name
	node.Description = description

	return &node, nil
}
