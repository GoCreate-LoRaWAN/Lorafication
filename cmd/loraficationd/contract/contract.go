// Package contract handles interfacing between the contract table and the lorafication
// daemon.
package contract

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Contract is a struct representing the structure of a row in the contract table
// of the database.
type Contract struct {
	ID            int       `db:"id"`
	NodePublicKey string    `db:"node_public_key"`
	EntityID      int       `db:"entity_id"`
	Created       time.Time `db:"created"`
	Modified      time.Time `db:"modified"`
}

// CreateContract takes a node public key and an entity ID and creates a row in the contract
// table.
func CreateContract(ctx context.Context, dbc *sqlx.DB, nodePublicKey string, entityID int) error {
	stmt, err := dbc.Preparex("INSERT INTO contract (node_public_key, entity_id) VALUES ($1, $2);")
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, nodePublicKey, entityID); err != nil {
		return fmt.Errorf("execute statement: %w", err)
	}

	return nil
}

// ResolvedContract represents a row returned in the complex query used in ResolveContracts.
type ResolvedContract struct {
	SMS   *int    `db:"sms"`
	Email *string `db:"email"`
}

// ResolveContracts takes a node public key and resolves all of the notification contracts
// that are paired with it. The returned result is each entity that is subscribed to said
// node's Email and/or SMS number.
func ResolveContracts(ctx context.Context, dbc *sqlx.DB, nodePublicKey string) ([]ResolvedContract, error) {
	stmt, err := dbc.PreparexContext(ctx, `SELECT
  sms,
  email
FROM
  contract
  INNER JOIN node ON contract.node_public_key = node.public_key
  INNER JOIN entity ON contract.entity_id = entity.id
WHERE
  node.public_key = $1;`)

	if err != nil {
		return nil, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, nodePublicKey)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()

	var contract ResolvedContract
	var contracts []ResolvedContract

	for rows.Next() {
		if err := rows.StructScan(&contract); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		contracts = append(contracts, contract)
	}

	return contracts, nil
}
