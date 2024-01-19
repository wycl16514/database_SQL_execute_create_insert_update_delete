package planner

import (
	metadata_manager "metadata_management"
	"parser"
	"query"
	"tx"
)

type BasicUpdatePlanner struct {
	mdm *metadata_manager.MetaDataManager
}

func NewBasicUpdatePlanner(mdm *metadata_manager.MetaDataManager) *BasicUpdatePlanner {
	return &BasicUpdatePlanner{
		mdm: mdm,
	}
}

func (b *BasicUpdatePlanner) ExecuteDelete(data *parser.DeleteData, tx *tx.Transation) int {
	/*
		先 scan 出要处理的记录，然后执行删除
	*/
	tablePlan := NewTablePlan(tx, data.TableName(), b.mdm)
	selectPlan := NewSelectPlan(tablePlan, data.Pred())
	scan := selectPlan.Open()
	updateScan := scan.(*query.SelectionScan)
	count := 0
	for updateScan.Next() {
		updateScan.Delete()
		count += 1
	}
	updateScan.Close()
	return count
}

func (b *BasicUpdatePlanner) ExecuteModify(data *parser.ModifyData, tx *tx.Transation) int {
	/*
		先 scan 出选中的记录，然后修改记录的信息
	*/
	tablePlan := NewTablePlan(tx, data.TableName(), b.mdm)
	selectPlan := NewSelectPlan(tablePlan, data.Pred())
	scan := selectPlan.Open()
	updateScan := scan.(*query.SelectionScan)
	count := 0
	for updateScan.Next() {
		val := data.NewValue().Evaluate(scan.(query.Scan))
		updateScan.SetVal(data.TargetField(), val)
		count += 1
	}
	return count
}

func (b *BasicUpdatePlanner) ExecuteInsert(data *parser.InsertData, tx *tx.Transation) int {
	tablePlan := NewTablePlan(tx, data.TableName(), b.mdm)
	updateScan := tablePlan.Open().(*query.TableScan)
	updateScan.Insert()
	insertFields := data.Fields()
	insertedVals := data.Vals()

	for i := 0; i < len(insertFields); i++ {
		updateScan.SetVal(insertFields[i], insertedVals[i])
	}

	updateScan.Close()
	return 1
}

func (b *BasicUpdatePlanner) ExecuteCreateTable(data *parser.CreateTableData, tx *tx.Transation) int {
	b.mdm.CreateTable(data.TableName(), data.NewSchema(), tx)
	return 0
}

func (b *BasicUpdatePlanner) ExecuteView(data *parser.ViewData, tx *tx.Transation) int {
	b.mdm.CreateView(data.ViewName(), data.ViewDef(), tx)
	return 0
}

func (b *BasicUpdatePlanner) ExecuteIndex(data *parser.IndexData, tx *tx.Transation) int {
	//b.mdm.CreateIndex
	//TODO

	return 0
}
