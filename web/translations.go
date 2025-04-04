package web

import "html"

var translations = map[string]map[string]string{
	"en": {
		"how-to":          "How to use? <b>(READ BEFORE USING)</b>",
		"how-to.1":        "Open <i>Scheduler</i>-><i>Scheduler export</i> in <i>STiNE</i>",
		"how-to.2":        "Select the first month of the desired semester",
		"how-to.3":        "Leave <i>Calendar week</i> empty",
		"how-to.4":        "Click on <i>Export Appointments</i>",
		"how-to.5":        "Click on <i>Calendar file</i> to download all events for the selected month as an <i>.ics</i> file",
		"how-to.6":        "Repeat for all months in the desired semester (<i>ignore months without any appointments</i>)",
		"how-to.7":        "Click on the file input below on this page and select <b>all</b> downloaded <i>.ics</i> files",
		"how-to.8":        "Click <i>Format Files</i>",
		"how-to.9":        "Import the resulting file into your calendar application of choice",
		"how-to.warning":  "DO NOT USE FOR FILES NOT EXPORTED FROM STINE!",
		"what-does-it-do": "What does this tool do?",
		"what-does-it-do.preamble": `
			ICS files exported from the scheduler export tool in STiNE do not
			conform to the official iCalendar spec (<i>RFC 5545</i>). <br/>
			Many calendar application therefore have problems opening them.
		`,
		"what-does-it-do.list-title": "This tool…",
		"what-does-it-do.1":          "…removes <i>NULL</i> bytes",
		"what-does-it-do.2":          "…converts file encoding from <i>ISO8859-1</i> to <i>UTF8</i>",
		"what-does-it-do.3":          "…removes empty lines and wraps long lines",
		"what-does-it-do.4":          "…ensures line endings are <i>CRLF</i>",
		"what-does-it-do.5":          "…combines all files into one",
		"what-does-it-do.6":          "…detects recurring events and merges them using <i>RRULEs</i>",
		"what-does-it-do.more-info":  `For more info see <a class="link" href="https://github.com/juho05/stine-ical-formatter/tree/main/formatter">here</a>.`,
		"select-files":               "Select iCalendar (.ics) files:",
		"format-files":               "Format Files",
		"error.rate-limit":           "Rate limit reached. Please wait a few seconds.",
		"error.unexpected":           "An unexpected error occured.",
		"error.files-too-large":      "Files too large (sum must be <5MB).",
		"error.no-files":             "No files selected.",
		"error.not-ics":              "At least one file is not an iCalendar file (extension: .ics)",
		"error.failed-to-open":       "Files could not be opened.",
		"error.format-failed":        "Formatting failed. Are you sure you uploaded the correct files?",
	},
	"de": {
		"how-to":          "Anleitung <b>(VORM NUTZEN LESEN)</b>",
		"how-to.1":        "Öffne <i>Termine</i>-><i>Terminexport</i> in <i>STiNE</i>",
		"how-to.2":        "Wähle den ersten Monat des gewünschten Semesters",
		"how-to.3":        "Lass <i>Kalenderwoche</i> unausgefüllt",
		"how-to.4":        "Drücke <i>Termine exportieren</i>",
		"how-to.5":        "Drücke <i>Kalenderdatei</i>, um alle Termine des gewählten Monats als eine <i>.ics</i> Datei herunterzuladen",
		"how-to.6":        "Wiederhole die Anweisungen für alle Monate des gewünschten Semesters (<i>ignoriere Monate ohne Termine</i>)",
		"how-to.7":        "Drücke auf die Dateieingabe auf dieser Seite und selektiere <b>alle</b> heruntergeladenen <i>.ics</i> Dateien",
		"how-to.8":        "Drücke <i>Dateien Formatieren</i>",
		"how-to.9":        "Importiere die resultierende Datei in ein Kalenderprogramm deiner Wahl",
		"how-to.warning":  "NUTZE DIESES TOOL NUR FÜR DATEIEN, DIE AUS STINE EXPORTIERT WURDEN!",
		"what-does-it-do": "Was tut dieses Tool?",
		"what-does-it-do.preamble": `
			ICS Dateien aus STiNEs Terminexport entsprechen nicht
			der offiziellen iCalendar Spezifikation (<i>RFC 5545</i>). <br/>
			Deswegen haben viele Kalenderprogramme Probleme, diese Dateien zu öffnen.
		`,
		"what-does-it-do.list-title": "Dieses Programm…",
		"what-does-it-do.1":          "…entfernt <i>NULL</i> Bytes",
		"what-does-it-do.2":          "…konvertiert die Dateicodierung von <i>ISO8859-1</i> zu <i>UTF8</i>",
		"what-does-it-do.3":          "…entfernt leere Zeilen und bricht lange Zeilen um",
		"what-does-it-do.4":          "…sichert, dass Zeilenenden als <i>CRLF</i> codiert sind",
		"what-does-it-do.5":          "…fügt alle Dateien zu einer zusammen",
		"what-does-it-do.6":          "…erkennt wiederkehrende Termine und fügt sie zu einem wiederkehrenden Termin mit <i>RRULEs</i> zusammen",
		"what-does-it-do.more-info":  `Für mehr Infos, siehe <a class="link" href="https://github.com/juho05/stine-ical-formatter/tree/main/formatter">hier</a>.`,
		"select-files":               "Wähle iCalendar (.ics) Dateien:",
		"format-files":               "Dateien Formatieren",
		"error.rate-limit":           "Ratelimit erreicht. Bitte warte ein paar Sekunden.",
		"error.unexpected":           "Ein unerwarteter Fehler ist aufgetreten.",
		"error.files-too-large":      "Dateien zu groß (Summe muss <5MB sein).",
		"error.no-files":             "Keine Dateien ausgewählt.",
		"error.not-ics":              "Mindestens eine Datei ist keine iCalendar Datei (Dateiendung: .ics)",
		"error.failed-to-open":       "Dateien konnten nicht geöffnet werden.",
		"error.format-failed":        "Formatierung fehlgeschlagen. Bist du sicher, dass du die richtigen Dateien hochgeladen hast?",
	},
}

func Translate(lang, key string) string {
	t, ok := translations[lang]
	if !ok {
		if lang != "en" {
			return Translate("en", key)
		}
		return html.EscapeString(key)
	}
	if str, ok := t[key]; ok {
		return str
	}
	if lang != "en" {
		return Translate("en", key)
	}
	return html.EscapeString(key)
}
