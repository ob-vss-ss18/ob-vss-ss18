package moviestore

import (
	"fmt"
)

// FSK is an unsigned 8-bit int
type FSK uint8

// Serial is an unsigned int
type Serial uint

// See https://de.wikipedia.org/wiki/Freiwillige_Selbstkontrolle_der_Filmwirtschaft
const (
	FSK0  = 0
	FSK6  = 6
	FSK12 = 12
	FSK16 = 16
	FSK18 = 18
)

// A Movie consists of a title, the fsk minimum age and a serial
type Movie struct {
	Title  string
	Fsk    FSK
	Serial Serial
}

// Returns a string representing the movie like that:
// 		"  23. Texas Chainsaw Massacre (Ab 18)"
// The serial field should be 4 digits wide.
func (m *Movie) String() string {
	return fmt.Sprintf("%4d. %s (Ab %d)", m.Serial, m.Title, m.Fsk)
}

// AllowedAtAge checks whether the movie is allowed at a given age or not.
func (m *Movie) AllowedAtAge(age Age) bool {
	return uint8(m.Fsk) <= uint8(age)
}
