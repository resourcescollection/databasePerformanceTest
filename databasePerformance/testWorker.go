package databasePerformance

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type iTestDatabase interface {
	DropTable(string) error
	CreateTestTable(string) error
	ExecSql(string) error
	Query(sql string) (*sql.Rows, error)
}

const GC_TestTableName = "testtableperformance" //"testtableperformance"
const GC_TestRecordCount = 10000

func testWorker(dbInst iTestDatabase) error {
	err := dbInst.DropTable(GC_TestTableName)
	if err != nil {
		fmt.Print("drop table error: ", err)
		return err
	}

	err = dbInst.CreateTestTable(GC_TestTableName)
	if err != nil {
		fmt.Print("create table error: ", err)
		return err
	}

	testInsertPerformance(dbInst)

	checkRecordCount(dbInst)

	testInsertRepeatData(dbInst)

	testUpdateData(dbInst)

	return nil
}

func testInsertPerformance(dbInst iTestDatabase) {
	timeBegin := time.Now().UnixMicro()
	strSql := ""
	writeTimes := 0
	for iLoop := 1; iLoop <= GC_TestRecordCount; iLoop++ {
		strSql += fmt.Sprintf("insert into %s(id, key, valuestr, valueint, valuefloat) values (%d,'key%d', 'test insert performance %d', %d, %f);\n",
			GC_TestTableName, iLoop, iLoop, iLoop, iLoop, float64(iLoop)*0.11)
		if iLoop%20 == 0 {
			err := dbInst.ExecSql(strSql)
			if err != nil {
				fmt.Print("insert data error: ", err, " id:", iLoop)
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

	textReport := "test start:" + strconv.FormatInt(timeBegin, 10)
	textReport += "\ntest   end:" + strconv.FormatInt(timeEnd, 10)
	textReport += "\ntest   use:" + strconv.FormatInt(timeUse, 10)
	textReport += "  (" + strconv.FormatInt(timeUse/(1000*1000), 10) + " seconds)"
	textReport += "\ntest writetimes:" + strconv.Itoa(writeTimes)

	fmt.Println(textReport)
}
func testInsertRepeatData(dbInst iTestDatabase) {
	timeBegin := time.Now().UnixMicro()
	strSql := ""
	writeTimes := 0
	iLoop := 1
	for iLoop <= GC_TestRecordCount {
		strSql = fmt.Sprintf("insert into %s(id, key, valuestr, valueint, valuefloat) values (%d,'key%d', 'test insert performance %d', %d, %f);\n",
			GC_TestTableName, iLoop, iLoop, iLoop, iLoop, float64(iLoop)*0.11)

		dbInst.ExecSql(strSql)
		writeTimes++
		iLoop += 20
	}
	timeEnd := time.Now().UnixMicro()

	timeUse := timeEnd - timeBegin

	textReport := "=====repeat data test: "
	textReport += "\ntest   use:" + strconv.FormatInt(timeUse, 10)
	textReport += "\nseconds:" + strconv.FormatInt(timeUse/(1000*1000), 10)
	textReport += "\ntest writetimes:" + strconv.Itoa(writeTimes)

	fmt.Println(textReport)
}
func testUpdateData(dbInst iTestDatabase) {
	timeBegin := time.Now().UnixMicro()
	strSql := ""
	writeTimes := 0
	iLoop := 1
	for iLoop <= GC_TestRecordCount {
		strSql = fmt.Sprintf("update %s set key='newkey%d', valuestr='test update performance%d', valueint=%d, valuefloat=%f where id=%d;\n",
			GC_TestTableName, iLoop, iLoop, iLoop, float64(iLoop)*11, iLoop)
		strSql += fmt.Sprintf("update %s set key='newkey%d', valuestr='test update performance%d', valueint=%d, valuefloat=%f where id=%d;\n",
			GC_TestTableName, iLoop, iLoop, iLoop, float64(iLoop)*11, iLoop+1)
		strSql += fmt.Sprintf("update %s set key='newkey%d', valuestr='test update performance%d', valueint=%d, valuefloat=%f where id=%d;",
			GC_TestTableName, iLoop, iLoop, iLoop, float64(iLoop)*11, iLoop+2)

		dbInst.ExecSql(strSql)
		writeTimes++
		iLoop += 20
	}
	timeEnd := time.Now().UnixMicro()

	timeUse := timeEnd - timeBegin

	textReport := "=====update data test: "
	textReport += "\ntest   use:" + strconv.FormatInt(timeUse, 10)
	textReport += "\nseconds:" + strconv.FormatInt(timeUse/(1000*1000), 10)
	textReport += "\ntest writetimes:" + strconv.Itoa(writeTimes)

	fmt.Println(textReport)
}

func checkRecordCount(dbInst iTestDatabase) error {
	sql := "select count(*) as rCount from " + GC_TestTableName + ";"
	rows, err := dbInst.Query(sql)
	if err != nil {
		return err
	}
	var nCount int = 0
	for rows.Next() {
		rows.Scan(&nCount)
	}

	if nCount == GC_TestRecordCount {
		fmt.Println("inserted datacount correct: ", nCount)
	} else {
		fmt.Println("inserted datacount incorrect: ", nCount)
	}

	return nil
}
