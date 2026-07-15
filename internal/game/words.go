package game

import (
	_ "embed"
	"strings"
)

// Full dictionary word pools, one word per line, pre-filtered to lowercase
// a-z (no accents, no ñ, no symbols), 3-14 letters, so any keyboard layout
// can play.
//
// Sources: dwyl/english-words (words_alpha) and olea/lemarios
// (lemario general del español).

//go:embed data/es.txt
var esRaw string

//go:embed data/en.txt
var enRaw string

// WordLists holds the words per language code.
var WordLists = map[string][]string{
	"es": strings.Fields(esRaw),
	"en": strings.Fields(enRaw),
}
