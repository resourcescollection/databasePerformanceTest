package databasePerformance

import "fmt"

func DP_TestPostgresql(connectstr, certFilePath string) error {
	dbInst := getTestPostgresql()
	err := dbInst.Initialize(connectstr, certFilePath)
	if err != nil {
		fmt.Println("postgresql test error: ", err)
		return err
	}

	testWorker1(dbInst)
	dbInst.Close()

	// work1 end;
	// work2 start;
	// work2 := newTestWorker2()
	// err := work2.start(connectstr, certFilePath)

	return err
}
