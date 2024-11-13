package databasePerformance

import (
	"dbtest/databaseSDK"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type cTestPostgresql struct {
	connectStr string
	dbInst     *databaseSDK.CPostgresql
}

var g_TestPostgresql *cTestPostgresql = &cTestPostgresql{connectStr: ""}

func getTestPostgresql() *cTestPostgresql {
	return g_TestPostgresql
}

func newTestProstgresql() *cTestPostgresql {
	return &cTestPostgresql{connectStr: ""}
}

func (pInst *cTestPostgresql) Initialize(connectstr, certFilePath string) error {
	if !strings.HasPrefix(strings.ToLower(connectstr), "postgresql://") {
		return errors.New("connection string type error")
	}
	pInst.dbInst = databaseSDK.NewPostgressql(connectstr)
	err := pInst.dbInst.Connect(certFilePath)
	if err != nil {
		return err
	} else {
		fmt.Println("database connect successful")
	}

	strVersion, err := pInst.dbInst.GetDatabaseVersion()
	if err != nil {
		return err
	}
	fmt.Println(strVersion)

	return nil
}
func (pInst *cTestPostgresql) Close() {
	pInst.dbInst.Close()
}
func (pInst *cTestPostgresql) DropTable(tableName string) error {
	return pInst.dbInst.DropTableIfExists(tableName)
}
func (pInst *cTestPostgresql) CreateTestTable(tableName string) error {
	strSql := fmt.Sprintf(C_TestCreateTable, tableName)
	_, err := pInst.dbInst.ExecSql(strSql)

	return err
}
func (pInst *cTestPostgresql) ExecSql(sql string) error {
	_, err := pInst.dbInst.ExecSql(sql)
	return err
}
func (pInst *cTestPostgresql) Query(sql string) (pgx.Rows, error) {
	return pInst.dbInst.Query(sql)
}

const C_TestCreateTable = ` CREATE TABLE IF NOT EXISTS %s(
			id                    BIGSERIAL         NOT NULL PRIMARY KEY,
			key varchar(128),
			valuestr TEXT,
			valueint BigInt,
			valuefloat REAL,
			date1 TIMESTAMP DEFAULT NOW());`
