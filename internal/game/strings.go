package game

// UIStrings holds every user-facing interface string for one language.
type UIStrings struct {
	MenuContinue string
	MenuNew      string
	MenuLang     string
	MenuOptions  string
	MenuQuit     string

	OptDisplay   string
	OptVolume    string
	OptBack      string
	DisplayModes [3]string

	MenuHint    string
	OptionsHint string
	ShopHint    string
	IntroHint   string

	ShopTitle  string
	UnknownCmd string

	UpgradeNames map[string]string // keyed by Upgrade.ID
}

var uiTables = map[string]UIStrings{
	"es": {
		MenuContinue: "CONTINUAR",
		MenuNew:      "NUEVA PARTIDA",
		MenuLang:     "IDIOMA",
		MenuOptions:  "OPCIONES",
		MenuQuit:     "SALIR",

		OptDisplay:   "PANTALLA",
		OptVolume:    "SONIDO",
		OptBack:      "VOLVER",
		DisplayModes: [3]string{"COMPLETA", "VENTANA", "SIN BORDE"},

		MenuHint:    "FLECHAS + ENTER",
		OptionsHint: "FLECHAS AJUSTAN · ESC VUELVE",
		ShopHint:    "ENTER COMPRA · ESC VUELVE",
		IntroHint:   "ENTER",

		ShopTitle:  "TIENDA",
		UnknownCmd: "COMANDO DESCONOCIDO",

		UpgradeNames: map[string]string{
			"mult":      "MULTIPLICADOR",
			"streakcap": "RACHA MÁXIMA",
		},
	},
	"en": {
		MenuContinue: "CONTINUE",
		MenuNew:      "NEW GAME",
		MenuLang:     "LANGUAGE",
		MenuOptions:  "OPTIONS",
		MenuQuit:     "QUIT",

		OptDisplay:   "DISPLAY",
		OptVolume:    "SOUND",
		OptBack:      "BACK",
		DisplayModes: [3]string{"FULLSCREEN", "WINDOWED", "BORDERLESS"},

		MenuHint:    "ARROWS + ENTER",
		OptionsHint: "ARROWS ADJUST · ESC BACK",
		ShopHint:    "ENTER BUYS · ESC BACK",
		IntroHint:   "ENTER",

		ShopTitle:  "SHOP",
		UnknownCmd: "UNKNOWN COMMAND",

		UpgradeNames: map[string]string{
			"mult":      "MULTIPLIER",
			"streakcap": "MAX STREAK",
		},
	},
}

// IntroScripts is the HR welcome story per language; both versions have
// the same number of lines so switching languages mid-intro stays valid.
var IntroScripts = map[string][]string{
	"es": {
		"Bienvenido al equipo.",
		"Gracias por aceptar el puesto de mecanógrafo. ¿Firmaste el contrato sin leerlo? Perfecto, aquí encajarás de maravilla.",
		"Este es tu nuevo equipo de trabajo: una pantalla, unas letras y un cursor que parpadea. Tecnología de punta. No lo rompas.",
		"Harás TODO desde aquí. ¿Reuniones? No hay. ¿Café? Tráelo de casa. ¿Palabras? Todas las que aguanten tus dedos.",
		"Recuerda: tienes una tienda escribiendo /tienda y el menú de pausa con /menu. No preguntes quién instaló eso en tu teclado. Nadie lo sabe.",
		"¿Que si te pagaremos por escribir palabras?... Así es. La economía está rara últimamente.",
		"Eso es todo. Suerte, mecanógrafo. Y pase lo que pase... no dejes de teclear.",
	},
	"en": {
		"Welcome to the team.",
		"Thanks for accepting the typist position. Signed the contract without reading it? Perfect, you'll fit right in.",
		"This is your new work equipment: a screen, some letters and a blinking cursor. Cutting-edge technology. Don't break it.",
		"You'll do EVERYTHING from here. Meetings? None. Coffee? Bring your own. Words? As many as your fingers can take.",
		"Remember: there's a shop if you type /shop and the pause menu at /menu. Don't ask who installed that in your keyboard. Nobody knows.",
		"Will we really pay you for typing words?... That's right. The economy has been weird lately.",
		"That's all. Good luck, typist. And whatever happens... don't stop typing.",
	},
}

// ReturnScripts are short HR quips shown on CONTINUAR, one at random
// (never the same twice in a row). Both languages have the same count so
// the picked index is valid in either.
var ReturnScripts = map[string][]string{
	"es": {
		"Ah, volviste. RRHH había apostado en contra.",
		"Tu silla sigue caliente. No preguntes por qué.",
		"Las palabras no se escribieron solas en tu ausencia. Lo intentamos.",
		"Bienvenido de vuelta. El café sigue sin existir.",
		"Informe de RRHH: te fuiste, nadie lo notó. Sigue tecleando.",
		"Volviste justo a tiempo. ¿A tiempo de qué? Ni idea.",
		"El cursor te extrañó. Nosotros... también, claro. Claro.",
		"Recuerda: cada palabra cuenta. Literalmente, las contamos.",
	},
	"en": {
		"Oh, you're back. HR had bet against it.",
		"Your chair is still warm. Don't ask why.",
		"The words didn't type themselves while you were gone. We tried.",
		"Welcome back. The coffee still doesn't exist.",
		"HR report: you left, nobody noticed. Keep typing.",
		"You're back just in time. In time for what? No idea.",
		"The cursor missed you. We did... too, of course. Of course.",
		"Remember: every word counts. Literally, we count them.",
	},
}

// UI returns the interface strings for the selected language.
func (g *Game) UI() UIStrings {
	return uiTables[Languages[g.LangIndex].Code]
}

// UpgradeName returns the localized shop name of an upgrade.
func (g *Game) UpgradeName(id string) string {
	return g.UI().UpgradeNames[id]
}
