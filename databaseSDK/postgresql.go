package databaseSDK

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
)

type CPostgresql struct {
	connectStr string
	database   *sql.DB
}

func NewPostgressql(constr string) *CPostgresql {
	return &CPostgresql{connectStr: constr}
}

func (pInst *CPostgresql) Connect(certFilePath string) error {
	conn, _ := url.Parse(pInst.connectStr)
	//strSSL := "sslmode=disable"
	//strSSL := "sslmode=require"
	//if certFilePath != "" {
	//	strSSL = "sslmode=verify-full&sslrootcert=" + certFilePath
	//}
	conn.RawQuery = certFilePath

	db, err := sql.Open("postgres", conn.String())
	if err != nil {
		return err
	}

	pInst.database = db

	return nil
}
func (pInst *CPostgresql) Close() {
	pInst.database.Close()
}

func (pInst *CPostgresql) ExecSql(sql string) (sql.Result, error) {
	return pInst.database.Exec(sql)
}
func (pInst *CPostgresql) Query(sql string) (*sql.Rows, error) {
	return pInst.database.Query(sql)
}

func (pInst *CPostgresql) GetDatabaseVersion() (string, error) {
	rows, err := pInst.database.Query("SELECT version()")
	if err != nil {
		return "query error: " + err.Error(), err
	}

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

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = '%s')", tableName)
	err := pInst.database.QueryRow(query).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}
func (pInst *CPostgresql) DropTableIfExists(tableName string) error {
	strSql := "DROP TABLE IF EXISTS " + tableName
	_, err := pInst.database.Exec(strSql)
	return err
}
