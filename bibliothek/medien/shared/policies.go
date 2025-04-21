package shared

import (
	"time"
)

type LeihregelPolicy interface {
	DauerFuer(mediumTyp MediumType) time.Duration
}

type StandardLeihregelPolicy struct{}

func NewStandardLeihregelPolicy() LeihregelPolicy {
	return &StandardLeihregelPolicy{}
}

func (p *StandardLeihregelPolicy) DauerFuer(mediumTyp MediumType) time.Duration {
	switch mediumTyp {
	case MediumTypBuch:
		return 21 * 24 * time.Hour // 3 Wochen
	case MediumTypDVD:
		return 7 * 24 * time.Hour // 1 Woche
	case MediumTypZeitschrift:
		return 14 * 24 * time.Hour // 2 Wochen
	case MediumTypEBook:
		return 30 * 24 * time.Hour // 1 Monat
	default:
		return 14 * 24 * time.Hour // Standardwert
	}
}
