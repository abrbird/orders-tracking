package migrations

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/config"
	"log"
	"strings"
)

func init() {
	goose.AddMigration(upAddOrderHistoryRecords, downAddOrderHistoryRecords)
}

const (
	TableName        = "order_history_records"
	shardServerNameF = "shard_%d"
	shardTableNameF  = "order_history_records_shard_%d"
)

func upAddOrderHistoryRecords(tx *sql.Tx) error {
	cfg, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	mainTableCreation := fmt.Sprintf(`
		CREATE TABLE public.%s (
		    id serial,
		    order_id bigint NOT NULL,
			status VARCHAR NOT NULL,
			confirmation VARCHAR NOT NULL,
			updated_at TIMESTAMP NOT NULL default current_timestamp
		)
		PARTITION BY hash (order_id);
	`, TableName)
	shardCreationF := `
		CREATE SERVER IF NOT EXISTS %s FOREIGN DATA WRAPPER postgres_fdw
			OPTIONS (
				dbname '%s',
				host '%s',
				port '%d'
			);
		CREATE USER MAPPING IF NOT EXISTS FOR %s SERVER %s 
			OPTIONS (user '%s', password '%s');
	`
	shardTableCreationF := `
		CREATE FOREIGN TABLE IF NOT EXISTS public.%s
		PARTITION OF public.%s
		FOR VALUES WITH (modulus %d, remainder %d) 
		server %s;
	`

	queryList := []string{
		mainTableCreation,
		`CREATE EXTENSION IF NOT EXISTS postgres_fdw;`,
		//fmt.Sprintf(`GRANT USAGE ON FOREIGN DATA WRAPPER postgres_fdw to %s;`, cfg.Database.User),
	}

	for i, shardParam := range cfg.Database.Shards {
		shardServerName := fmt.Sprintf(shardServerNameF, i)
		shardTableName := fmt.Sprintf(shardTableNameF, i)

		queryList = append(
			queryList,
			fmt.Sprintf(
				shardCreationF,
				shardServerName,
				cfg.Database.DB,
				shardParam.Host,
				shardParam.Port,
				cfg.Database.User,
				shardServerName,
				cfg.Database.User,
				cfg.Database.Password,
			),
			fmt.Sprintf(
				shardTableCreationF,
				shardTableName,
				TableName,
				len(cfg.Database.Shards),
				i,
				shardServerName,
			),
		)
	}
	queryList = append(
		queryList,
	)

	query := strings.Join(queryList, "")

	_, err = tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func downAddOrderHistoryRecords(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`DROP TABLE %s;`, TableName))
	if err != nil {
		return err
	}
	return nil
}
