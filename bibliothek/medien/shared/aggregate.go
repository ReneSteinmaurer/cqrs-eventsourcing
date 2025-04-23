package shared

import (
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"time"
)

var MediumAggregateType = "medium"

type MediumAggregate struct {
	MediumId            string
	ISBN                string
	MediumType          MediumType
	Name                string
	Genre               string
	Katalogisiert       bool
	VerliehenVon        *time.Time
	VerliehenBis        *time.Time
	VerliehenAnNutzerId string
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
	case MediumErworbenEventType:
		var e MediumErworbenEvent
		_ = json.Unmarshal(event.Payload, &e)
		m.ISBN = e.ISBN
		m.MediumId = e.MediumId
		m.MediumType = e.MediumType
		m.Name = e.Name
		m.Genre = e.Genre
	case MediumKatalogisiertEventType:
		var e MediumKatalogisiertEvent
		_ = json.Unmarshal(event.Payload, &e)
		m.Katalogisiert = true
	case MediumVerliehenEventType:
		var e MediumVerliehenEvent
		_ = json.Unmarshal(event.Payload, &e)
		m.VerliehenVon = &e.Von
		m.VerliehenBis = &e.Bis
		m.VerliehenAnNutzerId = e.NutzerId
	case MediumZurueckgegebenEventType:
		var e MediumZurueckgegebenEvent
		_ = json.Unmarshal(event.Payload, &e)
		m.VerliehenVon = nil
		m.VerliehenBis = nil
		m.VerliehenAnNutzerId = ""
	}
}

func (m *MediumAggregate) HandleMediumErwerben(event MediumErworbenEvent) error {
	if m.alreadyExists(event.ISBN) {
		return errors.New("ein Medium mit dieser Id existiert bereits im Bestand")
	}
	return nil
}

func (m *MediumAggregate) HandleMediumVerleihen(event MediumVerliehenEvent) error {
	if m.MediumId != event.MediumId {
		return errors.New("es gibt noch kein Medium mit dieser Id im System")
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

func (m *MediumAggregate) HandleMediumZurueckgegeben(event MediumZurueckgegebenEvent) error {
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

func (m *MediumAggregate) HandleMediumKatalogisieren(event MediumKatalogisiertEvent) error {
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
