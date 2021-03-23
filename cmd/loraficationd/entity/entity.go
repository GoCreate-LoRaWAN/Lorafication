// Package entity interfaces between the entity table in the database and the
// lorafication daemon.
package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Entity is a struct representing the structure of a row in the entity table
// of the database.
type Entity struct {
	ID       int       `db:"id"`
	Name     string    `db:"name"`
	Email    *string   `db:"email"`
	SMS      *int      `db:"sms"`
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}

// CreateEntity takes a name, email, and sms where email and sms are both optional
// (but at least one needs provided due to a database constraint) and creates a row
// in the entity table in the database.
func CreateEntity(ctx context.Context, dbc *sqlx.DB, name string, email *string, sms *int) (*Entity, error) {
	stmt, err := dbc.Preparex("INSERT INTO entity (\"name\", email, sms) VALUES ($1, $2, $3);")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, name, email, sms); err != nil {
		return nil, fmt.Errorf("execute statement: %w", err)
	}

	// None of the other fields should matter to the client.
	return &Entity{
		Name:  name,
		Email: email,
		SMS:   sms,
	}, nil
}
