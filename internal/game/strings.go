package game

// UIStrings holds every user-facing interface string for one language.
type UIStrings struct {
	MenuContinue string
	MenuNew      string
	MenuLang     string
	MenuOptions  string
	MenuReset    string
	ResetConfirm string // RESET armed: shown until the confirming Enter
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

	MenuHR       string // main-menu meta-shop entry (shows HR points)
	HRTitle      string
	HRHint       string
	HRPointsLbl  string // currency label inside the meta shop
	DayLabel     string // HUD "DÍA"
	QuotaLabel   string // HUD "CUOTA"
	PaydayTitle  string // end-of-day summary, quota covered
	FiredTitle   string // end-of-day summary, fired
	NextPayLabel string // "next payroll: day N" hint on non-payroll days
	DayEndHint   string
	MaxedLabel   string // shown instead of a price on capped upgrades

	StreakLabel     string // combo readout under the word
	GoldenLabel     string // shown above a golden word
	FrenzyLabel     string // coffee-frenzy banner (with its countdown)
	InspectionLabel string // urgent-word banner (with its countdown)
	MenuStats       string // main-menu stats entry
	StatsTitle      string
	StatsWords      string // words-typed counter (day summary + global)
	StatsFails      string // mistakes counter
	StatsWorst      string // most-failed word (day summary)
	StatsStreak     string // best streak (day summary + global)
	StatsBestDay    string // highest day reached (global)
	StatsTop        string // top-3 most-failed header (global)
	StatsHint       string
	NoneLabel       string // placeholder when there is no worst word yet

	CorruptLabel string // shown above a corrupt word

	// hidden story terminal
	MenuTerminal    string // main-menu entry, only once discovered
	TermBanner      string
	TermHelpHint    string
	TermHelp        []string
	TermUnknown     string
	TermLocked      string // reading a keyed doc without its key
	TermLockedTag   string // listing placeholder for locked docs
	TermDocsHeader  string
	TermKeyAccepted string
	TermNoDocs      string
	TermHint        string // scene footer

	UpgradeNames map[string]string // keyed by Upgrade.ID
}

