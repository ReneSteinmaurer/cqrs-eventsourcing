package shared

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MediumType string

const (
	MediumTypBuch        MediumType = "BUCH"
	MediumTypDVD         MediumType = "DVD"
	MediumTypeCD         MediumType = "CD"
	MediumTypZeitschrift MediumType = "ZEITSCHRIFT"
	MediumTypEBook       MediumType = "EBOOK"
)

func (m *MediumType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch MediumType(strings.ToUpper(s)) {
	case MediumTypBuch, MediumTypeCD, MediumTypDVD, MediumTypZeitschrift, MediumTypEBook:
		*m = MediumType(strings.ToUpper(s))
		return nil
	default:
		return fmt.Errorf("ungültiger MediumType: %s", s)
	}
}
