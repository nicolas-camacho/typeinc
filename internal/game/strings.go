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

	MenuHR      string // main-menu meta-shop entry (shows HR points)
	HRTitle     string
	HRHint      string
	HRPointsLbl string // currency label inside the meta shop
	DayLabel    string // HUD "DÍA"
	QuotaLabel  string // HUD "CUOTA"
	PaydayTitle string // end-of-day summary, quota covered
	FiredTitle  string // end-of-day summary, fired
	DayEndHint  string
	MaxedLabel  string // shown instead of a price on capped upgrades

	StreakLabel string // combo readout under the word
	MenuStats   string // main-menu stats entry
	StatsTitle  string
	StatsWords  string // words-typed counter (day summary + global)
	StatsFails  string // mistakes counter
	StatsWorst  string // most-failed word (day summary)
	StatsTop    string // top-3 most-failed header (global)
	StatsHint   string
	NoneLabel   string // placeholder when there is no worst word yet

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

		MenuHR:      "RRHH",
		HRTitle:     "OFICINA DE RRHH",
		HRHint:      "ENTER COMPRA · ESC VUELVE",
		HRPointsLbl: "PUNTOS RRHH",
		DayLabel:    "DÍA",
		QuotaLabel:  "CUOTA",
		PaydayTitle: "FIN DEL DÍA",
		FiredTitle:  "DESPEDIDO",
		DayEndHint:  "ENTER",
		MaxedLabel:  "MÁXIMO",

		StreakLabel: "RACHA",
		MenuStats:   "ESTADÍSTICAS",
		StatsTitle:  "ESTADÍSTICAS",
		StatsWords:  "PALABRAS",
		StatsFails:  "FALLOS",
		StatsWorst:  "PEOR PALABRA",
		StatsTop:    "TOP FALLOS",
		StatsHint:   "ESC VUELVE",
		NoneLabel:   "—",

		UpgradeNames: map[string]string{
			"mult":      "MULTIPLICADOR",
			"streakcap": "RACHA MÁXIMA",
			"hrmult":    "BONO PERMANENTE",
			"hrday":     "JORNADA EXTENDIDA",
			"hrquota":   "CUOTA REDUCIDA",
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

		MenuHR:      "HR",
		HRTitle:     "HR OFFICE",
		HRHint:      "ENTER BUYS · ESC BACK",
		HRPointsLbl: "HR POINTS",
		DayLabel:    "DAY",
		QuotaLabel:  "QUOTA",
		PaydayTitle: "END OF DAY",
		FiredTitle:  "FIRED",
		DayEndHint:  "ENTER",
		MaxedLabel:  "MAXED",

		StreakLabel: "STREAK",
		MenuStats:   "STATS",
		StatsTitle:  "STATISTICS",
		StatsWords:  "WORDS",
		StatsFails:  "FAILS",
		StatsWorst:  "WORST WORD",
		StatsTop:    "TOP FAILS",
		StatsHint:   "ESC BACK",
		NoneLabel:   "—",

		UpgradeNames: map[string]string{
			"mult":      "MULTIPLIER",
			"streakcap": "MAX STREAK",
			"hrmult":    "PERMANENT BONUS",
			"hrday":     "EXTENDED SHIFT",
			"hrquota":   "REDUCED QUOTA",
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

// PaydayScripts are sarcastic HR one-liners for a covered quota, one at
// random per day end (never the same twice in a row). Both languages have
// the same count so the picked index is valid in either.
var PaydayScripts = map[string][]string{
	"es": {
		"Su comisión ha sido procesada. De nada.",
		"RRHH agradece su contribución involuntaria al fondo de fiestas que no haremos.",
		"Cuota cobrada. El dinero va a un lugar seguro. Confíe.",
		"Pago recibido. Su lealtad ha sido anotada en una hoja de cálculo.",
		"Otro día productivo. Para nosotros, principalmente.",
		"Transacción completada. Sonría, es política de la empresa.",
	},
	"en": {
		"Your commission has been processed. You're welcome.",
		"HR thanks you for your involuntary contribution to the party fund for parties we won't throw.",
		"Quota collected. The money goes somewhere safe. Trust us.",
		"Payment received. Your loyalty has been noted in a spreadsheet.",
		"Another productive day. For us, mostly.",
		"Transaction complete. Smile, it's company policy.",
	},
}

// FiredScripts are HR one-liners for a missed quota (game over). Both
// languages have the same count.
var FiredScripts = map[string][]string{
	"es": {
		"RRHH le desea éxito en sus futuros proyectos. Lejos de aquí.",
		"No es usted, es su rendimiento. Bueno, sí es usted.",
		"Su puesto ha sido eliminado. Su silla ya tiene otro dueño.",
		"El comité revisó su caso durante casi tres segundos. Despedido.",
		"Entregue el teclado a la salida. Las letras son de la empresa.",
		"Su carrera aquí termina hoy. Su cuota, lamentablemente, no se pagó sola.",
	},
	"en": {
		"HR wishes you success in your future endeavors. Far from here.",
		"It's not you, it's your performance. Well, it is you.",
		"Your position has been eliminated. Your chair already has a new owner.",
		"The committee reviewed your case for almost three seconds. Fired.",
		"Hand in the keyboard on your way out. The letters belong to the company.",
		"Your career here ends today. Your quota, sadly, didn't pay itself.",
	},
}

// HRExplainScripts replaces the random payday quip at the end of the
// first day: HR explains (in its own way) what HR points are, one line
// per Enter, intro-style. Both languages have the same line count.
var HRExplainScripts = map[string][]string{
	"es": {
		"Sobrevivió a su primer día. Felicidades. El listón estaba en el suelo.",
		"Hablemos de su compensación: los PUNTOS RRHH.",
		"Son una moneda inventada por el departamento para no pagarle con dinero real. No cotizan en bolsa, no compran café y no pueden heredarse.",
		"Pero canjéelos en la OFICINA DE RRHH (comando /rrhh) y su vida laboral dolerá un poco menos.",
		"Y tranquilo: incluso si lo despedimos, se los quedará. Somos crueles, no ladrones.",
	},
	"en": {
		"You survived your first day. Congratulations. The bar was on the floor.",
		"Let's talk about your compensation: HR POINTS.",
		"They're a currency the department invented to avoid paying you real money. Not listed on any stock exchange, they don't buy coffee, and they can't be inherited.",
		"But redeem them at the HR OFFICE (/hr command) and your work life will hurt a little less.",
		"And relax: even if we fire you, you keep them. We're cruel, not thieves.",
	},
}

// DayStartScripts are plain HR morning quips shown when a new day begins.
// Both languages have the same count.
var DayStartScripts = map[string][]string{
	"es": {
		"Nuevo día. Mismas letras. A trabajar.",
		"RRHH le recuerda que la motivación viene incluida en su salario. Que no hay.",
		"Comienza otra jornada productiva. Defina 'productiva' como prefiera.",
		"El teclado ya está encendido. Usted también, suponemos.",
		"Hoy es un gran día para teclear. Todos los días lo son. No hay más días.",
	},
	"en": {
		"New day. Same letters. Get to work.",
		"HR reminds you that motivation is included in your salary. Which doesn't exist.",
		"Another productive day begins. Define 'productive' however you like.",
		"The keyboard is already on. So are you, we assume.",
		"Today is a great day for typing. Every day is. There are no other days.",
	},
}

// DayStartStatScripts are morning quips that mock yesterday's numbers.
// Placeholders: {words}, {fails}, {worst}, {worstn}. Both languages have
// the same count.
var DayStartStatScripts = map[string][]string{
	"es": {
		"Ayer escribió {words} palabras. RRHH esperaba más. RRHH siempre espera más.",
		"Registro de ayer: {fails} fallos. La tecla correcta estaba justo al lado. Siempre lo está.",
		"La palabra '{worst}' le ganó {worstn} veces ayer. Hoy es su revancha. O la de ella.",
		"{words} palabras y {fails} fallos ayer. La hoja de cálculo no juzga. Nosotros sí.",
	},
	"en": {
		"Yesterday you typed {words} words. HR expected more. HR always expects more.",
		"Yesterday's record: {fails} misses. The right key was right next door. It always is.",
		"The word '{worst}' beat you {worstn} times yesterday. Today is your rematch. Or its.",
		"{words} words and {fails} misses yesterday. The spreadsheet doesn't judge. We do.",
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
