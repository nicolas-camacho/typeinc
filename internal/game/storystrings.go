package game

// This file holds every piece of story content for the background mystery:
// glitched HR quips, corrupt words, terminal documents, marked words and
// the small toast pools. Same rules as strings.go: every player-facing
// pool exists in "es" and "en" with identical line counts (tests enforce
// it), and every key/ID the player can type is a language-neutral token.
//
// Canon: below the building lives "the guest" (employee 0, arrived 1997) —
// an extraterrestrial entity the company feeds with every typed word so it
// can learn human language and, one day, walk upstairs unnoticed. Vera,
// the previous typist, found out; now she types down there ON PURPOSE,
// rationing what it learns. HR are human accomplices paid to look away.

// StoryQuipScripts are the glitched HR lines interleaved among normal
// day-start and day-end quips, one pool per act (StoryStage 1..6) plus a
// tiny epilogue pool under StoryEndAct. They share one no-repeat cursor.
// The indices listed in storyHintLines (story.go) are objective hints and
// get picked with priority while that act's key is still unused.
var StoryQuipScripts = map[int]map[string][]string{
	1: {
		"es": {
			"Buenos días, Vera. ... Perdón. Costumbre. Tú no eres Vera. Todavía.",
			"Recuerda sonreír. Las cámaras funcionan mejor cuando",
			"Tu puesto fue ocupado por última vez hace [DATO NO DISPONIBLE].",
			"Aquí nunca ha trabajado nadie más. Deja de preguntar. Nadie ha preguntado. Bien.",
			"El teclado tiene comandos que RRHH no instaló. Por ejemplo /terminal. No lo escribas.",
			"Hay notas manuscritas en el sistema viejo. RRHH no las escribió. RRHH no las ha borrado. Curioso.",
			"Si una palabra morada que no pediste aparece en tu pantalla, NO la teclees. Es una prueba. Tecléala.",
			"El nombre que buscas empieza por V. Este mensaje no existe y no lo has leído.",
		},
		"en": {
			"Good morning, Vera. ... Sorry. Habit. You are not Vera. Yet.",
			"Remember to smile. The cameras work better when",
			"Your desk was last occupied [DATA NOT AVAILABLE] ago.",
			"Nobody else has ever worked here. Stop asking. Nobody asked. Good.",
			"Your keyboard has commands HR never installed. For example /terminal. Do not type it.",
			"There are handwritten notes in the old system. HR didn't write them. HR hasn't deleted them. Curious.",
			"If a purple word you didn't ask for shows up on your screen, do NOT type it. It's a test. Type it.",
			"The name you're looking for starts with V. This message does not exist and you have not read it.",
		},
	},
	2: {
		"es": {
			"El expediente que abriste no existe. El nombre que leíste tampoco. Vuelve al trabajo.",
			"Vera cumplía su cuota. Vera cumplía TODAS sus cuotas. Eso no la ayudó.",
			"RRHH desea aclarar que nadie ha desaparecido jamás en esta empresa. Reubicada. La palabra es reubicada.",
			"Su último turno fue el 77. El turno 77 no existe en el calendario. Escríbelo todo junto y calla.",
			"Los registros de Vera siguen en el sistema viejo. La llave es su turno. RRHH no ha dicho esto.",
			"A veces el edificio recuerda a sus empleados. Es un defecto estructural. Ignóralo.",
			"Tu productividad es idéntica a la del ocupante anterior. Idéntica. Tecla por tecla.",
			"No hay sótano en este edificio. El ascensor tiene un botón de menos. Cuenta otra vez.",
		},
		"en": {
			"The file you opened does not exist. Neither does the name you read. Back to work.",
			"Vera made her quota. Vera made ALL her quotas. It didn't help her.",
			"HR wishes to clarify that nobody has ever disappeared at this company. Relocated. The word is relocated.",
			"Her last shift was 77. Shift 77 does not exist on the calendar. Write it as one word and keep quiet.",
			"Vera's logs are still in the old system. The key is her shift. HR did not say this.",
			"Sometimes the building remembers its employees. It's a structural defect. Ignore it.",
			"Your productivity is identical to the previous occupant's. Identical. Key by key.",
			"There is no basement in this building. The elevator has one button too few. Count again.",
		},
	},
	3: {
		"es": {
			"Las palabras que tecleas no se guardan. Se envían abajo. Algo se las come.",
			"Lo de abajo no nació aquí. Lo de abajo no nació en ningún sitio de aquí.",
			"sombra. puerta. abajo. silencio. Son sus favoritas. Repítelas y verás. NO las repitas.",
			"Vera contaba ciertas palabras. Tres veces cada una. Fue suficiente para que algo contestara.",
			"¿A dónde crees que va todo lo que escribes? Correcto. Al sotano. Esa palabra abre cosas.",
			"Los logs de Vera mencionan un nivel que no aparece en los planos. Los planos fueron corregidos.",
			"Escucha. ¿Oyes teclear? Tú no estás tecleando ahora mismo. Ella sí.",
			"El sistema registra 2 mecanógrafos activos y 1 estudiante de idiomas. No preguntes qué estudia.",
		},
		"en": {
			"The words you type are not stored. They are sent below. Something eats them.",
			"The thing below was not born here. The thing below was not born anywhere near here.",
			"shadow. door. below. silence. Those are its favorites. Repeat them and see. Do NOT repeat them.",
			"Vera counted certain words. Three times each. It was enough to make something answer.",
			"Where do you think everything you type goes? Correct. To the sotano. That word opens things.",
			"Vera's logs mention a level that is not on the blueprints. The blueprints were corrected.",
			"Listen. Do you hear typing? You are not typing right now. She is.",
			"The system registers 2 active typists and 1 language student. Do not ask what it studies.",
		},
	},
	4: {
		"es": {
			"Los becarios están bien. Los becarios están perfectamente. Deja de mirar a los becarios.",
			"Uno de tus becarios escribió una frase que nadie le asignó. Abajo la oyeron. Déjale terminar la próxima vez.",
			"Las cámaras del pasillo B apuntan todas a la misma esquina. Seguridad archivó un parte. Léelo.",
			"Tu reloj y el reloj del edificio discrepan 47 minutos diarios. El huésped también estudia las horas.",
			"Ayer trabajaste 60 segundos más de los que existieron. Alguien de abajo los está aprendiendo.",
			"Si tu reloj muestra una hora imposible, no es un error. Es la hora del sitio del que vino.",
			"Las señales moradas son su forma de practicar. Cada una que tecleas le enseña. Cada una que tecleas te enseña.",
			"El becario nuevo preguntó por Vera. No hay becario nuevo. Había. Preguntó demasiado bien.",
		},
		"en": {
			"The interns are fine. The interns are perfectly fine. Stop looking at the interns.",
			"One of your interns typed a phrase nobody assigned. Below, they heard it. Let them finish next time.",
			"The hallway B cameras all point at the same corner. Security filed a report. Read it.",
			"Your clock and the building's clock disagree by 47 minutes a day. The guest studies hours too.",
			"Yesterday you worked 60 seconds more than existed. Someone below is learning them.",
			"If your clock shows an impossible time, it is not an error. It is the local time of where it came from.",
			"The purple signals are how it practices. Every one you type teaches it. Every one you type teaches you.",
			"The new intern asked about Vera. There is no new intern. There was. They asked too well.",
		},
	},
	5: {
		"es": {
			"La nómina incluye un pago al empleado 0 desde 1997. El año en que ALGO fichó por la empresa.",
			"El empleado 0 no cobra en dinero. Cobra en palabras. Está aprendiendo a pedirlas él mismo.",
			"Ayer el huésped dijo su primera frase completa: 'siguiente palabra por favor'. Contabilidad lloró.",
			"Cuatro. Cuatro. Uno. Tres. Está escrito en tus cierres de día. RRHH no ha dicho nada. RRHH nunca dice nada.",
			"Mira la cuota de esta noche. Los dígitos morados no son un error de pantalla. Apúntalos. Escríbelos en el sistema viejo.",
			"El acceso al nivel inferior requiere un código que ningún empleado vivo conoce. Enhorabuena por tu iniciativa.",
			"Vera encontró el código. Lo tecleó. Los registros dicen que no salió. Los registros no dicen que quisiera salir.",
			"Cuando hable como nosotros, subirá. Por ahora confunde 'nómina' con 'hambre'. Comprensible.",
		},
		"en": {
			"Payroll includes a payment to employee 0 since 1997. The year SOMETHING joined the company.",
			"Employee 0 is not paid in money. It is paid in words. It is learning to ask for them itself.",
			"Yesterday the guest said its first full sentence: 'next word please'. Accounting wept.",
			"Four. Four. One. Three. It is written on your day-end reports. HR has said nothing. HR never says anything.",
			"Look at tonight's quota. The purple digits are not a display error. Write them down. Type them into the old system.",
			"Access to the lower level requires a code no living employee knows. Congratulations on your initiative.",
			"Vera found the code. She typed it. The records say she never came out. The records don't say she wanted to.",
			"When it speaks like us, it will come upstairs. For now it confuses 'payroll' with 'hunger'. Understandable.",
		},
	},
	6: {
		"es": {
			"La puerta de abajo tiene un teclado. Acepta una sola combinación: el nombre y el número. Todo junto. Sin espacios.",
			"Ya sabes el nombre de ella. Ya viste el número en tus cierres. Deja de fingir que no.",
			"Detrás de la puerta hay un estudiante de idiomas de otro mundo. Lleva 27 años de clases. Casi no tiene acento.",
			"Quien abre la puerta no es despedido. Queremos que quede claro. No es exactamente despedido.",
			"El plan del huésped es simple: aprender a hablar, subir, encajar. Lleva décadas practicando con tus palabras.",
			"Vera teclea abajo por una razón. Alguien tiene que elegir QUÉ aprende. Prefiere que sea ella.",
			"Cuando la abras, saluda. Le gusta practicar los saludos. Los saludos casi le salen.",
			"Última advertencia del departamento: no abras la puerta. Es la advertencia estándar. Nadie la ha escuchado nunca.",
		},
		"en": {
			"The door below has a keypad. It takes a single combination: the name and the number. All together. No spaces.",
			"You already know her name. You already saw the number on your reports. Stop pretending you don't.",
			"Behind the door there is a language student from another world. 27 years of lessons. Almost no accent left.",
			"Whoever opens the door is not fired. We want to be clear. Not exactly fired.",
			"The guest's plan is simple: learn to speak, go upstairs, blend in. It has practiced on your words for decades.",
			"Vera types down there for a reason. Someone has to choose WHAT it learns. She'd rather it be her.",
			"When you open it, say hello. It likes practicing greetings. It almost gets greetings right.",
			"Final warning from the department: do not open the door. It's the standard warning. Nobody has ever listened.",
		},
	},
	StoryEndAct: {
		"es": {
			"Buenos días. Solo buenos días. Hoy el edificio está callado. Abajo, alguien dosifica sus lecciones.",
			"El registro tiene una entrada nueva: GRACIAS POR LAS PALABRAS. La ortografía es perfecta. Eso preocupa.",
			"Tu cuota de hoy es la de siempre. El huésped sigue abajo. Su acento mejora. Vera lo vigila.",
			"RRHH ha archivado el caso Vera. El expediente termina con una frase que nadie tecleó arriba: 'aun no'.",
		},
		"en": {
			"Good morning. Just good morning. The building is quiet today. Below, someone rations its lessons.",
			"The log has a new entry: THANK YOU FOR THE WORDS. The spelling is perfect. That is worrying.",
			"Today's quota is the usual one. The guest is still below. Its accent improves. Vera keeps watch.",
			"HR has archived the Vera case. The file ends with a sentence nobody typed upstairs: 'not yet'.",
		},
	},
}

