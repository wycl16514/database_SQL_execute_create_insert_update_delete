package parser

import (
	"record_manager"
)

type CreateTableData struct {
	tblName string
	sch     *record_manager.Schema
}

func NewCreateTableData(tblName string, sch *record_manager.Schema) *CreateTableData {
	return &CreateTableData{
		tblName: tblName,
		sch:     sch,
	}
}

func (c *CreateTableData) TableName() string {
	return c.tblName
}

func (c *CreateTableData) NewSchema() *record_manager.Schema {
	return c.sch
}
