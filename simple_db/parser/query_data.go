package parser

//QueryData 用来描述select语句的操作信息
import (
	"query"
)

type QueryData struct {
	fields []string
	tables []string
	pred   *query.Predicate
}

func NewQueryData(fields []string, tables []string, pred *query.Predicate) *QueryData {
	return &QueryData{
		fields: fields,
		tables: tables,
		pred:   pred,
	}
}

func (q *QueryData) Fields() []string {
	return q.fields
}

func (q *QueryData) Tables() []string {
	return q.tables
}

func (q *QueryData) Pred() *query.Predicate {
	return q.pred
}

func (q *QueryData) ToString() string {
	result := "select "
	for _, fldName := range q.fields {
		result += fldName + ", "
	}

	// 去掉最后一个逗号
	result = result[:len(result)-1]
	result += " from "
	for _, tableName := range q.tables {
		result += tableName + ", "
	}
	// 去掉最后一个逗号
	result = result[:len(result)-1]
	predStr := q.pred.ToString()
	if predStr != "" {
		result += " where " + predStr
	}

	return result
}
