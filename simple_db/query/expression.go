package query

import (
	"record_manager"
)

//expr -> expr op term

type Expression struct {
	val     *Constant
	fldName string
}

func NewExpressionWithConstant(val *Constant) *Expression {
	return &Expression{
		val:     val,
		fldName: "",
	}
}

func NewExpressionWithString(fldName string) *Expression {
	return &Expression{
		val:     nil,
		fldName: fldName,
	}
}

func (e *Expression) IsFieldName() bool {
	return e.fldName != ""
}

func (e *Expression) AsConstant() *Constant {
	return e.val
}

func (e *Expression) AsFieldName() string {
	return e.fldName
}

func (e *Expression) Evaluate(s Scan) *Constant {
	/*
		expression 有可能对应一个常量，或者对应一个字段名，如果是后者，那么我们需要查询该字段对应的具体值
	*/
	if e.val != nil {
		return e.val
	}

	return s.GetVal(e.fldName)
}

func (e *Expression) AppliesTo(sch *record_manager.Schema) bool {
	if e.val != nil {
		return true
	}

	return sch.HasFields(e.fldName)
}

func (e *Expression) ToString() string {
	if e.val != nil {
		return e.val.ToString()
	}

	return e.fldName
}
