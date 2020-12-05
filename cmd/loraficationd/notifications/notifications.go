// Package notifications handles interfacing between the notification_contracts
// table and the lorafication daemon.
package notifications

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Contract represents a row in the notification_contracts table in the database.
type Contract struct {
	ID       int       `db:"id"`
	NodeID   int       `db:"node_id"`
	EntityID int       `db:"entity_id"`
	Email    *string   `db:"email"`
	SMS      *int      `db:"sms"`
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}

// ResolveContracts takes a nodeID and resolves all of the contracts stored in the
// database associated with this node ID in the notification_contracts table.
func ResolveContracts(dbc *sqlx.DB, nodeID int) ([]Contract, error) {
	stmt, err := dbc.Preparex("SELECT * FROM notification_contracts WHERE node_id=$1")
	if err != nil {
		return nil, fmt.Errorf("prepare statement to get notification contracts: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(nodeID)
	if err != nil {
		return nil, fmt.Errorf("query database for notification contracts: %w", err)
	}
	defer rows.Close()

	var contract Contract
	var contracts []Contract

	for rows.Next() {
		if err := rows.StructScan(&contract); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		contracts = append(contracts, contract)
	}

	return contracts, nil
}