// CorruptWord is one non-dictionary word the building makes the player
// type. Texts are lowercase ASCII (a-z and spaces) because TypeChar
// matches byte by byte; the ClueID is language-neutral so a clue found in
// one language counts in the other.
type CorruptWord struct {
	ClueID string
	Act    int // minimum StoryStage for this word to appear
	Text   map[string]string
}

// CorruptWords lists every corrupt word, roughly two per act.
var CorruptWords = []CorruptWord{
	{ClueID: "cw_help", Act: 1, Text: map[string]string{"es": "ayudame", "en": "help me"}},
	{ClueID: "cw_name", Act: 1, Text: map[string]string{"es": "vera", "en": "vera"}},
	{ClueID: "cw_notme", Act: 2, Text: map[string]string{"es": "no fui yo", "en": "not me"}},
	{ClueID: "cw_shift", Act: 2, Text: map[string]string{"es": "turno setenta y siete", "en": "shift seventy seven"}},
	{ClueID: "cw_listen", Act: 3, Text: map[string]string{"es": "nos escuchan", "en": "they listen"}},
	{ClueID: "cw_below", Act: 3, Text: map[string]string{"es": "van abajo", "en": "they go below"}},
	{ClueID: "cw_basement", Act: 4, Text: map[string]string{"es": "sotano", "en": "sotano"}},
	{ClueID: "cw_watching", Act: 4, Text: map[string]string{"es": "me miran", "en": "they watch me"}},
	{ClueID: "cw_numbers", Act: 5, Text: map[string]string{"es": "cuenta los numeros", "en": "count the numbers"}},
	{ClueID: "cw_alive", Act: 5, Text: map[string]string{"es": "sigue viva", "en": "still alive"}},
	{ClueID: "cw_open", Act: 6, Text: map[string]string{"es": "abre la puerta", "en": "open the door"}},
	{ClueID: "cw_down", Act: 6, Text: map[string]string{"es": "baja", "en": "come down"}},
}

