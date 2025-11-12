package cfbr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/jackc/pgx/v5"
)

var tableName = "cfbr_seasons"

type CfbrRepo struct {
	conn *pgx.Conn
}

func CfbrRepository() (CfbrRepo, error) {
	connString := core.EnvGetOrPanic("NEON_CONNECTION")

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return CfbrRepo{}, err
	}

	repo := CfbrRepo{
		conn: conn,
	}

	repo.createTable()

	return repo, nil
}

func (repo *CfbrRepo) createTable() error {

	qry := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v (
    id TEXT PRIMARY KEY,
    year INT NOT NULL,
    created_at BIGINT NOT NULL,
    compressed_season BYTEA NOT NULL
	)`, tableName)

	_, err := repo.conn.Exec(context.Background(), qry)
	if err != nil {
		return err
	}

	log.Println("Table created (if not exists).")
	return nil
}

func (repo *CfbrRepo) insert(szn SerializeableCompressedSeason) error {
	qry := fmt.Sprintf(`INSERT INTO %v (id, year, created_at, compressed_season)
     VALUES ($1, $2, $3, $4)`, tableName)

	_, err := repo.conn.Exec(
		context.Background(), qry,
		szn.Id,
		szn.Year,
		time.Now().Unix(),
		szn.CompressedSeason,
	)

	return err
}

func (repo *CfbrRepo) get(id string) (SerializeableCompressedSeason, error) {
	qry := fmt.Sprintf(`SELECT id, year, created_at, compressed_season
	 FROM %v 
	 WHERE id=$1`,
		tableName)

	szn := SerializeableCompressedSeason{}

	err := repo.conn.QueryRow(context.Background(), qry, id).Scan(
		&szn.Id,
		&szn.Year,
		&szn.CreatedAt,
		&szn.CompressedSeason)
	if err != nil {
		return szn, err
	}
	return szn, nil
}
