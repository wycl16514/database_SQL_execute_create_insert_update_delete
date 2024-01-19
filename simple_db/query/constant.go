package query

import (
	"strconv"
)

type Constant struct {
	ival *int
	sval *string
}

func NewConstantWithInt(ival *int) *Constant {
	return &Constant{
		ival: ival,
		sval: nil,
	}
}

func NewConstantWithString(sval *string) *Constant {
	return &Constant{
		ival: nil,
		sval: sval,
	}
}

func (c *Constant) AsInt() int {
	return *c.ival
}

func (c *Constant) AsString() string {
	return *c.sval
}

func (c *Constant) Equals(obj *Constant) bool {
	if c.ival != nil && obj.ival != nil {
		return *c.ival == *obj.ival
	}

	if c.sval != nil && obj.sval != nil {
		return *c.sval == *obj.sval
	}

	return false
}

func (c *Constant) ToString() string {
	if c.ival != nil {
		return strconv.FormatInt((int64)(*c.ival), 10)
	}

	return *c.sval
}