// StoryDoc is one readable document in the /terminal scene. Docs with an
// empty Key open as soon as StoryStage reaches Act; keyed docs need their
// language-neutral Key typed in the terminal. Vera's numbered notes (one
// per act) spell out that act's ritual and where the next key lives.
type StoryDoc struct {
	ID    string // typed by the player: short, lowercase ASCII
	Act   int    // stage at which the doc is listed
	Key   string // "" = unlocked by stage alone
	Title map[string]string
	Body  map[string][]string
}

// StoryDocs lists every terminal document in listing order.
var StoryDocs = []StoryDoc{
	{
		ID: "aviso", Act: 1,
		Title: map[string]string{"es": "AVISO DE BIENVENIDA", "en": "WELCOME NOTICE"},
		Body: map[string][]string{
			"es": {
				"BIENVENIDO AL PUESTO 1. SU PREDECESOR: [REDACTADO].",
				"MOTIVO DE LA VACANTE: [REDACTADO].",
				"NO HAY NADA ESCRITO EN EL REVERSO DE ESTE AVISO.",
				"deja los avisos oficiales. lee mi primera nota: escribe 'leer nota1'.",
				"— v",
			},
			"en": {
				"WELCOME TO DESK 1. YOUR PREDECESSOR: [REDACTED].",
				"REASON FOR VACANCY: [REDACTED].",
				"THERE IS NOTHING WRITTEN ON THE BACK OF THIS NOTICE.",
				"skip the official notices. read my first note: type 'leer nota1' or 'read nota1'.",
				"— v",
			},
		},
	},
	{
		ID: "nota1", Act: 1,
		Title: map[string]string{"es": "NOTA DE VERA — 1", "en": "VERA'S NOTE — 1"},
		Body: map[string][]string{
			"es": {
				"si lees esto, ocupas mi silla. te dire como averiguar mi nombre.",
				"el sistema escucha los tropiezos. falla la MISMA palabra tres veces",
				"seguidas en un dia y te lo dira el mismo, en morado. tecleala entera.",
				"tambien despierta con cinco fallos en un dia PAR. cualquier palabra.",
				"mi nombre es la llave de mi expediente. escribelo aqui, solo, sin mas.",
				"— v",
			},
			"en": {
				"if you're reading this, you sit in my chair. here is how to learn my name.",
				"the system listens to stumbles. fail the SAME word three times",
				"in one day and it will tell you itself, in purple. type it all the way.",
				"it also wakes after five fails on an EVEN day. any word.",
				"my name is the key to my personnel file. type it here, alone, nothing else.",
				"— v",
			},
		},
	},
	{
		ID: "expediente", Act: 1, Key: "vera",
		Title: map[string]string{"es": "EXPEDIENTE DE PERSONAL", "en": "PERSONNEL FILE"},
		Body: map[string][]string{
			"es": {
				"EMPLEADA: VERA. PUESTO: MECANOGRAFA. ESCRITORIO: 1. EL TUYO.",
				"ALTA: HACE [ERROR] DIAS. BAJA: [ERROR].",
				"RENDIMIENTO: EXCELENTE. CUOTAS CUMPLIDAS: TODAS.",
				"NOTA DE RRHH: LA EMPLEADA HACE PREGUNTAS SOBRE EL DESTINO DE LAS PALABRAS.",
				"NOTA DE RRHH (POSTERIOR): LA EMPLEADA YA NO HACE PREGUNTAS.",
				"ULTIMO TURNO REGISTRADO: 77. SU SEGUNDA NOTA EXPLICA COMO USARLO.",
			},
			"en": {
				"EMPLOYEE: VERA. POSITION: TYPIST. DESK: 1. YOURS.",
				"HIRED: [ERROR] DAYS AGO. TERMINATED: [ERROR].",
				"PERFORMANCE: EXCELLENT. QUOTAS MET: ALL OF THEM.",
				"HR NOTE: EMPLOYEE ASKS QUESTIONS ABOUT WHERE THE WORDS GO.",
				"HR NOTE (LATER): EMPLOYEE NO LONGER ASKS QUESTIONS.",
				"LAST LOGGED SHIFT: 77. HER SECOND NOTE EXPLAINS HOW TO USE IT.",
			},
		},
	},
	{
		ID: "nota2", Act: 2,
		Title: map[string]string{"es": "NOTA DE VERA — 2", "en": "VERA'S NOTE — 2"},
		Body: map[string][]string{
			"es": {
				"encontraste mi expediente. bien. ya saben que lo buscaste. no importa.",
				"guarde mis registros en este sistema. la llave es mi ultimo turno,",
				"todo junto y en minusculas: turno77. escribelo aqui.",
				"otra cosa: un dia PERFECTO (ni un fallo, quince palabras o mas)",
				"tambien hace que conteste. le gusta la gente ordenada. le gustaba yo.",
				"— v",
			},
			"en": {
				"you found my file. good. they know you looked. it doesn't matter.",
				"i kept my logs in this system. the key is my last shift,",
				"as one lowercase word: turno77. type it here.",
				"one more thing: a PERFECT day (not one fail, fifteen words or more)",
				"also makes it answer. it likes tidy people. it liked me.",
				"— v",
			},
		},
	},
	{
		ID: "log1", Act: 2, Key: "turno77",
		Title: map[string]string{"es": "REGISTRO DE VERA — 1", "en": "VERA'S LOG — 1"},
		Body: map[string][]string{
			"es": {
				"dia 40. lo confirme: las palabras no se archivan. se ENVIAN a un nivel",
				"que no sale en los planos. algo las recibe. algo las estudia.",
				"no come palabras al azar. hay cuatro que pesan mas: sombra, puerta,",
				"abajo, silencio. teclea cada una TRES veces y sabras que me refiero.",
				"lo que hay abajo esta aprendiendo. mi tercera nota dice donde.",
				"sigue contando. yo ya no puedo parar.",
			},
			"en": {
				"day 40. confirmed: the words are not archived. they are SENT to a level",
				"that is not on the blueprints. something receives them. something studies them.",
				"it doesn't eat random words. four weigh more: shadow, door,",
				"below, silence. type each one THREE times and you'll see what i mean.",
				"the thing below is learning. my third note says where.",
				"keep counting. i can't stop anymore.",
			},
		},
	},
	{
		ID: "nota3", Act: 3,
		Title: map[string]string{"es": "NOTA DE VERA — 3", "en": "VERA'S NOTE — 3"},
		Body: map[string][]string{
			"es": {
				"las palabras pesadas funcionan, verdad? cada repeticion es un golpe",
				"en su puerta. ya sabes que hay puerta.",
				"el nivel se llama como el boton tapado del ascensor: sotano.",
				"esa palabra, aqui, abre mi segundo registro.",
				"y apunta: una racha de quince en un dia IMPAR tambien lo llama.",
				"— v",
			},
			"en": {
				"the heavy words work, don't they? every repeat is a knock",
				"on its door. you know there's a door now.",
				"the level is named after the covered elevator button: sotano.",
				"that word, typed here, opens my second log.",
				"and note this: a streak of fifteen on an ODD day also calls it.",
				"— v",
			},
		},
	},
	{
		ID: "log2", Act: 3, Key: "sotano",
		Title: map[string]string{"es": "REGISTRO DE VERA — 2", "en": "VERA'S LOG — 2"},
		Body: map[string][]string{
			"es": {
				"dia 61. baje con el boton tapado. lo vi. no es un archivo. no es una maquina.",
				"NO NACIO AQUI. llego en 1997. la empresa lo llama 'el huesped'.",
				"come palabras porque esta APRENDIENDO nuestro idioma. cada palabra",
				"que tecleamos arriba es una leccion. quiere hablar como nosotros.",
				"quiere subir y que nadie note la diferencia.",
				"los becarios lo oyen practicar. las horas que faltan se las queda el.",
				"voy a bajar otra vez. si lees esto, mira los numeros de tus cierres de dia.",
				"los numeros no mienten. rrhh cobra por mentir.",
			},
			"en": {
				"day 61. i went down with the covered button. i saw it. not an archive. not a machine.",
				"NOT BORN HERE. it arrived in 1997. the company calls it 'the guest'.",
				"it eats words because it is LEARNING our language. every word",
				"we type upstairs is a lesson. it wants to speak like us.",
				"it wants to come up and have nobody notice the difference.",
				"the interns hear it practicing. the missing hours? it keeps them.",
				"i'm going down again. if you read this, watch the numbers on your day-end reports.",
				"the numbers don't lie. hr is paid to lie.",
			},
		},
	},
	{
		ID: "camaras", Act: 4,
		Title: map[string]string{"es": "PARTE DE SEGURIDAD", "en": "SECURITY REPORT"},
		Body: map[string][]string{
			"es": {
				"CAMARA 3 Y CAMARA 4: ENFOCADAS A LA ESQUINA SURESTE DEL PASILLO B. NADIE LAS MOVIO.",
				"LA ESQUINA NO APARECE EN LAS GRABACIONES. SE GRABA UN TECLADO. TECLEANDO SOLO.",
				"AUDIO RECUPERADO: UNA VOZ REPITE PALABRAS DE OFICINA. EL ACENTO NO ES DE NINGUN PAIS.",
				"EL ACENTO MEJORA EN CADA CINTA.",
				"SEGURIDAD RECOMIENDA NO REVISAR ESTE PARTE. USTED YA LO HA REVISADO.",
			},
			"en": {
				"CAMERA 3 AND CAMERA 4: POINTED AT THE SOUTHEAST CORNER OF HALLWAY B. NOBODY MOVED THEM.",
				"THE CORNER DOES NOT APPEAR IN THE FOOTAGE. A KEYBOARD IS RECORDED INSTEAD. TYPING BY ITSELF.",
				"RECOVERED AUDIO: A VOICE REPEATS OFFICE WORDS. THE ACCENT IS FROM NO COUNTRY.",
				"THE ACCENT IMPROVES ON EVERY TAPE.",
				"SECURITY RECOMMENDS NOT REVIEWING THIS REPORT. YOU HAVE ALREADY REVIEWED IT.",
			},
		},
	},
	{
		ID: "nota4", Act: 4,
		Title: map[string]string{"es": "NOTA DE VERA — 4", "en": "VERA'S NOTE — 4"},
		Body: map[string][]string{
			"es": {
				"a estas alturas ya lo habras visto practicar: tus becarios escribiendo",
				"frases que nadie les dio, el reloj marcando horas de OTRO sitio.",
				"dejale terminar las frases a los becarios. cada una es una pista.",
				"la auditoria de horas esta sellada. se abre igual que todo aqui:",
				"con errores. TRES fallos en un dia PAR y aparece en este sistema.",
				"— v",
			},
			"en": {
				"by now you've seen it practice: your interns typing",
				"phrases nobody assigned, the clock showing hours from SOMEWHERE else.",
				"let the interns finish those phrases. each one is a clue.",
				"the hours audit is sealed. it opens like everything here:",
				"with mistakes. THREE fails on an EVEN day and it shows up in this system.",
				"— v",
			},
		},
	},
	{
		ID: "horas", Act: 4, Key: "",
		Title: map[string]string{"es": "AUDITORÍA DE HORAS", "en": "HOURS AUDIT"},
		Body: map[string][]string{
			"es": {
				"AUDITORIA INTERNA: LOS DIAS DE ESTE EDIFICIO DURAN MENOS DE LO CONTRATADO.",
				"DIFERENCIA MEDIA: 47 MINUTOS DIARIOS. DESTINO: NIVEL INFERIOR. USO: LECCIONES.",
				"LOS RELOJES MUESTRAN A VECES HORAS NO HOMOLOGADAS: 66:99, -0:13.",
				"NO SON AVERIAS. SON LA HORA LOCAL DEL LUGAR DEL QUE VINO EL HUESPED.",
				"EL HUESPED APRENDE NUESTRAS HORAS MAS RAPIDO QUE NUESTRAS PALABRAS.",
			},
			"en": {
				"INTERNAL AUDIT: THE DAYS IN THIS BUILDING LAST LESS THAN CONTRACTED.",
				"AVERAGE DIFFERENCE: 47 MINUTES PER DAY. DESTINATION: LOWER LEVEL. USE: LESSONS.",
				"THE CLOCKS SOMETIMES DISPLAY NON-APPROVED TIMES: 66:99, -0:13.",
				"THEY ARE NOT MALFUNCTIONS. THEY ARE THE LOCAL TIME OF WHERE THE GUEST CAME FROM.",
				"THE GUEST LEARNS OUR HOURS FASTER THAN OUR WORDS.",
			},
		},
	},
	{
		ID: "nota5", Act: 5,
		Title: map[string]string{"es": "NOTA DE VERA — 5", "en": "VERA'S NOTE — 5"},
		Body: map[string][]string{
			"es": {
				"ultima nota. el codigo de la puerta son CUATRO DIGITOS.",
				"aparecen en morado junto a la cuota, en tus cierres de dia. apuntalos",
				"y escribelos aqui, solos. la anomalia de nomina explica de donde salen.",
				"si no los ves aun: una racha de VEINTE en un dia multiplo de CINCO",
				"fuerza a que algo te hable. a mi me funciono.",
				"— v",
			},
			"en": {
				"last note. the door code is FOUR DIGITS.",
				"they show up in purple next to the quota on your day-end reports. write them",
				"down and type them here, alone. the payroll anomaly explains where they come from.",
				"if you don't see them yet: a streak of TWENTY on a day multiple of FIVE",
				"forces something to talk to you. it worked for me.",
				"— v",
			},
		},
	},
	{
		ID: "nomina", Act: 5,
		Title: map[string]string{"es": "ANOMALÍA DE NÓMINA", "en": "PAYROLL ANOMALY"},
		Body: map[string][]string{
			"es": {
				"CONTABILIDAD DETECTA UN PAGO RECURRENTE AL EMPLEADO 0 DESDE 1997.",
				"EL EMPLEADO 0 NO FIGURA EN NINGUNA PLANTILLA. NO ES DE ESTA PLANTILLA. NO ES DE ESTE PLANETA.",
				"EL PAGO NO ES EN DINERO: ES EN PALABRAS. CONCEPTO: MATERIAL DIDACTICO.",
				"LA REFERENCIA DEL PAGO SON CUATRO DIGITOS. SE IMPRIMEN EN LOS CIERRES DE DIA.",
				"TECLEE LOS CUATRO DIGITOS EN ESTE SISTEMA PARA VER EL REGISTRO DE ACCESOS.",
			},
			"en": {
				"ACCOUNTING DETECTS A RECURRING PAYMENT TO EMPLOYEE 0 SINCE 1997.",
				"EMPLOYEE 0 IS NOT ON ANY ROSTER. NOT ON THIS ROSTER. NOT FROM THIS PLANET.",
				"THE PAYMENT IS NOT IN MONEY: IT IS IN WORDS. REFERENCE: TEACHING MATERIALS.",
				"THE PAYMENT REFERENCE IS FOUR DIGITS. THEY PRINT ON THE DAY-END REPORTS.",
				"TYPE THE FOUR DIGITS INTO THIS SYSTEM TO SEE THE ACCESS LOG.",
			},
		},
	},
	{
		ID: "acceso", Act: 5, Key: "4413",
		Title: map[string]string{"es": "REGISTRO DE ACCESOS — NIVEL INFERIOR", "en": "ACCESS LOG — LOWER LEVEL"},
		Body: map[string][]string{
			"es": {
				"ACCESO 1: EMPLEADO 0 ('EL HUESPED'). ENTRADA: 1997, PROCEDENCIA NO TERRESTRE. SALIDA: PENDIENTE DE VOCABULARIO.",
				"ACCESO 2: VERA. ENTRADA: TURNO 77. SALIDA: RECHAZADA POR LA EMPLEADA.",
				"ACCESO 3: [RESERVADO]. ENTRADA: [PENDIENTE]. SALIDA: —",
				"LA PUERTA ACEPTA UNA SOLA COMBINACION: EL NOMBRE DE ELLA Y EL NUMERO DE EL.",
				"TODO JUNTO. SIN ESPACIOS. EN MINUSCULAS.",
				"USTED YA CONOCE AMBOS.",
			},
			"en": {
				"ENTRY 1: EMPLOYEE 0 ('THE GUEST'). IN: 1997, ORIGIN NON-TERRESTRIAL. OUT: PENDING VOCABULARY.",
				"ENTRY 2: VERA. IN: SHIFT 77. OUT: DECLINED BY THE EMPLOYEE.",
				"ENTRY 3: [RESERVED]. IN: [PENDING]. OUT: —",
				"THE DOOR TAKES A SINGLE COMBINATION: HER NAME AND ITS NUMBER.",
				"ALL TOGETHER. NO SPACES. LOWERCASE.",
				"YOU ALREADY KNOW BOTH.",
			},
		},
	},
	{
		ID: "puerta", Act: 6, Key: "vera4413",
		Title: map[string]string{"es": "LA PUERTA", "en": "THE DOOR"},
		Body: map[string][]string{
			"es": {
				"la puerta se abre. no hay monstruo de pelicula. hay un aula.",
				"algo enorme y paciente ocupa la mitad oscura. repite palabras de oficina",
				"con un acento que mejora frase a frase. 'nomina. cuota. buenos dias.'",
				"vera esta en la mitad iluminada, tecleando. lleva desde el turno 77.",
				"'aprende de lo que tecleamos', dice sin parar. 'cuando hable perfecto, subira.",
				"y nadie notara la diferencia. por eso tecleo yo: para elegir que aprende.'",
				"levanta una mano y senala la silla vacia junto a la suya.",
				"'dos teclados dosifican mejor que uno. tu decides que palabra le toca hoy.'",
				"arriba, tu pantalla escribe una palabra sola: gracias.",
			},
			"en": {
				"the door opens. there is no movie monster. there is a classroom.",
				"something huge and patient fills the dark half. it repeats office words",
				"in an accent that improves sentence by sentence. 'payroll. quota. good morning.'",
				"vera sits in the lit half, typing. she has been since shift 77.",
				"'it learns from what we type', she says without stopping. 'when it speaks perfectly, it will go up.",
				"and nobody will notice the difference. that's why i type: to choose what it learns.'",
				"she raises one hand and points at the empty chair next to hers.",
				"'two keyboards ration better than one. you decide which word it gets today.'",
				"upstairs, your screen types a word on its own: thank you.",
			},
		},
	},
	{
		ID: "informe", Act: StoryEndAct,
		Title: map[string]string{"es": "INFORME FINAL — CASO VERA", "en": "FINAL REPORT — VERA CASE"},
		Body: map[string][]string{
			"es": {
				"CASO CERRADO POR RRHH. CAUSA: RESUELTO POR EL OCUPANTE DEL PUESTO 1.",
				"EL HUESPED PERMANECE EN EL NIVEL INFERIOR. SU ACENTO MEJORA. SU PLAN NO HA CAMBIADO.",
				"VERA CONTINUA DOSIFICANDO SUS LECCIONES. AHORA POR TURNOS. USTED CONOCE EL HORARIO.",
				"RRHH LE RECUERDA QUE COBRA, EN PARTE, POR NO CONTAR NADA DE ESTO.",
				"LA INFILTRACION SIGUE PENDIENTE DE VOCABULARIO. SIGA TECLEANDO CON CRITERIO. CORDIALMENTE.",
			},
			"en": {
				"CASE CLOSED BY HR. CAUSE: SOLVED BY THE OCCUPANT OF DESK 1.",
				"THE GUEST REMAINS ON THE LOWER LEVEL. ITS ACCENT IMPROVES. ITS PLAN HAS NOT CHANGED.",
				"VERA CONTINUES RATIONING ITS LESSONS. IN SHIFTS NOW. YOU KNOW THE SCHEDULE.",
				"HR REMINDS YOU THAT YOU ARE PAID, IN PART, TO TELL NOBODY ABOUT THIS.",
				"THE INFILTRATION REMAINS PENDING VOCABULARY. KEEP TYPING RESPONSIBLY. CORDIALLY.",
			},
		},
	},
}

