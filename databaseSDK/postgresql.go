package databaseSDK

import (
	"context"
	"fmt"
	"time"

	//_ "github.com/lib/pq"
	//_ "github.com/Kount/pq-timeouts"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type CPostgresql struct {
	connectStr string
	database   *pgx.Conn
}

func NewPostgressql(constr string) *CPostgresql {
	return &CPostgresql{connectStr: constr}
}

func (pInst *CPostgresql) Connect(certFilePath string) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := pgx.Connect(context, pInst.connectStr)

	//conn, _ := url.Parse(pInst.connectStr)
	//strSSL := "sslmode=disable"
	//strSSL := "sslmode=require"
	//if certFilePath != "" {
	//	strSSL = "sslmode=verify-full&sslrootcert=" + certFilePath
	//}
	//conn.RawQuery = certFilePath

	//db, err := sql.Open("pq-timeouts", pInst.connectStr+"?read_timeout=1500&write_timeout=2000")
	if err != nil {
		return err
	}

	pInst.database = conn

	return nil
}
func (pInst *CPostgresql) Close() {
	pInst.database.Close(context.Background())
}

func (pInst *CPostgresql) ExecSql(sql string) (pgconn.CommandTag, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return pInst.database.Exec(context, sql)
}
func (pInst *CPostgresql) Query(sql string) (pgx.Rows, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//return pInst.database.Exec(context, sql)
	return pInst.database.Query(context, sql)
}

func (pInst *CPostgresql) GetDatabaseVersion() (string, error) {
	rows, err := pInst.Query("SELECT version()") //pInst.database.Query("SELECT version()")
	if err != nil {
		return "query error: " + err.Error(), err
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return err.Error(), err
		}
	}
	return result, nil
}

func (pInst *CPostgresql) CheckTableExists(tableName string) bool {

	var exists bool
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = '%s')", tableName)
	err := pInst.database.QueryRow(context, query).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}
func (pInst *CPostgresql) DropTableIfExists(tableName string) error {
	strSql := "DROP TABLE IF EXISTS " + tableName
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := pInst.database.Exec(context, strSql)
	return err
}
