package main

import (
	bmg "buffer_manager"
	fm "file_manager"
	"fmt"
	lm "log_manager"
	metadata_manager "metadata_management"
	"parser"
	"planner"
	"query"
	"tx"
)

func PrintStudentTable(tx *tx.Transation, mdm *metadata_manager.MetaDataManager) {
	queryStr := "select name, majorid, gradyear from STUDENT"
	p := parser.NewSQLParser(queryStr)
	queryData := p.Query()
	test_planner := planner.CreateBasicQueryPlanner(mdm)
	test_plan := test_planner.CreatePlan(queryData, tx)
	test_interface := (test_plan.Open())
	test_scan, _ := test_interface.(query.Scan)
	for test_scan.Next() {
		fmt.Printf("name: %s, majorid: %d, gradyear: %d\n",
			test_scan.GetString("name"), test_scan.GetInt("majorid"),
			test_scan.GetInt("gradyear"))
	}
}

func CreateInsertUpdateByUpdatePlanner() {
	file_manager, _ := fm.NewFileManager("student", 2048)
	log_manager, _ := lm.NewLogManager(file_manager, "logfile.log")
	buffer_manager := bmg.NewBufferManager(file_manager, log_manager, 3)
	tx := tx.NewTransation(file_manager, log_manager, buffer_manager)
	mdm := metadata_manager.NewMetaDataManager(false, tx)

	updatePlanner := planner.NewBasicUpdatePlanner(mdm)
	createTableSql := "create table STUDENT (name varchar(16), majorid int, gradyear int)"
	p := parser.NewSQLParser(createTableSql)
	tableData := p.UpdateCmd().(*parser.CreateTableData)
	updatePlanner.ExecuteCreateTable(tableData, tx)

	insertSQL := "insert into STUDENT (name, majorid, gradyear) values(\"tylor\", 30, 2020)"
	p = parser.NewSQLParser(insertSQL)
	insertData := p.UpdateCmd().(*parser.InsertData)
	updatePlanner.ExecuteInsert(insertData, tx)
	insertSQL = "insert into STUDENT (name, majorid, gradyear) values(\"tom\", 35, 2023)"
	p = parser.NewSQLParser(insertSQL)
	insertData = p.UpdateCmd().(*parser.InsertData)
	updatePlanner.ExecuteInsert(insertData, tx)

	fmt.Println("table after insert:")
	PrintStudentTable(tx, mdm)

	updateSQL := "update STUDENT set majorid=20 where majorid=30 and gradyear=2020"
	p = parser.NewSQLParser(updateSQL)
	updateData := p.UpdateCmd().(*parser.ModifyData)
	updatePlanner.ExecuteModify(updateData, tx)

	fmt.Println("table after update:")
	PrintStudentTable(tx, mdm)

	deleteSQL := "delete from STUDENT where majorid=35"
	p = parser.NewSQLParser(deleteSQL)
	deleteData := p.UpdateCmd().(*parser.DeleteData)
	updatePlanner.ExecuteDelete(deleteData, tx)

	fmt.Println("table after delete")
	PrintStudentTable(tx, mdm)
}

func main() {
	CreateInsertUpdateByUpdatePlanner()
}
