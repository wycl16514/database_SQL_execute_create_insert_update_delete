package planner

import (
	"query"
	"record_manager"
)

type SelectPlan struct {
	p    Plan
	pred *query.Predicate
}

func NewSelectPlan(p Plan, pred *query.Predicate) *SelectPlan {
	return &SelectPlan{
		p:    p,
		pred: pred,
	}
}

func (s *SelectPlan) Open() interface{} {
	scan := s.p.Open()
	updateScan, ok := scan.(query.UpdateScan)
	if !ok {
		updateScanWrapper := query.NewUpdateScanWrapper(scan.(query.Scan))
		return query.NewSelectionScan(updateScanWrapper, s.pred)
	}
	return query.NewSelectionScan(updateScan, s.pred)
}

func (s *SelectPlan) BlocksAccessed() int {
	return s.p.BlocksAccessed()
}

func (s *SelectPlan) RecordsOutput() int {
	/*
			这里是一种预估,假设 student 有 90 条记录，根据原来我们在 StatInfo 中做的假设，
		    也就是给定字段取不同值的数量是总数的 1/3，于是假设表中有字段 age，那么根据假设，他有 31 种
		    不同取值可能，于是当 where 过滤条件为 age=20,那么我们预计满足条件的记录有 90/31=2 条
	*/
	return s.p.RecordsOutput() / s.pred.ReductionFactor(s.p)
}

func (s *SelectPlan) min(a int, b int) int {
	if a <= b {
		return a
	}

	return b
}

func (s *SelectPlan) DistinctValues(fldName string) int {
	if s.pred.EquatesWithConstant(fldName) != nil {
		//如果查询是 A=c 类型，A 是字段，c 是常量，那么查询结果返回一条数据
		return 1
	} else {
		//如果查询是 A=B 类型，A,B 都是字段，那么查询结果返回不同类型数值较小的那个字段
		fldName2 := s.pred.EquatesWithField(fldName)
		if fldName2 != "" {
			return s.min(s.p.DistinctValues(fldName), s.p.DistinctValues(fldName2))
		} else {
			return s.p.DistinctValues(fldName)
		}
	}
}

func (s *SelectPlan) Schema() record_manager.SchemaInterface {
	return s.p.Schema()
}
