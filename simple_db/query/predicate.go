package query

import (
	"record_manager"
)

type Predicate struct {
	terms []*Term
}

func NewPredicate() *Predicate {
	return &Predicate{}
}

func NewPredicateWithTerms(t *Term) *Predicate {
	predicate := &Predicate{}
	predicate.terms = make([]*Term, 0)
	predicate.terms = append(predicate.terms, t)
	return predicate
}

func (p *Predicate) ConjoinWith(pred *Predicate) {
	p.terms = append(p.terms, pred.terms...)
}

func (p *Predicate) IsSatisfied(s Scan) bool {
	for _, t := range p.terms {
		if !t.IsSatisfied(s) {
			return false
		}
	}

	return true
}

func (p *Predicate) ReductionFactor(plan Plan) int {

	factor := 1
	for _, t := range p.terms {
		factor *= t.ReductionFactor(plan)
	}

	return factor
}

func (p *Predicate) SelectedSubPred(sch *record_manager.Schema) *Predicate {
	result := NewPredicate()
	for _, t := range p.terms {
		if t.AppliesTo(sch) {
			result.terms = append(result.terms, t)
		}
	}

	if len(result.terms) == 0 {
		return nil
	}

	return result
}

func (p *Predicate) JoinSubPred(sch1 *record_manager.Schema, sch2 *record_manager.Schema) *Predicate {
	result := NewPredicate()
	newSch := record_manager.NewSchema()
	newSch.AddAll(sch1)
	newSch.AddAll(sch2)
	for _, t := range p.terms {
		if !t.AppliesTo(sch1) && !t.AppliesTo(sch2) && t.AppliesTo(newSch) {
			result.terms = append(result.terms, t)
		}
	}

	if len(result.terms) == 0 {
		return nil
	}
	return result
}

func (p *Predicate) EquatesWithConstant(fldName string) *Constant {
	for _, t := range p.terms {
		c := t.EquatesWithConstant(fldName)
		if c != nil {
			return c
		}
	}

	return nil
}

func (p *Predicate) EquatesWithField(fldName string) string {
	for _, t := range p.terms {
		s := t.EquatesWithField(fldName)
		if s != "" {
			return s
		}
	}

	return ""
}

func (p *Predicate) ToString() string {
	result := ""
	for _, t := range p.terms {
		result += " and " + t.ToString()
	}

	return result
}
