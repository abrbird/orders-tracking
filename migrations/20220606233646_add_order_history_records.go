package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddOrderHistoryRecords, downAddOrderHistoryRecords)
}

func upAddOrderHistoryRecords(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE order_history_records (
		    id serial PRIMARY KEY,
		    order_id bigint NOT NULL,
			status VARCHAR NOT NULL,
			confirmation VARCHAR NOT NULL,
			updated_at TIMESTAMP NOT NULL default current_timestamp,
		    UNIQUE (order_id, status)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downAddOrderHistoryRecords(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE order_history_records;")
	if err != nil {
		return err
	}
	return nil
}
