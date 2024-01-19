package query

import (
	"math"
	"record_manager"
)

type Term struct {
	lhs *Expression
	rhs *Expression
}

func NewTerm(lhs *Expression, rhs *Expression) *Term {
	return &Term{
		lhs,
		rhs,
	}
}

func (t *Term) IsSatisfied(s Scan) bool {
	lhsVal := t.lhs.Evaluate(s)
	rhsVal := t.rhs.Evaluate(s)
	return rhsVal.Equals(lhsVal)
}

func (t *Term) AppliesTo(sch *record_manager.Schema) bool {
	return t.lhs.AppliesTo(sch) && t.rhs.AppliesTo(sch)
}

func (t *Term) ReductionFactor(p Plan) int {
	//Plan是后面我们研究SQL解析执行时才创建的对象，
	lhsName := ""
	rhsName := ""
	if t.lhs.IsFieldName() && t.rhs.IsFieldName() {
		lhsName = t.lhs.AsFieldName()
		rhsName = t.rhs.AsFieldName()
		if p.DistinctValues(lhsName) > p.DistinctValues(rhsName) {
			return p.DistinctValues(lhsName)
		}
		return p.DistinctValues(rhsName)
	}

	if t.lhs.IsFieldName() {
		lhsName = t.lhs.AsFieldName()
		return p.DistinctValues(lhsName)
	}

	if t.rhs.IsFieldName() {
		rhsName = t.rhs.AsFieldName()
		return p.DistinctValues(rhsName)
	}

	if t.lhs.AsConstant().Equals(t.rhs.AsConstant()) {
		return 1
	} else {
		return math.MaxInt
	}
}

func (t *Term) EquatesWithConstant(fldName string) *Constant {
	if t.lhs.IsFieldName() && t.lhs.AsFieldName() == fldName && !t.rhs.IsFieldName() {
		return t.rhs.AsConstant()
	} else if t.rhs.IsFieldName() && t.rhs.AsFieldName() == fldName && !t.lhs.IsFieldName() {
		return t.lhs.AsConstant()
	} else {
		return nil
	}
}

func (t *Term) EquatesWithField(fldName string) string {
	if t.lhs.IsFieldName() && t.lhs.AsFieldName() == fldName && t.rhs.IsFieldName() {
		return t.rhs.AsFieldName()
	} else if t.rhs.IsFieldName() && t.rhs.AsFieldName() == fldName && t.lhs.IsFieldName() {
		return t.lhs.AsFieldName()
	}

	return ""
}

func (t *Term) ToString() string {
	return t.lhs.ToString() + "=" + t.rhs.ToString()
}
