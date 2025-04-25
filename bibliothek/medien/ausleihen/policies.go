package ausleihen

import (
	"cqrs-playground/bibliothek/medien/erwerben/events"
	"time"
)

type LeihregelPolicy interface {
	DauerFuer(mediumTyp events.MediumType) time.Duration
}

type StandardLeihregelPolicy struct{}

func NewStandardLeihregelPolicy() LeihregelPolicy {
	return &StandardLeihregelPolicy{}
}

func (p *StandardLeihregelPolicy) DauerFuer(mediumTyp events.MediumType) time.Duration {
	switch mediumTyp {
	case events.MediumTypBuch:
		return 21 * 24 * time.Hour // 3 Wochen
	case events.MediumTypDVD:
		return 7 * 24 * time.Hour // 1 Woche
	case events.MediumTypeCD:
		return 7 * 24 * time.Hour // 1
	case events.MediumTypZeitschrift:
		return 14 * 24 * time.Hour // 2 Wochen
	case events.MediumTypEBook:
		return 30 * 24 * time.Hour // 1 Monat
	default:
		return 14 * 24 * time.Hour // Standardwert
	}
}
