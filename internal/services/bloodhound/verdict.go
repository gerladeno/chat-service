package bloodhound

import "fmt"

var VerdictOK = NewVerdict()

// Verdict contains facts about the analyzed message.
// If the message is "clean" than verdict equals VerdictOK.
type Verdict struct {
	facts map[Fact]struct{}
}

type Fact int

const (
	FactContainsCardCVC Fact = iota
	FactContainsCardNumber
	FactContainsCardDate
	FactContainsSMSCode
	FactContainsSuspiciousPhrases
)

func NewVerdict(facts ...Fact) Verdict {
	v := Verdict{facts: make(map[Fact]struct{}, len(facts))}
	for _, fact := range facts {
		v.facts[fact] = struct{}{}
	}
	return v
}

func (v Verdict) String() string {
	return fmt.Sprintf("%v", v.facts)
}

func (v Verdict) Equals(rhs Verdict) bool {
	if v.FactsNumber() != rhs.FactsNumber() {
		return false
	}
	for fact := range v.facts {
		if !rhs.HasFact(fact) {
			return false
		}
	}
	return true
}

func (v *Verdict) AddFact(f Fact) {
	if v.facts == nil {
		v.facts = make(map[Fact]struct{})
	}
	v.facts[f] = struct{}{}
}

func (v Verdict) HasFact(f Fact) bool {
	_, ok := v.facts[f]
	return ok
}

func (v Verdict) FactsNumber() int {
	return len(v.facts)
}
