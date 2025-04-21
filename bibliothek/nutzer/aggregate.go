package nutzer

import (
	"cqrs-playground/bibliothek/nutzer/events"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"net/mail"
)

type NutzerAggregate struct {
	NutzerId string
	Email    string
	Vorname  string
	Nachname string
}

func NewNutzerAggregate(events []shared.Event) *NutzerAggregate {
	n := &NutzerAggregate{}
	for _, event := range events {
		n.Apply(event)
	}
	return n
}

func (n *NutzerAggregate) Apply(event shared.Event) {
	switch event.Type {
	case events.NutzerRegistriertEventType:
		var e events.NutzerRegistriertEvent
		_ = json.Unmarshal(event.Payload, &e)
		n.NutzerId = e.NutzerId
		n.Vorname = e.Vorname
		n.Nachname = e.Nachname
		n.Email = e.Email
	}
}

func (n *NutzerAggregate) HandleRegistriereNutzer(event events.NutzerRegistriertEvent) error {
	if !n.isEmailValid(event.Email) {
		return errors.New("die Email ist nicht valide")
	}
	if n.isEmailAlreadyRegistered(event.Email) {
		return errors.New("es ist bereits ein Nutzer mit dieser Email registriert")
	}
	return nil
}

func (n *NutzerAggregate) isEmailAlreadyRegistered(email string) bool {
	return n.Email == email
}

func (n *NutzerAggregate) isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
