package erwerben

import (
	"cqrs-playground/bibliothek/medien/erwerben/events"
)

type ErwerbeMediumCommand struct {
	ISBN       string
	MediumType events.MediumType
	Name       string
	Genre      string
}
