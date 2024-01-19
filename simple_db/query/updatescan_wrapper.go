package query

import (
	"record_manager"
)

type UpdateScanWrapper struct {
	scan Scan
}

func NewUpdateScanWrapper(s Scan) *UpdateScanWrapper {
	return &UpdateScanWrapper{
		scan: s,
	}
}

func (u *UpdateScanWrapper) GetScan() Scan {
	return u.scan
}

func (u *UpdateScanWrapper) SetInt(fldName string, val int) {
	//DO NOTHING
}

func (u *UpdateScanWrapper) SetString(fldName string, val string) {
	//DO NOTHING
}

func (u *UpdateScanWrapper) SetVal(fldName string, val *Constant) {
	//DO NOTHING
}

func (u *UpdateScanWrapper) Insert() {
	//DO NOTHING
}

func (u *UpdateScanWrapper) Delete() {
	//DO NOTHING
}

func (u *UpdateScanWrapper) GetRid() *record_manager.RID {
	return nil
}

func (u *UpdateScanWrapper) MoveToRid(rid *record_manager.RID) {
	// DO NOTHING
}
