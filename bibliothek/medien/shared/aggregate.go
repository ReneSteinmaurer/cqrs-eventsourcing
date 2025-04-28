package shared

import (
	"cqrs-playground/bibliothek/medien/ausleihen/events"
	events2 "cqrs-playground/bibliothek/medien/erwerben/events"
	events3 "cqrs-playground/bibliothek/medien/katalogisieren/events"
	events4 "cqrs-playground/bibliothek/medien/rueckgeben/events"
	events5 "cqrs-playground/bibliothek/medien/verlieren/events"
	events6 "cqrs-playground/bibliothek/medien/wiederaufgefunden/events"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"time"
)

var MediumAggregateType = "medium"

type MediumAggregate struct {
	MediumId            string
	ISBN                string
	MediumType          events2.MediumType
	Name                string
	Genre               string
	Katalogisiert       bool
	VerliehenVon        *time.Time
	VerliehenBis        *time.Time
	VerliehenAnNutzerId string
	Verloren            bool
	VerlorenVonNutzerId string
}

func NewMediumAggregate(events []shared.Event) *MediumAggregate {
	m := &MediumAggregate{}
	for _, event := range events {
		m.Apply(event)
	}
	return m
}

func (m *MediumAggregate) Apply(event shared.Event) {
	switch event.Type {
	case events2.MediumErworbenEventType:
		var e events2.MediumErworbenEvent
		_ = json.Unmarshal(event.Payload, &e)
		m.ISBN = e.ISBN
		m.MediumId = e.MediumId
		m.MediumType = e.MediumType
		m.Name = e.Name
		m.Genre = e.Genre
	case events3.MediumKatalogisiertEventType:
		m.Katalogisiert = true
	case events.MediumVerliehenEventType:
		var e events.MediumVerliehenEvent
		_ = json.Unmarshal(event.Payload, &e)
		m.VerliehenVon = &e.Von
		m.VerliehenBis = &e.Bis
		m.VerliehenAnNutzerId = e.NutzerId
	case events4.MediumZurueckgegebenEventType:
		m.VerliehenVon = nil
		m.VerliehenBis = nil
		m.VerliehenAnNutzerId = ""
	case events5.MediumVerlorenDurchBenutzerEventType:
		var e events5.MediumVerlorenDurchBenutzerEvent
		_ = json.Unmarshal(event.Payload, &e)
		m.VerliehenVon = nil
		m.VerliehenBis = nil
		m.VerliehenAnNutzerId = ""
		m.Verloren = true
		m.VerlorenVonNutzerId = e.NutzerId
	case events5.MediumBestandsverlustEventType:
		m.VerliehenVon = nil
		m.VerliehenBis = nil
		m.VerliehenAnNutzerId = ""
		m.Verloren = true
	case events6.MediumWiederaufgefundenEventType:
		m.Verloren = false
	}
}

func (m *MediumAggregate) HandleMediumErwerben(event events2.MediumErworbenEvent) error {
	if m.alreadyExists(event.ISBN) {
		return errors.New("ein Medium mit dieser Id existiert bereits im Bestand")
	}
	return nil
}

func (m *MediumAggregate) HandleMediumVerleihen(event events.MediumVerliehenEvent) error {
	if m.MediumId != event.MediumId {
		return errors.New("es gibt noch kein Medium mit dieser Id im System")
	}
	if m.Verloren {
		return errors.New("das medium ist als verloren markiert und kann nicht verleihen werden")
	}
	if !m.isKatalogisiert() {
		return errors.New("das Medium ist noch nicht katalogisiert")
	}
	if m.VerliehenAnNutzerId == event.NutzerId {
		return errors.New("dieser nutzer hat das medium bereits ausgeliehen")
	}
	if m.isVerliehen() {
		return errors.New("das Medium wurde bereits verliehen und ist nicht verfügbar")
	}

	return nil
}

func (m *MediumAggregate) HandleMediumZurueckgegeben(event events4.MediumZurueckgegebenEvent) error {
	if m.MediumId != event.MediumId {
		return errors.New("es gibt noch kein Medium mit dieser Id im System")
	}
	if !m.isKatalogisiert() {
		return errors.New("das Medium ist noch nicht katalogisiert")
	}
	if m.VerliehenAnNutzerId != event.NutzerId {
		return errors.New("dieser nutzer hat das medium nicht ausgeliehen")
	}
	if !m.isVerliehen() {
		return errors.New("das Medium ist derzeit nicht verliehen")
	}

	return nil
}

func (m *MediumAggregate) HandleMediumVerlorenDurchBenutzer(event events5.MediumVerlorenDurchBenutzerEvent) error {
	if m.Verloren {
		return errors.New("das medium ist als verloren markiert und kann nicht verloren werden")
	}
	if m.MediumId != event.MediumId {
		return errors.New("es gibt noch kein Medium mit dieser Id im System")
	}
	if m.VerliehenAnNutzerId != event.NutzerId {
		return errors.New("dieser nutzer hat das medium nicht ausgeliehen")
	}
	if !m.isVerliehen() {
		return errors.New("das Medium ist derzeit nicht verliehen")
	}

	return nil
}

func (m *MediumAggregate) HandleMediumBestandsverlust(event events5.MediumBestandsverlustEvent) error {
	if m.Verloren {
		return errors.New("das medium ist als verloren markiert und kann nicht verloren werden")
	}
	if m.MediumId != event.MediumId {
		return errors.New("es gibt noch kein Medium mit dieser Id im System")
	}
	return nil
}

func (m *MediumAggregate) HandleMediumWiederaufgefunden(event events6.MediumWiederaufgefundenEvent) error {
	if m.MediumId != event.MediumId {
		return errors.New("es gibt noch kein Medium mit dieser Id im System")
	}
	if !m.Verloren {
		return errors.New("das Medium ist nicht als verloren markiert und kann somit nicht wiederaufgefunden werden")
	}
	if m.isVerliehen() {
		return errors.New("das Medium ist verliehen und kann somit nicht als wiederaufgefunden markiert werden")
	}
	return nil
}

func (m *MediumAggregate) HandleMediumKatalogisieren(event events3.MediumKatalogisiertEvent) error {
	if m.MediumId != event.MediumId {
		return errors.New("es gibt noch kein Medium mit dieser Id im System")
	}
	if m.isKatalogisiert() {
		return errors.New("das Medium ist bereits katalogisiert")
	}

	return nil
}

func (m *MediumAggregate) isVerliehen() bool {
	return m.VerliehenVon != nil && m.VerliehenBis != nil
}

func (m *MediumAggregate) isKatalogisiert() bool {
	return m.Katalogisiert
}

func (m *MediumAggregate) alreadyExists(isbn string) bool {
	return m.ISBN == isbn && m.ISBN != ""
}
