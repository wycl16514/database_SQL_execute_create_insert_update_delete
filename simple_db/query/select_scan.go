package query

import (
	"record_manager"
)

type SelectionScan struct {
	updateScan UpdateScan
	pred       *Predicate
}

func NewSelectionScan(s UpdateScan, pred *Predicate) *SelectionScan {
	return &SelectionScan{
		updateScan: s,
		pred:       pred,
	}
}

func (s *SelectionScan) BeforeFirst() {
	s.updateScan.GetScan().BeforeFirst()
}

func (s *SelectionScan) Next() bool {
	for s.updateScan.GetScan().Next() {
		if s.pred.IsSatisfied(s.updateScan.GetScan()) {
			return true
		}
	}

	return false
}

func (s *SelectionScan) GetInt(fldName string) int {
	return s.updateScan.GetScan().GetInt(fldName)
}

func (s *SelectionScan) GetString(fldName string) string {
	return s.updateScan.GetScan().GetString(fldName)
}

func (s *SelectionScan) GetVal(fldName string) *Constant {
	return s.updateScan.GetScan().GetVal(fldName)
}

func (s *SelectionScan) HasField(fldName string) bool {
	return s.updateScan.GetScan().HasField(fldName)
}

func (s *SelectionScan) Close() {
	s.updateScan.GetScan().Close()
}

func (s *SelectionScan) SetInt(fldName string, val int) {
	s.updateScan.SetInt(fldName, val)
}

func (s *SelectionScan) SetString(fldName string, val string) {
	s.updateScan.SetString(fldName, val)
}

func (s *SelectionScan) SetVal(fldName string, val *Constant) {
	s.updateScan.SetVal(fldName, val)
}

func (s *SelectionScan) Delete() {
	s.updateScan.Delete()
}

func (s *SelectionScan) Insert() {
	s.updateScan.Insert()
}

func (s *SelectionScan) GetRid() *record_manager.RID {
	return s.updateScan.GetRid()
}

func (s *SelectionScan) MoveToRID(rid *record_manager.RID) {
	s.updateScan.MoveToRid(rid)
}
