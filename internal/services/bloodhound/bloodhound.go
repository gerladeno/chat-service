package bloodhound

import (
	"regexp"
	"strings"
)

var (
	patternCardCVC    = regexp.MustCompile(`(\D|^|\n)\d{3}(\D|$|\n)`)
	patternCardNumber = regexp.MustCompile(`\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}`)
	patternCardDate   = regexp.MustCompile(`\d{2}/\d{2}`)
	patternSMSCode    = regexp.MustCompile(`10-\d{4}`)

	suspiciousPhrases = map[string]struct{}{
		`login`:    {},
		`password`: {},
		`логин`:    {},
		`парол`:    {},
	}
)

// Search tries to find sensitive information in the message and returns Verdict.
func Search(msg string) (Verdict, error) {
	v := NewVerdict()
	if patternCardCVC.MatchString(msg) {
		v.AddFact(FactContainsCardCVC)
	}
	if patternCardNumber.MatchString(msg) {
		v.AddFact(FactContainsCardNumber)
	}
	if patternCardDate.MatchString(msg) {
		v.AddFact(FactContainsCardDate)
	}
	if patternSMSCode.MatchString(msg) {
		v.AddFact(FactContainsSMSCode)
	}
	if containsSuspiciousPhrases(msg) {
		v.AddFact(FactContainsSuspiciousPhrases)
	}
	return v, nil
}

func containsSuspiciousPhrases(s string) bool {
	for k := range suspiciousPhrases {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}
