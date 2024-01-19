package parser

import "fmt"

type ViewData struct {
	viewName  string
	queryData *QueryData
}

func NewViewData(viewName string, qd *QueryData) *ViewData {
	return &ViewData{
		viewName:  viewName,
		queryData: qd,
	}
}

func (v *ViewData) ViewName() string {
	return v.viewName
}

func (v *ViewData) ViewDef() string {
	return v.queryData.ToString()
}

func (v *ViewData) ToString() string {
	s := fmt.Sprintf("view name %s, viewe def: %s", v.viewName, v.ViewDef())
	return s
}
