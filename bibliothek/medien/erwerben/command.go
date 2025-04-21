package erwerben

import "cqrs-playground/bibliothek/medien/shared"

type ErwerbeMediumCommand struct {
	ISBN       string
	MediumType shared.MediumType
	Name       string
	Genre      string
}
