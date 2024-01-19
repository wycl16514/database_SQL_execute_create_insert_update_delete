package query

type ProjectScan struct {
	scan      Scan
	fieldList []string
}

func NewProjectScan(s Scan, fieldList []string) *ProjectScan {
	return &ProjectScan{
		scan:      s,
		fieldList: fieldList,
	}
}

func (p *ProjectScan) BeforeFirst() {
	p.scan.BeforeFirst()
}

func (p *ProjectScan) Next() bool {
	return p.scan.Next()
}

func (p *ProjectScan) GetInt(fldName string) int {
	if p.scan.HasField(fldName) {
		return p.scan.GetInt(fldName)
	}

	panic("Field Not Found")
}

func (p *ProjectScan) GetString(fldName string) string {
	if p.scan.HasField(fldName) {
		return p.scan.GetString(fldName)
	}

	panic("Field Not Found")
}

func (p *ProjectScan) GetVal(fldName string) *Constant {
	if p.scan.HasField(fldName) {
		return p.scan.GetVal(fldName)
	}

	panic("Field Not Found")
}

func (p *ProjectScan) HasField(fldName string) bool {
	for _, s := range p.fieldList {
		if s == fldName {
			return true
		}
	}

	return false
}

func (p *ProjectScan) Close() {
	p.scan.Close()
}