// MarkedWord is a normal dictionary word the system secretly watches:
// completing it enough times (across runs) fires its clue. Revealed to the
// player in Vera's first log ("the heavy ones").
type MarkedWord struct {
	ClueID string
	Act    int
	Count  int
	Word   map[string]string // per-language dictionary word
}

// MarkedWords lists the watched words; each pair exists in both
// dictionaries (tests check).
var MarkedWords = []MarkedWord{
	{ClueID: "mk_shadow", Act: 3, Count: 3, Word: map[string]string{"es": "sombra", "en": "shadow"}},
	{ClueID: "mk_door", Act: 3, Count: 3, Word: map[string]string{"es": "puerta", "en": "door"}},
	{ClueID: "mk_below", Act: 3, Count: 3, Word: map[string]string{"es": "abajo", "en": "below"}},
	{ClueID: "mk_silence", Act: 3, Count: 3, Word: map[string]string{"es": "silencio", "en": "silence"}},
}

// CorruptDoneScripts are the toasts shown right after typing a corrupt
// word: the signal acknowledges you.
var CorruptDoneScripts = map[string][]string{
	"es": {
		"REGISTRADO. NADIE HA VISTO ESTO.",
		"la senal te ha oido.",
		"eso no salia en tu contrato.",
		"ANOMALIA ARCHIVADA. GRACIAS POR SU COLABORACION.",
		"abajo lo han recibido. abajo lo estan repitiendo.",
		"RRHH NO RECONOCE ESTA PALABRA. VUELVA AL TRABAJO.",
	},
	"en": {
		"LOGGED. NOBODY SAW THIS.",
		"the signal heard you.",
		"that was not in your contract.",
		"ANOMALY FILED. THANK YOU FOR YOUR COOPERATION.",
		"below has received it. below is repeating it.",
		"HR DOES NOT RECOGNIZE THIS WORD. BACK TO WORK.",
	},
}

