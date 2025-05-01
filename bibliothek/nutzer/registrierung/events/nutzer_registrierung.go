package events

const NutzerRegistriertEventType = "NutzerRegistriertEvent"

type NutzerRegistriertEvent struct {
	Vorname  string
	Nachname string
	Email    string
	NutzerId string
}

func NewNutzerRegistrierungEvent(Email, Vorname, Nachname, NutzerId string) NutzerRegistriertEvent {
	return NutzerRegistriertEvent{
		Vorname:  Vorname,
		Nachname: Nachname,
		Email:    Email,
		NutzerId: NutzerId,
	}
}
