package parser

import "fmt"

type IndexData struct {
	idxName string
	tblName string
	fldName string
}

func NewIndexData(idxName string, tblName string, fldName string) *IndexData {
	return &IndexData{
		idxName: idxName,
		tblName: tblName,
		fldName: fldName,
	}
}

func (i *IndexData) IdexName() string {
	return i.idxName
}

func (i *IndexData) tableName() string {
	return i.tblName
}

func (i *IndexData) fieldName() string {
	return i.fldName
}

func (i *IndexData) ToString() string {
	str := fmt.Sprintf("index name: %s, table name: %s, field name: %s", i.idxName, i.tblName, i.fldName)
	return str
}