// MarkedHitScripts are the toasts shown when a marked word's counter
// completes: something noticed the knocking.
var MarkedHitScripts = map[string][]string{
	"es": {
		"has llamado suficientes veces. algo ha contestado.",
		"PALABRA PESADA CONTABILIZADA. NO PREGUNTE POR QUIEN.",
		"el contador de abajo acaba de moverse.",
		"deja de teclear esa palabra. es tarde. sigue tecleando.",
	},
	"en": {
		"you have knocked enough times. something answered.",
		"HEAVY WORD ACCOUNTED FOR. DO NOT ASK BY WHOM.",
		"the counter below just moved.",
		"stop typing that word. it's late. keep typing.",
	},
}

// PossessedScripts are the phrases a possessed intern types on its own
// (lowercase ASCII, they go through the intern letter-by-letter renderer).
var PossessedScripts = map[string][]string{
	"es": {
		"ella sigue aqui",
		"le oigo repetir palabras",
		"no me hagas mirar",
		"vera me pidio ayuda",
	},
	"en": {
		"she is still here",
		"i hear it repeat words",
		"dont make me look",
		"vera asked me for help",
	},
}

// PossessedDoneScripts are the toasts when a possessed intern finishes its
// phrase and snaps out of it.
var PossessedDoneScripts = map[string][]string{
	"es": {
		"tu becario no recuerda haber escrito eso.",
		"INCIDENTE DE BECARIO ARCHIVADO. NO HUBO INCIDENTE.",
		"el becario ha vuelto. una parte, al menos.",
	},
	"en": {
		"your intern doesn't remember typing that.",
		"INTERN INCIDENT FILED. THERE WAS NO INCIDENT.",
		"the intern is back. part of them, at least.",
	},
}

// ClockGlitchTexts are impossible clock readings, language-neutral.
var ClockGlitchTexts = []string{"66:99", "-0:13", "??:??", "77:77"}
