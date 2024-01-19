package parser

import (
	"query"
)

type InsertData struct {
	tblName string
	flds    []string
	vals    []*query.Constant
}

func NewInsertData(tblName string, flds []string, vals []*query.Constant) *InsertData {
	return &InsertData{
		tblName: tblName,
		flds:    flds,
		vals:    vals,
	}
}

func (i *InsertData) TableName() string {
	return i.tblName
}

func (i *InsertData) Fields() []string {
	return i.flds
}

func (i *InsertData) Vals() []*query.Constant {
	return i.vals
}
