package cfbr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/jackc/pgx/v5"
)

// type CfbrRepo = mongo.MongoConnection[SerializeableCompressedSeason]

// func CfbrRepository() CfbrRepo {

// 	connString := core.EnvGetOrPanic("NEON_CONNECTION")

// 	conn, err := pgx.Connect(context.Background(), connString)
// 	if err != nil {
// 		log.Fatalf("unable to connect to database: %v", err)
// 	}

// 	repo := CfbrRepo{
// 		conn: conn,
// 	}

// 	repo.createTable()

// 	// dbName := core.EnvGetOrPanic("DB_NAME_CFBR")
// 	// return CfbrRepo{
// 	// 	Db:         dbName,
// 	// 	Collection: "seasons",
// 	// }

// 	return repo
// }

var tableName = "cfbr_seasons"

type CfbrRepo struct {
	conn *pgx.Conn
}

func CfbrRepository() CfbrRepo {
	connString := core.EnvGetOrPanic("NEON_CONNECTION")

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	repo := CfbrRepo{
		conn: conn,
	}

	repo.createTable()

	// dbName := core.EnvGetOrPanic("DB_NAME_CFBR")
	// return CfbrRepo{
	// 	Db:         dbName,
	// 	Collection: "seasons",
	// }

	return repo
}

func (repo *CfbrRepo) createTable() {

	qry := fmt.Sprintf(`CREATE TABLE %v (
    id TEXT PRIMARY KEY,
    year INT NOT NULL,
    created_at BIGINT NOT NULL,
    compressed_season BYTEA NOT NULL`, tableName)

	_, err := repo.conn.Exec(context.Background(), qry)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	log.Println("Table created (if not exists).")
}

func (repo *CfbrRepo) insert(szn SerializeableCompressedSeason) {
	qry := fmt.Sprintf(`INSERT INTO %v (id, year, created_at, compressed_season)
     VALUES ($1, $2, $3, $4)`, tableName)

	_, err := repo.conn.Exec(
		context.Background(), qry,
		"season-2025",
		2025,
		time.Now().Unix(),
		szn,
	)
	if err != nil {
		log.Fatal(err)
	}
}
