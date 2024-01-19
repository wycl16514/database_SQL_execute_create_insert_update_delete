package planner

import (
	metadata_manager "metadata_management"
	"parser"
	"tx"
)

type BasicQueryPlanner struct {
	mdm *metadata_manager.MetaDataManager
}

func CreateBasicQueryPlanner(mdm *metadata_manager.MetaDataManager) QueryPlanner {
	return &BasicQueryPlanner{
		mdm: mdm,
	}
}

func (b *BasicQueryPlanner) CreatePlan(data *parser.QueryData, tx *tx.Transation) Plan {
	//1,直接创建 QueryData 对象中的表
	plans := make([]Plan, 0)
	tables := data.Tables()
	for _, tblname := range tables {
		//获取该表对应视图的 sql 代码
		viewDef := b.mdm.GetViewDef(tblname, tx)
		if viewDef != "" {
			//直接创建表对应的视图
			parser := parser.NewSQLParser(viewDef)
			viewData := parser.Query()
			//递归的创建对应表的规划器
			plans = append(plans, b.CreatePlan(viewData, tx))
		} else {
			plans = append(plans, NewTablePlan(tx, tblname, b.mdm))
		}
	}

	//将所有表执行 Product 操作，注意表的次序会对后续查询效率有重大影响，但这里我们不考虑表的次序，只是按照
	//给定表依次执行 Product 操作，后续我们会在这里进行优化
	p := plans[0]
	plans = plans[1:]

	for _, nextPlan := range plans {
		p = NewProductPlan(p, nextPlan)
	}

	p = NewSelectPlan(p, data.Pred())

	return NewProjectPlan(p, data.Fields())
}
