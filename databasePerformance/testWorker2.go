package databasePerformance

import (
	"fmt"
	"strconv"
	"time"
)

type cTestWorker2 struct {
	dbList   []*cTestPostgresql
	timeList [c_threadCount]int64
	//logList  [c_threadCount]string
}

const c_threadCount = 3
const c_testRecordCount2 = GC_TestRecordCount / (c_threadCount + 2)

func newTestWorker2() *cTestWorker2 {
	return &cTestWorker2{}
}

func (pInst *cTestWorker2) start(connectstr, certFilePath string) error {
	iDBCount := 0
	var lastError error = nil
	for iLoop := 0; iLoop < c_threadCount; iLoop++ {
		dbInst := newTestProstgresql()
		err := dbInst.Initialize(connectstr, certFilePath)
		if err != nil {
			lastError = err
			fmt.Println("-----------------------------")
			fmt.Println("db connect error, index: ", iLoop+1, " error info: ", err)
			break
		}
		pInst.timeList[iLoop] = 0
		iDBCount++
		pInst.dbList = append(pInst.dbList, dbInst)
		time.Sleep(30 * time.Millisecond)
	}
	if iDBCount < 1 {
		fmt.Println("database connect error:", lastError)
		return lastError
	}
	fmt.Println("create db connecting: ", len(pInst.dbList))

	for iLoop := 0; iLoop < iDBCount; iLoop++ {
		//pInst.logList[iLoop] = ""
		go pInst.testInsertPerformance(iLoop, pInst.dbList[iLoop])
		time.Sleep(30 * time.Millisecond)
	}

	for iLoop := 0; iLoop < 100; iLoop++ {
		time.Sleep(2 * time.Second)
		var bFinished bool = true
		var indexNotYet int = -1
		for iLoop1 := 0; iLoop1 < iDBCount; iLoop1++ {
			if pInst.timeList[iLoop1] == 0 {
				indexNotYet = iLoop1
				bFinished = false
				break
			}
		}
		if bFinished {
			break
		}
		fmt.Println("thread still working: ", indexNotYet) //, pInst.logList[indexNotYet])
	}

	textReport := ""
	for iLoop := 0; iLoop < iDBCount; iLoop++ {
		textReport += "\n[" + strconv.Itoa(iLoop) + "] test   use:" + strconv.FormatInt(pInst.timeList[iLoop], 10)
		textReport += "  (" + strconv.FormatInt(pInst.timeList[iLoop]/(1000*1000), 10) + " seconds)"
	}

	fmt.Println(textReport)

	return nil
}

func (pInst *cTestWorker2) testInsertPerformance(index int, dbInst iTestDatabase) {
	//pInst.logList[index] += fmt.Sprintf("thread start %d\n", index)
	var tableName string = GC_TestTableName + strconv.Itoa(index+1)
	err := pInst.dbList[index].DropTable(tableName)
	if err != nil {
		pInst.timeList[index] = -1
		fmt.Print("drop table error: ", err, "index:", index)
		//pInst.logList[index] += fmt.Sprintf("tdrop table error: %d\n", index)
		return
	}
	//pInst.logList[index] += fmt.Sprintf("table droped %d\n", index)

	err = pInst.dbList[index].CreateTestTable(tableName)
	if err != nil {
		pInst.timeList[index] = -1
		fmt.Print("create table error: ", err, "index:", index)
		//pInst.logList[index] += fmt.Sprintf("create table error: %d\n", index)
		return
	}
	//pInst.logList[index] += fmt.Sprintf("table created %d\n", index)

	var idbase = 0 //GC_TestRecordCount * index
	timeBegin := time.Now().UnixMicro()
	strSql := ""
	writeTimes := 0
	for iLoop := 1; iLoop <= c_testRecordCount2; iLoop++ {
		iID := idbase + iLoop
		strSql += fmt.Sprintf("insert into %s(id, key, valuestr, valueint, valuefloat) values (%d,'key%d', 'test insert performance %d', %d, %f);\n",
			tableName, iID, iID, iID, iID, float64(iID)*0.11)
		if iLoop%10 == 0 {
			//pInst.logList[index] += fmt.Sprintf("start insert %d\n", index)
			err := dbInst.ExecSql(strSql)
			//pInst.logList[index] += fmt.Sprintf("insert done %d\n", index)
			if err != nil {
				fmt.Print("insert data error: ", err, " id:", iID)
				break
			}
			strSql = ""
			writeTimes++
		}
	}
	if strSql != "" {
		err := dbInst.ExecSql(strSql)
		writeTimes++
		if err != nil {
			fmt.Print("last insert data error: ", err)
		}
	}
	timeEnd := time.Now().UnixMicro()

	timeUse := timeEnd - timeBegin

	textReport := "[" + strconv.Itoa(index) + "] test start:" + strconv.FormatInt(timeBegin, 10)
	textReport += "\ntest   end:" + strconv.FormatInt(timeEnd, 10)
	textReport += "\ntest   use:" + strconv.FormatInt(timeUse, 10)
	textReport += "  (" + strconv.FormatInt(timeUse/(1000*1000), 10) + " seconds)"
	textReport += "\ntest writetimes:" + strconv.Itoa(writeTimes)

	pInst.timeList[index] = timeUse

	fmt.Println(textReport)
	err = pInst.checkRecordCount(index, dbInst)
	if err != nil {
		fmt.Println("check insert count error: ", err)
	}
	pInst.dbList[index].Close()
}

func (pInst *cTestWorker2) checkRecordCount(index int, dbInst iTestDatabase) error {
	var tableName string = GC_TestTableName + strconv.Itoa(index+1)
	//var idbase = GC_TestRecordCount * index
	//var idMax = GC_TestRecordCount * (index + 1)
	sql := "select count(*) as rCount from " + tableName + ";"
	rows, err := dbInst.Query(sql)
	if err != nil {
		return err
	}
	var nCount int = 0
	for rows.Next() {
		rows.Scan(&nCount)
	}

	if nCount == c_testRecordCount2 {
		fmt.Println("inserted datacount correct: ", nCount, " index: ", index)
	} else {
		fmt.Println("inserted datacount incorrect: ", nCount, " index: ", index)
	}

	return nil
}
