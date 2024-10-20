package databasePerformance

import "fmt"

func DP_TestPostgresql(connectstr, certFilePath string) error {
	dbInst := getTestPostgresql()
	err := dbInst.Initialize(connectstr, certFilePath)
	if err != nil {
		fmt.Println("postgresql test error: ", err)
		return err
	}

	return testWorker(dbInst)
}
