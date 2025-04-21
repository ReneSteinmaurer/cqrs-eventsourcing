package shared

type MediumErworbenEvent struct {
	ISBN       string
	MediumId   string
	MediumType MediumType
	Name       string
	Genre      string
}

func NewMediumErworbenEvent(isbn, mediumId, name, genre string, mediumType MediumType) MediumErworbenEvent {
	return MediumErworbenEvent{
		ISBN:       isbn,
		MediumId:   mediumId,
		MediumType: mediumType,
		Name:       name,
		Genre:      genre,
	}
}

var MediumErworbenEventType = "MediumErworbenEvent"
