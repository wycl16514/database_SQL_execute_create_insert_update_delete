package query

type ProductScan struct {
	scan1 Scan
	scan2 Scan
}

func NewProductScan(s1 Scan, s2 Scan) *ProductScan {
	p := &ProductScan{
		scan1: s1,
		scan2: s2,
	}

	p.scan1.Next()
	return p
}

func (p *ProductScan) BeforeFirst() {
	p.scan1.BeforeFirst()
	p.scan1.Next()
	p.scan2.BeforeFirst()
}

func (p *ProductScan) Next() bool {
	if p.scan2.Next() {
		return true
	} else {
		p.scan2.BeforeFirst()
		return p.scan2.Next() && p.scan1.Next()
	}
}

func (p *ProductScan) GetInt(fldName string) int {
	if p.scan1.HasField(fldName) {
		return p.scan1.GetInt(fldName)
	} else {
		return p.scan2.GetInt(fldName)
	}
}

func (p *ProductScan) GetString(fldName string) string {
	if p.scan1.HasField(fldName) {
		return p.scan1.GetString(fldName)
	} else {
		return p.scan2.GetString(fldName)
	}
}

func (p *ProductScan) GetVal(fldName string) *Constant {
	if p.scan1.HasField(fldName) {
		return p.scan1.GetVal(fldName)
	} else {
		return p.scan2.GetVal(fldName)
	}
}

func (p *ProductScan) HasField(fldName string) bool {
	return p.scan1.HasField(fldName) || p.scan2.HasField(fldName)
}

func (p *ProductScan) Close() {
	p.scan1.Close()
	p.scan2.Close()
}
