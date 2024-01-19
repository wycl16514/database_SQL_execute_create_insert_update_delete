package planner

import (
	"query"
	"record_manager"
)

type ProjectPlan struct {
	p      Plan
	schema *record_manager.Schema
}

func NewProjectPlan(p Plan, fieldList []string) *ProjectPlan {
	project_plan := ProjectPlan{
		p:      p,
		schema: record_manager.NewSchema(),
	}

	for _, field := range fieldList {
		project_plan.schema.Add(field, project_plan.p.Schema())
	}

	return &project_plan
}

func (p *ProjectPlan) Open() interface{} {
	s := p.p.Open()
	return query.NewProjectScan(s.(query.Scan), p.schema.Fields())
}

func (p *ProjectPlan) BlocksAccessed() int {
	return p.p.BlocksAccessed()
}

func (p *ProjectPlan) RecordsOutput() int {
	return p.p.RecordsOutput()
}

func (p *ProjectPlan) DistinctValues(fldName string) int {
	return p.p.DistinctValues(fldName)
}

func (p *ProjectPlan) Schema() record_manager.SchemaInterface {
	return p.schema
}
