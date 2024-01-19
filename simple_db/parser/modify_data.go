package parser

import "query"

type ModifyData struct {
	tblName string
	fldName string
	newVal  *query.Expression
	pred    *query.Predicate
}

func NewModifyData(tblName string, fldName string, newVal *query.Expression, pred *query.Predicate) *ModifyData {
	return &ModifyData{
		tblName: tblName,
		fldName: fldName,
		newVal:  newVal,
		pred:    pred,
	}
}

func (m *ModifyData) TableName() string {
	return m.tblName
}

func (m *ModifyData) TargetField() string {
	return m.fldName
}

func (m *ModifyData) NewValue() *query.Expression {
	return m.newVal
}

func (m *ModifyData) Pred() *query.Predicate {
	return m.pred
}