var uiTables = map[string]UIStrings{
	"es": {
		MenuContinue: "CONTINUAR",
		MenuNew:      "NUEVA PARTIDA",
		MenuLang:     "IDIOMA",
		MenuOptions:  "OPCIONES",
		MenuReset:    "REINICIAR",
		ResetConfirm: "¿SEGURO? ENTER BORRA TODO",
		MenuQuit:     "SALIR",

		OptDisplay:   "PANTALLA",
		OptVolume:    "SONIDO",
		OptBack:      "VOLVER",
		DisplayModes: [3]string{"COMPLETA", "VENTANA", "SIN BORDE"},

		MenuHint:    "FLECHAS + ENTER",
		OptionsHint: "FLECHAS AJUSTAN · ESC VUELVE",
		ShopHint:    "ENTER COMPRA · ESC VUELVE",
		IntroHint:   "ENTER · ESC SALTA",

		ShopTitle:  "TIENDA",
		UnknownCmd: "COMANDO DESCONOCIDO",

		MenuHR:       "RRHH",
		HRTitle:      "OFICINA DE RRHH",
		HRHint:       "ENTER COMPRA · ESC VUELVE",
		HRPointsLbl:  "PUNTOS RRHH",
		DayLabel:     "DÍA",
		QuotaLabel:   "CUOTA",
		PaydayTitle:  "FIN DEL DÍA",
		FiredTitle:   "DESPEDIDO",
		NextPayLabel: "PRÓXIMA NÓMINA: DÍA",
		DayEndHint:   "ENTER",
		MaxedLabel:   "MÁXIMO",

		StreakLabel:     "RACHA",
		GoldenLabel:     "¡DORADA!",
		FrenzyLabel:     "CAFÉ ×2",
		InspectionLabel: "INSPECCIÓN DE RRHH",
		MenuStats:       "ESTADÍSTICAS",
		StatsTitle:      "ESTADÍSTICAS",
		StatsWords:      "PALABRAS",
		StatsFails:      "FALLOS",
		StatsWorst:      "PEOR PALABRA",
		StatsStreak:     "MEJOR RACHA",
		StatsBestDay:    "MEJOR DÍA",
		StatsTop:        "TOP FALLOS",
		StatsHint:       "ESC VUELVE",
		NoneLabel:       "—",

		CorruptLabel: "SEÑAL DESCONOCIDA",

		MenuTerminal: "/TERMINAL",
		TermBanner:   "SISTEMA ANTIGUO v0.77 — ACCESO NO AUTORIZADO",
		TermHelpHint: "escribe 'ayuda' si te atreves",
		TermHelp: []string{
			"ayuda — esta lista",
			"archivos — documentos del sistema",
			"leer <id> — abre un documento",
			"salir — vuelve al trabajo",
			"cualquier otra cosa se interpreta como una llave",
		},
		TermUnknown:     "no reconocido. el sistema lo ha anotado igualmente.",
		TermLocked:      "BLOQUEADO. este documento requiere una llave.",
		TermLockedTag:   "[BLOQUEADO]",
		TermDocsHeader:  "DOCUMENTOS:",
		TermKeyAccepted: "LLAVE ACEPTADA. abriendo...",
		TermNoDocs:      "no hay documentos. todavía.",
		TermHint:        "ESC VUELVE",

		UpgradeNames: map[string]string{
			"mult":        "MULTIPLICADOR",
			"streakcap":   "RACHA MÁXIMA",
			"golden":      "PALABRAS DORADAS",
			"intern":      "BECARIO",
			"internspeed": "CAFÉ PARA EL BECARIO",
			"hrmult":      "BONO PERMANENTE",
			"hrday":       "JORNADA EXTENDIDA",
			"hrquota":     "CUOTA REDUCIDA",
			"hrintern":    "PLANTILLA DE BECARIOS",
			"hrend":       "COMANDO /TERMINAR",
		},
	},
	"en": {
		MenuContinue: "CONTINUE",
		MenuNew:      "NEW GAME",
		MenuLang:     "LANGUAGE",
		MenuOptions:  "OPTIONS",
		MenuReset:    "RESET",
		ResetConfirm: "SURE? ENTER WIPES EVERYTHING",
		MenuQuit:     "QUIT",

		OptDisplay:   "DISPLAY",
		OptVolume:    "SOUND",
		OptBack:      "BACK",
		DisplayModes: [3]string{"FULLSCREEN", "WINDOWED", "BORDERLESS"},

		MenuHint:    "ARROWS + ENTER",
		OptionsHint: "ARROWS ADJUST · ESC BACK",
		ShopHint:    "ENTER BUYS · ESC BACK",
		IntroHint:   "ENTER · ESC SKIPS",

		ShopTitle:  "SHOP",
		UnknownCmd: "UNKNOWN COMMAND",

		MenuHR:       "HR",
		HRTitle:      "HR OFFICE",
		HRHint:       "ENTER BUYS · ESC BACK",
		HRPointsLbl:  "HR POINTS",
		DayLabel:     "DAY",
		QuotaLabel:   "QUOTA",
		PaydayTitle:  "END OF DAY",
		FiredTitle:   "FIRED",
		NextPayLabel: "NEXT PAYROLL: DAY",
		DayEndHint:   "ENTER",
		MaxedLabel:   "MAXED",

		StreakLabel:     "STREAK",
		GoldenLabel:     "GOLDEN!",
		FrenzyLabel:     "COFFEE ×2",
		InspectionLabel: "HR INSPECTION",
		MenuStats:       "STATS",
		StatsTitle:      "STATISTICS",
		StatsWords:      "WORDS",
		StatsFails:      "FAILS",
		StatsWorst:      "WORST WORD",
		StatsStreak:     "BEST STREAK",
		StatsBestDay:    "BEST DAY",
		StatsTop:        "TOP FAILS",
		StatsHint:       "ESC BACK",
		NoneLabel:       "—",

		CorruptLabel: "UNKNOWN SIGNAL",

		MenuTerminal: "/TERMINAL",
		TermBanner:   "LEGACY SYSTEM v0.77 — UNAUTHORIZED ACCESS",
		TermHelpHint: "type 'help' if you dare",
		TermHelp: []string{
			"help — this list",
			"docs — system documents",
			"read <id> — opens a document",
			"exit — back to work",
			"anything else is treated as a key",
		},
		TermUnknown:     "not recognized. the system logged it anyway.",
		TermLocked:      "LOCKED. this document requires a key.",
		TermLockedTag:   "[LOCKED]",
		TermDocsHeader:  "DOCUMENTS:",
		TermKeyAccepted: "KEY ACCEPTED. opening...",
		TermNoDocs:      "no documents. yet.",
		TermHint:        "ESC BACK",

		UpgradeNames: map[string]string{
			"mult":        "MULTIPLIER",
			"streakcap":   "MAX STREAK",
			"golden":      "GOLDEN WORDS",
			"intern":      "INTERN",
			"internspeed": "COFFEE FOR THE INTERN",
			"hrmult":      "PERMANENT BONUS",
			"hrday":       "EXTENDED SHIFT",
			"hrquota":     "REDUCED QUOTA",
			"hrintern":    "INTERN HEADCOUNT",
			"hrend":       "/ENDDAY COMMAND",
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
		"Llegas tarde. A tu propio pasatiempo. Impresionante.",
		"Tu contraseña seguía siendo 'contraseña'. No tocamos nada.",
		"La cuota no te olvidó. Las cuotas nunca olvidan.",
		"Reiniciamos la cafetera en tu honor. Sigue rota.",
		"El departamento celebró tu ausencia. Fue una reunión breve.",
		"Volviste. El becario tenía esperanzas de ascenso. Ya no.",
		"Tu escritorio acumuló polvo corporativo. Es polvo premium.",
		"Fichaste tu entrada. El reloj te estaba juzgando.",
		"Te guardamos la silla. Te cobramos el almacenamiento.",
		"Los informes se apilaron solos. Es su magia favorita.",
		"RRHH actualizó tu expediente: 'reincidente'.",
		"Qué alegría verte. Transmite esa frase al departamento legal.",
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
		"You're late. To your own hobby. Impressive.",
		"Your password was still 'password'. We touched nothing.",
		"The quota didn't forget you. Quotas never forget.",
		"We restarted the coffee machine in your honor. Still broken.",
		"The department celebrated your absence. It was a short meeting.",
		"You're back. The intern had promotion hopes. Not anymore.",
		"Your desk gathered corporate dust. It's premium dust.",
		"You clocked in. The clock was judging you.",
		"We kept your chair. We billed you for storage.",
		"The reports piled up on their own. It's their favorite magic.",
		"HR updated your file: 'repeat offender'.",
		"So glad to see you. Forward that sentence to legal.",
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
		"Cuota abonada. El departamento de finanzas fingirá sorpresa.",
		"Día cerrado. Su esfuerzo fue registrado y luego archivado.",
		"Pago procesado. No pregunte en qué moneda.",
		"Objetivo cumplido. Mañana el objetivo lo cumple a usted.",
		"Enhorabuena. Ha pagado por el privilegio de seguir trabajando.",
		"La cuota está satisfecha. Es la única satisfecha aquí.",
		"Recibimos su dinero. Lo extrañaremos más que a usted.",
		"Balance cuadrado. Su vida, ya veremos.",
		"Su productividad fue notable. La nota fue 'suficiente'.",
		"Cuota cubierta. El jefe sonrió. Nadie sabe por qué.",
		"Un día más, una cuota menos. La matemática nunca descansa.",
		"Todo en orden. Qué sospechoso.",
	},
	"en": {
		"Your commission has been processed. You're welcome.",
		"HR thanks you for your involuntary contribution to the party fund for parties we won't throw.",
		"Quota collected. The money goes somewhere safe. Trust us.",
		"Payment received. Your loyalty has been noted in a spreadsheet.",
		"Another productive day. For us, mostly.",
		"Transaction complete. Smile, it's company policy.",
		"Quota paid. The finance department will feign surprise.",
		"Day closed. Your effort was recorded, then archived.",
		"Payment processed. Don't ask in which currency.",
		"Target met. Tomorrow the target meets you.",
		"Congratulations. You've paid for the privilege of continuing to work.",
		"The quota is satisfied. It's the only one here.",
		"We received your money. We'll miss it more than you.",
		"Books balanced. Your life, we'll see.",
		"Your productivity was remarkable. The remark was 'sufficient'.",
		"Quota covered. The boss smiled. Nobody knows why.",
		"One more day, one less quota. Math never rests.",
		"Everything in order. How suspicious.",
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
		"Recoja sus cosas. Ah, cierto, no tiene.",
		"Su acceso fue revocado antes de terminar esta frase.",
		"El algoritmo lo despidió. Nosotros solo aplaudimos.",
		"No llore sobre el teclado. Es de la empresa.",
		"Su plaza será ocupada por una hoja de cálculo.",
		"Fue un placer. Legalmente no podemos decir de quién.",
		"La puerta está donde siempre. Usted, ya no.",
		"Su rendimiento fue histórico. Historia antigua.",
		"Adjuntamos una carta de recomendación en blanco. Rellénela usted.",
		"El equipo lo despidió con aplausos. Eran de alivio.",
		"Números rojos, salida verde. Camine.",
		"Gracias por su servicio. La gratitud expira hoy.",
	},
	"en": {
		"HR wishes you success in your future endeavors. Far from here.",
		"It's not you, it's your performance. Well, it is you.",
		"Your position has been eliminated. Your chair already has a new owner.",
		"The committee reviewed your case for almost three seconds. Fired.",
		"Hand in the keyboard on your way out. The letters belong to the company.",
		"Your career here ends today. Your quota, sadly, didn't pay itself.",
		"Collect your things. Oh right, you have none.",
		"Your access was revoked before this sentence ended.",
		"The algorithm fired you. We just applauded.",
		"Don't cry on the keyboard. It's company property.",
		"Your position will be filled by a spreadsheet.",
		"It's been a pleasure. Legally we can't say whose.",
		"The door is where it's always been. You, not anymore.",
		"Your performance was historic. Ancient history.",
		"We attached a blank reference letter. Fill it in yourself.",
		"The team clapped you out. Claps of relief.",
		"Red numbers, green exit. Walk.",
		"Thank you for your service. Gratitude expires today.",
	},
}

// HRExplainScripts replaces the random payday quip on the first payroll
// day (day HRPayCycle): HR explains (in its own way) what HR points are,
// one line per Enter, intro-style. Both languages have the same count.
var HRExplainScripts = map[string][]string{
	"es": {
		"Tres días sobrevividos. Es hora de hablar de su compensación: los PUNTOS RRHH.",
		"Son una moneda inventada por el departamento para no pagarle con dinero real. No cotizan en bolsa, no compran café y no pueden heredarse.",
		"La nómina llega cada TRES días. ¿Que por qué no cada día? El papeleo, ya sabe. El papeleo.",
		"Si lo despedimos a mitad de ciclo, ese ciclo se pierde. Léalo como un incentivo para no ser despedido.",
		"Canjéelos en la OFICINA DE RRHH (comando /rrhh) y su vida laboral dolerá un poco menos.",
		"Y tranquilo: los puntos ya cobrados se los queda. Somos crueles, no ladrones.",
	},
	"en": {
		"Three days survived. Time to talk about your compensation: HR POINTS.",
		"They're a currency the department invented to avoid paying you real money. Not listed on any stock exchange, they don't buy coffee, and they can't be inherited.",
		"Payroll lands every THREE days. Why not daily? The paperwork, you know. The paperwork.",
		"If we fire you mid-cycle, that cycle is forfeited. Read it as an incentive to not get fired.",
		"Redeem them at the HR OFFICE (/hr command) and your work life will hurt a little less.",
		"And relax: points already collected are yours to keep. We're cruel, not thieves.",
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

// InternHiredScripts are HR quips shown in the shop right after hiring an
// intern. Both languages have the same count.
var InternHiredScripts = map[string][]string{
	"es": {
		"El becario no cobra. Legalmente es 'formación'.",
		"Su becario llegó. No le enseñe nada: aprende solo, como los hongos.",
		"RRHH aprobó la contratación en tiempo récord. Los gratis siempre pasan.",
		"El becario firmó su contrato sin leerlo. Nos recordó a alguien.",
		"El becario dice que estudió dos másteres. Le creímos uno.",
		"Trátelo bien. Es broma, es un becario.",
		"Viene con energía infinita y contrato finito.",
		"El becario pidió silla. Le dimos motivación.",
		"Su currículum decía 'proactivo'. Ya teclea sin preguntar.",
		"Becario incorporado. La fotocopiadora respira aliviada.",
		"Prometimos experiencia laboral. Esto cuenta, técnicamente.",
		"El becario no duerme. No preguntamos por qué.",
		"Nuevo becario: mismo sueldo que el anterior. Cero.",
		"RRHH le asignó un mentor. El mentor es usted. Suerte.",
		"Firmó por 'crecimiento profesional'. Crecerá su cuota.",
		"El becario teclea por amor al arte. El arte tampoco paga.",
	},
	"en": {
		"The intern isn't paid. Legally it's 'training'.",
		"Your intern has arrived. Don't teach them anything: they learn alone, like mushrooms.",
		"HR approved the hire in record time. Free ones always get through.",
		"The intern signed their contract without reading it. Reminded us of someone.",
		"The intern says they have two master's degrees. We believed one.",
		"Treat them well. Kidding, they're an intern.",
		"Comes with infinite energy and a finite contract.",
		"The intern asked for a chair. We gave them motivation.",
		"Their resume said 'proactive'. Already typing without asking.",
		"Intern onboarded. The photocopier breathes easier.",
		"We promised work experience. This counts, technically.",
		"The intern doesn't sleep. We don't ask why.",
		"New intern: same salary as the last one. Zero.",
		"HR assigned them a mentor. The mentor is you. Good luck.",
		"They signed for 'professional growth'. Your quota will grow.",
		"The intern types for the love of the craft. The craft doesn't pay either.",
	},
}

// EventCoffeeScripts announce a coffee frenzy. Both languages have the
// same count.
var EventCoffeeScripts = map[string][]string{
	"es": {
		"Alguien arregló la cafetera. RRHH ya está investigando.",
		"Café gratis en la sala de descanso. Corra, es un error administrativo.",
		"La oficina huele a café de verdad. Aproveche antes de que lo prohíban.",
		"El café es real. Repetimos: el café es real.",
		"Cafeína en vena. RRHH no se hace responsable de sus dedos.",
		"Hora feliz de café. La felicidad termina pronto, como todo aquí.",
		"El proveedor se equivocó de oficina. Beba rápido.",
		"Café del bueno. Presupuesto del malo. Disfrute la paradoja.",
		"Alerta: productividad espontánea detectada. Aproveche.",
		"La máquina escupe café y no facturas. Milagro trimestral.",
		"Doble ración de café. La factura llegará por triplicado.",
		"Huele a café y a oportunidad. Una de las dos es real.",
		"El jefe invitó al café. Revise que no sea descafeinado... o trampa.",
		"Café gratis: el error contable más querido del año.",
		"Suba el ritmo. El café baja solo.",
	},
	"en": {
		"Someone fixed the coffee machine. HR is already investigating.",
		"Free coffee in the break room. Run, it's an administrative error.",
		"The office smells like real coffee. Enjoy it before it gets banned.",
		"The coffee is real. We repeat: the coffee is real.",
		"Caffeine straight to the vein. HR takes no responsibility for your fingers.",
		"Coffee happy hour. The happiness ends soon, like everything here.",
		"The supplier got the wrong office. Drink fast.",
		"Good coffee. Bad budget. Enjoy the paradox.",
		"Alert: spontaneous productivity detected. Make it count.",
		"The machine is dispensing coffee, not invoices. Quarterly miracle.",
		"Double coffee ration. The bill will arrive in triplicate.",
		"Smells like coffee and opportunity. One of the two is real.",
		"The boss bought the coffee. Check it's not decaf... or a trap.",
		"Free coffee: the most beloved accounting error of the year.",
		"Pick up the pace. The coffee goes down on its own.",
	},
}

// EventInspectionScripts announce an HR inspection (the next word becomes
// urgent). Both languages have the same count.
var EventInspectionScripts = map[string][]string{
	"es": {
		"RRHH pasa por su mesa. Escriba esto. Rápido. Sonría.",
		"Inspección sorpresa. Sorpresa para usted; nosotros lo sabíamos hace semanas.",
		"Auditoría de productividad en curso. Demuestre que existe.",
		"Los de arriba están mirando. Literalmente, desde arriba.",
		"Control de calidad humano. Usted es el humano en duda.",
		"Escriba como si su contrato dependiera de ello. Depende.",
		"RRHH trae portapapeles. Nunca es buena señal.",
		"Demuestre productividad. El cronómetro ya corre.",
		"Visita sorpresa del comité. Sonría con los dedos.",
		"Esta palabra va directa a su expediente. Sin presión.",
		"Auditoría relámpago. Los relámpagos caen dos veces aquí.",
		"El inspector no parpadea. Usted tampoco debería.",
		"Palabra urgente, salario constante. Así funciona.",
		"Su evaluación anual se decide en segundos. Estos.",
		"Teclee natural. Como si no lo estuviéramos grabando.",
	},
	"en": {
		"HR is walking by your desk. Type this. Fast. Smile.",
		"Surprise inspection. Surprise for you; we knew weeks ago.",
		"Productivity audit in progress. Prove you exist.",
		"Upper management is watching. Literally, from above.",
		"Human quality control. You're the human in question.",
		"Type like your contract depends on it. It does.",
		"HR brought a clipboard. Never a good sign.",
		"Demonstrate productivity. The stopwatch is already running.",
		"Surprise committee visit. Smile with your fingers.",
		"This word goes straight to your file. No pressure.",
		"Lightning audit. Lightning strikes twice here.",
		"The inspector doesn't blink. Neither should you.",
		"Urgent word, constant salary. That's how it works.",
		"Your annual review is decided in seconds. These ones.",
		"Type naturally. As if we weren't recording you.",
	},
}

// UpgradeBoughtScripts are HR quips shown right after buying an upgrade,
// keyed by upgrade ID then language. Every pool has the same line count in
// both languages. The intern purchase uses InternHiredScripts instead.
var UpgradeBoughtScripts = map[string]map[string][]string{
	"mult": {
		"es": {
			"Multiplicador instalado. Su dedo vale más. Usted, igual.",
			"Cada palabra paga más. La inflación le manda saludos.",
			"Mejora aprobada. Finanzas lloró un poco.",
			"Ahora escribe oro. Sigue sin poder gastarlo en café.",
			"Su tarifa subió. Su categoría, no.",
			"Palabras premium. Mismo teclado de plástico.",
			"El contador de puntos pidió un aumento. Denegado.",
			"Más puntos por palabra. La cuota ya tomó nota.",
			"Rendimiento amplificado. Expectativas también.",
			"Compró productividad. El secreto mejor guardado del capitalismo.",
			"Sus palabras valen doble. Sus opiniones siguen valiendo cero.",
			"Multiplicador activo. No lo multiplique a usted, por favor.",
			"Inversión sensata. Primera vez que escribimos eso.",
			"El sistema ahora lo sobrevalora oficialmente.",
			"Más ingresos por tecla. El teclado exige comisión.",
			"Mejora registrada. Alguien en contabilidad se desmayó.",
			"Su producción vale más. Enhorabuena, activo empresarial.",
			"Puntos extra por el mismo esfuerzo. No se acostumbre.",
			"La empresa paga más por palabra. Recuperaremos el dinero en cuota.",
			"Multiplicador nuevo. La motivación sigue de baja.",
		},
		"en": {
			"Multiplier installed. Your finger is worth more. You, the same.",
			"Every word pays more. Inflation sends its regards.",
			"Upgrade approved. Finance cried a little.",
			"You type gold now. Still can't spend it on coffee.",
			"Your rate went up. Your job title didn't.",
			"Premium words. Same plastic keyboard.",
			"The score counter asked for a raise. Denied.",
			"More points per word. The quota already took note.",
			"Output amplified. Expectations too.",
			"You bought productivity. Capitalism's best-kept secret.",
			"Your words are worth double. Your opinions still worth zero.",
			"Multiplier active. Please don't multiply yourself.",
			"A sensible investment. First time we've ever written that.",
			"The system now officially overvalues you.",
			"More income per keystroke. The keyboard demands commission.",
			"Upgrade recorded. Someone in accounting fainted.",
			"Your output is worth more. Congratulations, corporate asset.",
			"Extra points for the same effort. Don't get used to it.",
			"The company pays more per word. We'll claw it back in quota.",
			"New multiplier. Motivation remains on sick leave.",
		},
	},
	"streakcap": {
		"es": {
			"Techo de racha elevado. Su constancia por fin tiene permiso.",
			"Ahora puede encadenar más aciertos. Encadenado, como nos gusta.",
			"Racha ampliada. El error sigue esperando su momento.",
			"Más combo, más presión. De nada.",
			"Su racha puede crecer. Su ego, contrólelo.",
			"Límite subido. Los límites existen para venderlos.",
			"Consistencia certificada. El certificado no es válido fuera.",
			"El contador de racha respiraba con dificultad. Lo ampliamos.",
			"Más margen para lucirse. Nadie está mirando.",
			"Racha máxima nueva. La usará una vez y fallará.",
			"Combo extendido. Extienda también su jornada, ya que está.",
			"Su cadena de aciertos tiene eslabones nuevos. Sin sindicato.",
			"Permiso para ser bueno más rato. Revocable.",
			"El multiplicador de racha agradece el espacio extra.",
			"Techo nuevo. El suelo sigue siendo el mismo.",
			"Ahora sus rachas caben. Sus sueños, aún no.",
			"Ampliación aprobada. RRHH apostó a que la desperdicia.",
			"Más racha, más puntos, más cuota. El ciclo de la vida.",
			"Suba esa racha. La gravedad corporativa hará el resto.",
			"Límite de racha extendido. Todo lo demás sigue limitado.",
		},
		"en": {
			"Streak ceiling raised. Your consistency finally has a permit.",
			"You can chain more hits now. Chained, just how we like it.",
			"Streak extended. The mistake is still waiting for its moment.",
			"More combo, more pressure. You're welcome.",
			"Your streak can grow. Your ego, keep it in check.",
			"Limit raised. Limits exist to be sold.",
			"Certified consistency. The certificate isn't valid outside.",
			"The streak counter was struggling to breathe. We expanded it.",
			"More room to shine. Nobody's watching.",
			"New max streak. You'll use it once and miss.",
			"Combo extended. Extend your shift too, while you're at it.",
			"Your chain of hits has new links. No union.",
			"Permission to be good for longer. Revocable.",
			"The streak multiplier appreciates the extra space.",
			"New ceiling. The floor is still the same.",
			"Now your streaks fit. Your dreams, not yet.",
			"Extension approved. HR bet you'll waste it.",
			"More streak, more points, more quota. The circle of life.",
			"Raise that streak. Corporate gravity will handle the rest.",
			"Streak limit extended. Everything else remains limited.",
		},
	},
	"golden": {
		"es": {
			"Palabras doradas mejoradas. El brillo no es oro. Bueno, sí.",
			"Más suerte comprada. La suerte también factura.",
			"El oro aparece más seguido. Hacienda ha sido notificada.",
			"Brillo ampliado. Gafas de sol no incluidas.",
			"Ahora el azar trabaja para usted. Sin contrato, como todos.",
			"Doradas más frecuentes. El becario pidió una. Denegado.",
			"Compró probabilidad. El casino corporativo agradece.",
			"Más palabras doradas. Siguen sin ser comestibles.",
			"La fortuna sonríe más. Es una sonrisa corporativa: vacía.",
			"Oro mejorado. El rey Midas era de RRHH, por cierto.",
			"Aumentamos su suerte. Del rendimiento se encarga usted.",
			"Palabras que brillan más y pagan más. Sospechosamente generoso.",
			"El destello dorado ya no es un error de pantalla.",
			"Estadísticamente más rico. Estadísticamente.",
			"Mejor chance, mejor pago. La cuota observa con interés.",
			"El oro llama al oro. La cuota llama más fuerte.",
			"Doradas mejoradas. Recuerde: el brillo se descuenta del alma.",
			"Suerte al alza. Riesgo de acostumbrarse: alto.",
			"Más oro en pantalla. El de la caja fuerte sigue siendo nuestro.",
			"Nivel dorado subido. Su aura ya casi refleja.",
		},
		"en": {
			"Golden words upgraded. All that glitters isn't gold. Well, it is.",
			"More luck purchased. Luck invoices too.",
			"Gold shows up more often. The tax office has been notified.",
			"Shine expanded. Sunglasses not included.",
			"Chance now works for you. No contract, like everyone.",
			"More frequent goldens. The intern asked for one. Denied.",
			"You bought probability. The corporate casino thanks you.",
			"More golden words. Still not edible.",
			"Fortune smiles more. It's a corporate smile: empty.",
			"Gold improved. King Midas worked in HR, by the way.",
			"We raised your luck. Performance is on you.",
			"Words that shine more and pay more. Suspiciously generous.",
			"The golden flash is no longer a display error.",
			"Statistically richer. Statistically.",
			"Better odds, better pay. The quota watches with interest.",
			"Gold calls to gold. The quota calls louder.",
			"Goldens upgraded. Remember: the shine is deducted from your soul.",
			"Luck trending up. Risk of getting used to it: high.",
			"More gold on screen. The one in the vault is still ours.",
			"Golden level up. Your aura is almost reflective now.",
		},
	},
	"internspeed": {
		"es": {
			"Café para el becario. Ahora vibra a 60Hz.",
			"Dosis de cafeína administrada. El becario ve sonidos.",
			"El becario teclea más rápido. Parpadear es opcional.",
			"Café asignado. Su sueño fue reasignado.",
			"Velocidad de becario aumentada. La normativa no aplica aquí.",
			"Más café, más letras. La ciencia corporativa es simple.",
			"El becario alcanzó una nueva marcha. No tiene frenos.",
			"Cafeína de grado industrial. Uso humano: discutible.",
			"El becario ya no camina, se desplaza. Escribiendo.",
			"Café premium para el becario. Lo barato salía caro.",
			"Sus manos son un borrón. Auditado y aprobado.",
			"Turbo activado. El becario pidió agua. Le dimos más café.",
			"El corazón del becario late en 4/4. Compás de trabajo.",
			"Nueva taza asignada. La taza es métrica de rendimiento.",
			"El becario descubrió la velocidad. La velocidad lo descubrió a él.",
			"Café mejorado. El temblor es una función, no un error.",
			"Teclea tan rápido que RRHH le pidió calmarse. Es broma: más rápido.",
			"Ración doble. El becario nos escribió 'gracias' en 0.2 segundos.",
			"Los dedos del becario alcanzaron velocidad crucero. Sin aterrizaje previsto.",
			"Inversión en cafeína. El retorno se mide en palabras por minuto.",
		},
		"en": {
			"Coffee for the intern. They now vibrate at 60Hz.",
			"Caffeine dose administered. The intern can see sounds.",
			"The intern types faster. Blinking is optional.",
			"Coffee assigned. Their sleep was reassigned.",
			"Intern speed increased. Regulations don't apply here.",
			"More coffee, more letters. Corporate science is simple.",
			"The intern found a new gear. There are no brakes.",
			"Industrial-grade caffeine. Human use: debatable.",
			"The intern no longer walks, they glide. While typing.",
			"Premium coffee for the intern. Cheap was getting expensive.",
			"Their hands are a blur. Audited and approved.",
			"Turbo engaged. The intern asked for water. We gave them more coffee.",
			"The intern's heart beats in 4/4. Work tempo.",
			"New mug assigned. The mug is a performance metric.",
			"The intern discovered speed. Speed discovered them back.",
			"Coffee upgraded. The trembling is a feature, not a bug.",
			"They type so fast HR asked them to slow down. Kidding: faster.",
			"Double ration. The intern typed us 'thanks' in 0.2 seconds.",
			"The intern's fingers reached cruising speed. No landing scheduled.",
			"Caffeine investment. Returns measured in words per minute.",
		},
	},
	"hrmult": {
		"es": {
			"Bono permanente. Permanente como su deuda con la empresa.",
			"Aumento vitalicio. La vida laboral, se entiende.",
			"RRHH aprobó pagarle más. Guarde este momento histórico.",
			"Bonificación eterna. La eternidad tiene letra pequeña.",
			"Su valor de mercado subió. El mercado somos nosotros.",
			"Mejora permanente adquirida. Ni despedirlo la borra. Qué error.",
			"Puntos mejor pagados para siempre. 'Siempre' según contrato.",
			"Este bono lo acompañará entre despidos. Qué romántico.",
			"Sueldo base mejorado. La base sigue bajo tierra.",
			"Un lujo permanente. Su primera posesión real.",
			"El bono sobrevive a todo. Cucarachas y bonos.",
			"Firmado y sellado: vale usted un 10% más. Enmarcable.",
			"Mejora indeleble. Como las manchas de café del contrato.",
			"RRHH invirtió en usted. Todos en la oficina perdieron una apuesta.",
			"Bono perpetuo activado. Perpetuo suena caro. Lo fue.",
			"Su multiplicador trasciende despidos. Usted no.",
			"Beneficio de por vida. Revisamos: sí, es legal.",
			"Más puntos por siempre jamás. Fin del cuento.",
			"El bono es suyo. Lo único aquí que lo es.",
			"Mejora permanente registrada en un papel que no perderemos. Probablemente.",
		},
		"en": {
			"Permanent bonus. Permanent like your debt to the company.",
			"A lifetime raise. Work life, that is.",
			"HR approved paying you more. Treasure this historic moment.",
			"Eternal bonus. Eternity has fine print.",
			"Your market value went up. We are the market.",
			"Permanent upgrade acquired. Not even firing you erases it. What a mistake.",
			"Better-paid points forever. 'Forever' as per contract.",
			"This bonus will follow you between firings. How romantic.",
			"Base salary improved. The base remains underground.",
			"A permanent luxury. Your first real possession.",
			"The bonus survives everything. Cockroaches and bonuses.",
			"Signed and sealed: you're worth 10% more. Frame it.",
			"Indelible upgrade. Like the coffee stains on your contract.",
			"HR invested in you. Everyone in the office lost a bet.",
			"Perpetual bonus active. Perpetual sounds expensive. It was.",
			"Your multiplier transcends firings. You don't.",
			"A benefit for life. We checked: yes, it's legal.",
			"More points forever and ever. The end.",
			"The bonus is yours. The only thing here that is.",
			"Permanent upgrade recorded on a paper we won't lose. Probably.",
		},
	},
	"hrday": {
		"es": {
			"Jornada extendida. Celebre: más tiempo aquí dentro.",
			"Compró más horas de trabajo. Léalo otra vez. Con calma.",
			"El día ahora dura más. El sol no fue consultado.",
			"Más segundos por jornada. Cada uno cuenta. Nosotros los contamos.",
			"Tiempo extra aprobado. El tiempo libre presentó su renuncia.",
			"Su día laboral creció. Su día personal fue absorbido.",
			"Jornada más larga. El reloj pidió horas extra. Irónico.",
			"Extendimos su día. Newton se revuelve en su tumba.",
			"Más tiempo para teclear. Los sueños pueden esperar.",
			"El calendario protestó. Lo despedimos también.",
			"Día ampliado. La noche presentó una queja formal.",
			"Compró tiempo. El truco más viejo de la empresa, al revés.",
			"Más minutos de productividad. La palabra favorita de alguien.",
			"Jornada premium: ahora con más jornada.",
			"El reloj corre más lento para usted. De miedo, creemos.",
			"Tiempo adicional concedido. Adminístrelo mal, como siempre.",
			"Su turno se estiró. Usted también debería: ergonomía.",
			"Más horas frente a la pantalla. La pantalla está encantada.",
			"El día rinde más. Usted decide qué significa 'rendir'.",
			"Tiempo extendido. Einstein lo llamó relatividad. Nosotros, negocio.",
		},
		"en": {
			"Extended shift. Celebrate: more time in here.",
			"You bought more working hours. Read that again. Slowly.",
			"The day lasts longer now. The sun was not consulted.",
			"More seconds per shift. Each one counts. We count them.",
			"Overtime approved. Free time handed in its resignation.",
			"Your work day grew. Your personal day was absorbed.",
			"Longer shift. The clock asked for overtime. Ironic.",
			"We extended your day. Newton spins in his grave.",
			"More time to type. Dreams can wait.",
			"The calendar protested. We fired it too.",
			"Day expanded. The night filed a formal complaint.",
			"You bought time. The company's oldest trick, in reverse.",
			"More minutes of productivity. Somebody's favorite word.",
			"Premium shift: now with more shift.",
			"The clock runs slower for you. Out of fear, we believe.",
			"Additional time granted. Manage it poorly, as usual.",
			"Your shift stretched. You should too: ergonomics.",
			"More hours at the screen. The screen is delighted.",
			"The day yields more. You decide what 'yields' means.",
			"Time extended. Einstein called it relativity. We call it business.",
		},
	},
	"hrquota": {
		"es": {
			"Cuota reducida. RRHH exige menos. Sospeche.",
			"Expectativas rebajadas oficialmente. Un sueño cumplido.",
			"La cuota adelgazó. La primera dieta que funciona aquí.",
			"Menos exigencia. El listón bajó y nadie se golpeó.",
			"Cuota recortada. Finanzas está en duelo.",
			"Pedimos menos. No se lo cuente a los demás departamentos.",
			"Rebaja aprobada. El formulario tenía 47 páginas. Las firmó alguien.",
			"La cuota perdió peso. Usted ganó margen.",
			"Menos números rojos en su horizonte. Siguen ahí, pero menos.",
			"Exigencia reducida. La pereza institucional lo celebra.",
			"Su cuota encogió en la lavadora corporativa.",
			"El monstruo de fin de día ahora muerde más suave.",
			"Descuento en obligaciones. Sin devoluciones.",
			"Cuota más humana. Todo lo humana que puede ser una cuota.",
			"Menos presión diaria. La presión mensual no existe. Aún.",
			"El objetivo bajó un escalón. No baje usted dos.",
			"Cuota reducida de por vida. La vida laboral, insistimos.",
			"RRHH afloja la soga. Sigue siendo una soga.",
			"Pagará menos cada día. Alguien arriba perdió esa negociación.",
			"La cuota cede terreno. Territorio reconquistado, mecanógrafo.",
		},
		"en": {
			"Reduced quota. HR demands less. Be suspicious.",
			"Expectations officially lowered. A dream come true.",
			"The quota slimmed down. The first diet that works here.",
			"Less demanded. The bar dropped and nobody got hit.",
			"Quota cut. Finance is in mourning.",
			"We ask for less. Don't tell the other departments.",
			"Discount approved. The form had 47 pages. Somebody signed them.",
			"The quota lost weight. You gained margin.",
			"Fewer red numbers on your horizon. Still there, just fewer.",
			"Requirements reduced. Institutional laziness celebrates.",
			"Your quota shrank in the corporate washing machine.",
			"The end-of-day monster bites softer now.",
			"A discount on obligations. No refunds.",
			"A more humane quota. As humane as a quota can be.",
			"Less daily pressure. Monthly pressure doesn't exist. Yet.",
			"The target stepped down one rung. Don't you step down two.",
			"Quota reduced for life. Work life, we insist.",
			"HR loosens the rope. It's still a rope.",
			"You'll pay less each day. Someone upstairs lost that negotiation.",
			"The quota gives ground. Territory reclaimed, typist.",
		},
	},
	"hrintern": {
		"es": {
			"Plantilla ampliada. Más becarios, mismo presupuesto: cero.",
			"Cupo de becarios subido. La fotocopiadora tiembla.",
			"Autorización firmada: puede coleccionar becarios.",
			"Más manos gratis. Legal dice que no lo digamos así.",
			"El aula de becarios ahora tiene más pupitres.",
			"RRHH aprueba más becarios. El comité ético estaba de vacaciones.",
			"Ampliación de plantilla no remunerada. El futuro del empleo.",
			"Más becarios por partida. Coleccione los cuatro.",
			"Su ejército de formación crece. Sin armas, con teclados.",
			"Nuevo cupo aprobado. Los becarios vienen en packs, como el folio.",
			"Puede contratar más becarios. Ellos no pueden negarse. Simetría.",
			"Plantilla extendida. La palabra 'plantilla' hace mucho trabajo ahí.",
			"Más becarios disponibles. La cantera no se agota. Lo comprobamos.",
			"Cupo elevado. Cada becario incluye su propia taza. Es mentira.",
			"El límite de becarios subió. El límite moral sigue en revisión.",
			"Más sillas junto a la suya. Sillas metafóricas.",
			"Becarios adicionales autorizados. El organigrama ya no cabe en la pizarra.",
			"Su equipo puede crecer. La palabra 'equipo' es generosa.",
			"Contrate más aprendices. Aprenderán lo que usted: a teclear.",
			"Cupo nuevo de becarios. RRHH los cría, usted los explota. Sinergia.",
		},
		"en": {
			"Headcount expanded. More interns, same budget: zero.",
			"Intern cap raised. The photocopier trembles.",
			"Authorization signed: you may collect interns.",
			"More free hands. Legal says don't phrase it that way.",
			"The intern classroom has more desks now.",
			"HR approves more interns. The ethics committee was on vacation.",
			"Unpaid headcount expansion. The future of employment.",
			"More interns per run. Collect all four.",
			"Your training army grows. No weapons, just keyboards.",
			"New cap approved. Interns come in packs, like paper.",
			"You can hire more interns. They can't refuse. Symmetry.",
			"Headcount extended. The word 'headcount' does a lot of work there.",
			"More interns available. The talent pool never runs dry. We checked.",
			"Cap raised. Each intern includes their own mug. That's a lie.",
			"The intern limit went up. The moral limit is still under review.",
			"More chairs next to yours. Metaphorical chairs.",
			"Additional interns authorized. The org chart no longer fits the whiteboard.",
			"Your team may grow. The word 'team' is generous.",
			"Hire more trainees. They'll learn what you did: typing.",
			"New intern cap. HR raises them, you exploit them. Synergy.",
		},
	},
	"hrend": {
		"es": {
			"Comando /terminar desbloqueado. El poder de rendirse con estilo.",
			"Ahora puede cerrar el día antes. La cuota no espera: cobra igual.",
			"Botón de escape adquirido. Escape hacia otra cuota.",
			"El día termina cuando usted diga. Frase nunca antes oída aquí.",
			"Poder ejecutivo concedido. Ejecútese con cuidado.",
			"/terminar activo. Úselo con confianza y saldo suficiente.",
			"Puede irse antes. Los valientes cobran, los cortos... adiós.",
			"Autonomía comprada. La vendemos carísima por algo.",
			"El reloj ahora le obedece. Qué envidia.",
			"Cierre anticipado habilitado. RRHH cobra igual: siempre cobra.",
			"Comando premium instalado. Escríbalo sin errores, sería irónico.",
			"Terminar el día es un privilegio. Ya es oficialmente privilegiado.",
			"La jornada ya no lo retiene. La nómina sí.",
			"Salida rápida desbloqueada. Rápida no significa gratuita.",
			"Ahora decide cuándo rendir cuentas. Las cuentas no cambian.",
			"El interruptor del día es suyo. La factura de la luz también.",
			"Fin de jornada a la carta. El menú tiene un solo plato: cuota.",
			"Adelante, termine el día. La cuota estará esperando en la puerta.",
			"Control del tiempo adquirido. Úselo mejor que nosotros.",
			"/terminar es suyo. Recuerde: rendirse temprano también se factura.",
		},
		"en": {
			"/endday command unlocked. The power to give up in style.",
			"You can close the day early now. The quota doesn't wait: it charges anyway.",
			"Escape button acquired. Escape into another quota.",
			"The day ends when you say. A sentence never heard here before.",
			"Executive power granted. Execute with care.",
			"/endday active. Use it with confidence and sufficient balance.",
			"You may leave early. The brave cash out, the short... goodbye.",
			"Autonomy purchased. There's a reason we price it so high.",
			"The clock obeys you now. How enviable.",
			"Early closure enabled. HR charges anyway: it always charges.",
			"Premium command installed. Type it without typos, that'd be ironic.",
			"Ending the day is a privilege. You're now officially privileged.",
			"The shift no longer holds you. The payroll does.",
			"Fast exit unlocked. Fast doesn't mean free.",
			"You now decide when to settle accounts. The accounts don't change.",
			"The day's off-switch is yours. So is the electricity bill.",
			"End-of-shift à la carte. The menu has one dish: quota.",
			"Go ahead, end the day. The quota will be waiting at the door.",
			"Time control acquired. Use it better than we do.",
			"/endday is yours. Remember: giving up early is also billable.",
		},
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
