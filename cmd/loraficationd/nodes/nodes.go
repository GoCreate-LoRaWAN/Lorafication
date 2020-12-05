// Package nodes interfaces between the nodes table in the database and the
// lorafication daemon.
package nodes

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Node is a struct representing the structure of a row in the nodes table
// of the database.
type Node struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Created     time.Time `db:"created"`
	Modified    time.Time `db:"modified"`
}

// ResolveNode takes an ID and finds the corresponding row in the nodes table
// and returns it.
func ResolveNode(dbc *sqlx.DB, id int) (*Node, error) {
	stmt, err := dbc.Preparex("SELECT * FROM nodes WHERE id=$1")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowx(id)

	var node Node
	if err := row.StructScan(&node); err != nil {
		return nil, fmt.Errorf("retrieve record from table: %w", err)
	}

	return &node, nil
}
