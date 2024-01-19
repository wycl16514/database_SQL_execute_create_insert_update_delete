package planner

import (
	"parser"
	"record_manager"
	"tx"
)

type Plan interface {
	Open() interface{}
	BlocksAccessed() int               //对应 B(s)
	RecordsOutput() int                //对应 R(s)
	DistinctValues(fldName string) int //对应 V(s,F)
	Schema() record_manager.SchemaInterface
}

type QueryPlanner interface {
	CreatePlan(data *parser.QueryData, tx *tx.Transation) Plan
}

type UpdatePlanner interface {
	/*
		解释执行 insert 语句，返回被修改的记录条数
	*/
	ExecuteInsert(data *parser.InsertData, tx *tx.Transation) int

	/*
		解释执行 delete 语句，返回被删除的记录数
	*/
	ExecuteDelete(data *parser.DeleteData, tx *tx.Transation) int

	/*
		解释执行 update 语句，返回被修改的记录数
	*/
	ExecuteModify(data *parser.ModifyData, tx *tx.Transation) int

	/*
		解释执行 create table 语句，返回新建表中的记录数
	*/
	ExecuteCreateTable(data *parser.CreateTableData, tx *tx.Transation) int

	/*
		解释执行 create index 语句，返回当前建立了索引的记录数

	*/
	ExecuteCreateIndex(data *parser.IndexData, tx *tx.Transation) int
}
