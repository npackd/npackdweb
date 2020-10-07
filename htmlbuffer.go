package main

import (
	"bytes"
	"html"
	"time"
)

// Puffer für die HTML-Ausgabe.
type HTMLBuffer struct {
	buf bytes.Buffer
}

// Schließt einen Tag.
func (b *HTMLBuffer) E(tag string) *HTMLBuffer {
	b.Unencoded("</")
	b.Unencoded(tag)
	return b.Unencoded(">")
}

// Gibt den gegebenen Text aus.
func (b *HTMLBuffer) T(txt string) *HTMLBuffer {
	b.buf.WriteString(html.EscapeString(txt))
	return b
}

// Gibt den gegebenen Text aus, ohne die HTML-Kodierung durchzuführen.
func (b *HTMLBuffer) Unencoded(txt string) *HTMLBuffer {
	b.buf.WriteString(txt)
	return b
}

func (b *HTMLBuffer) CollapsedDateTime(t time.Time) *HTMLBuffer {
	b.Unencoded(`
<a data-toggle="collapse" href="#collapseExample" aria-expanded="false" aria-controls="collapseExample">
  ...
	</a><div class="collapse" id="collapseExample">Additional text</div>`)
		return b
}

// Tag-Anfang
//
// tag: Name des Tags
// params: Namen und Werte für die Attribute. Der letzte Wert wird als Inhalt
// 	des Tags ausgegeben, falls die Anzahl der Werte ungerade ist. Die
//      Attribute mit leeren Namen werden ignoriert.
func (b *HTMLBuffer) B(tag string, params ...string) *HTMLBuffer {
	b.Unencoded("<")
	b.Unencoded(tag)
	n := len(params)
	if n % 2 != 0 {
		n--
	}
	for i := 1; i < n; i += 2 {
		name := params[i - 1]
		if name != "" {
			value := params[i]
			b.Unencoded(" ").Unencoded(name).Unencoded("=")
			b.Unencoded("\"").T(value).Unencoded("\"")
		}
	}
	b.Unencoded(">")
	if len(params) % 2 != 0 {
		b.T(params[len(params) - 1])
	}
	return b
}

// Kompletter Tag
//
// tag: Name des Tags
// params: Namen und Werte für die Attribute. Der letzte Wert wird als Inhalt
// 	des Tags ausgegeben, falls die Anzahl der Werte ungerade ist. Die
//      Attribute mit leeren Namen werden ignoriert.
func (b *HTMLBuffer) BE(tag string, params ...string) *HTMLBuffer {
	b.B(tag, params...)
	b.E(tag)
	return b
}

// "input" vom Typ "hidden"
func (b *HTMLBuffer) Hidden(id, name, value string) *HTMLBuffer {
	return b.BE("input", "type", "hidden", "name", name, "value", value, "id", id)
}

// "input" vom Typ "submit"
func (b *HTMLBuffer) SubmitButton(title string) *HTMLBuffer {
	return b.BE("input", "type", "submit", "class", "btn btn-default", "value", title)
}

// "input" vom Typ "submit"
func (b *HTMLBuffer) SubmitButton2(name, value, title string) *HTMLBuffer {
	return b.BE("button", "type", "submit", "class", "btn btn-default", "name", name, "value", value, title)
}

// "input" vom Typ "button" mit window.location.href=
func (b *HTMLBuffer) ChangeLocationButton(title string, loc string) *HTMLBuffer {
	return b.BE("input", "type", "button", "class", "btn btn-default",
		"value", title, "onclick", "window.location.href='" + loc + "'")
}

// "input" vom Typ "button" mit JavaScript
func (b *HTMLBuffer) JSButton(title string, code string) *HTMLBuffer {
	return b.BE("input", "type", "button", "class", "btn btn-default",
		"value", title, "onclick", code)
}

// Eingabefeld
func (b *HTMLBuffer) TextInput(title string, name string, value string) *HTMLBuffer {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.BE("input", "type", "text", "class", "form-control", "id", name,
		"name", name, "value", value)
	b.E("div")
	b.E("div")
	return b
}

// Datei hochladen
func (b *HTMLBuffer) FileInput(title string, name string) *HTMLBuffer {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.BE("input", "type", "file", "id", name,
		"name", name)
	b.E("div")
	b.E("div")
	return b
}

// Text
func (b *HTMLBuffer) TextArea(title string, name string, value string) *HTMLBuffer {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.BE("textarea", "rows", "2", "class", "form-control", "id", name,
		"name", name, value)
	b.E("div")
	b.E("div")
	return b
}

// Eingabefeld für ein Datum
// value: yyyy-mm-dd
func (b *HTMLBuffer) DateInput(title string, name string, value string) *HTMLBuffer {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.BE("input", "type", "date", "class", "form-control", "id", name,
		"name", name, "value", value)
	b.E("div")
	b.E("div")
	return b
}

// Eingabefeld für ein Datum mit Uhrzeit
// value: yyyy-mm-ddTHH:MM
func (b *HTMLBuffer) DateTimeInput(title string, name string, value string) *HTMLBuffer {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.BE("input", "type", "datetime-local", "class", "form-control", "id", name,
		"name", name, "value", value)
	b.E("div")
	b.E("div")
	return b
}

// Eingabefeld für eine Zahl
// value: 123123.123
func (b *HTMLBuffer) NumberInput(title string, name string, value string) *HTMLBuffer {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.BE("input", "type", "number", "class", "form-control", "id", name,
		"name", name, "value", value, "step", "0.00001")
	b.E("div")
	b.E("div")
	return b
}

// Mehrere Radiobuttons.
// title: Überschrift
// name: Feldname
// labels: Überschriften für die Werte
// values: Werte (werden auch für "id" verwendet)
// selected: der ausgewählte Wert
func (b *HTMLBuffer) RadioButtons(title string, name string,
	labels []string, values [] string, selected string) {
	b.B("div", "class", "form-group")
	b.BE("label", "for", values[0], "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	for i := range labels {
		checked := ""
		if values[i] == selected {
			checked = "checked"
		}
		b.B("label", "class", "radio-inline")
		b.BE("input", "type", "radio", "name", name,
			"value", values[i], "id", values[i],
			checked, "checked", labels[i])
		b.E("label")
	}
	b.E("div")
	b.E("div")
}

// Ausgabe der Überschrift und HTML Tags vor einem Eingabefeld.
// title: Überschrift
// name: Feldname und ID für "select"
func (b *HTMLBuffer) BeforeInput(title string, name string) {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
}

// Ausgabe der HTML Tags nach einem Eingabefeld.
func (b *HTMLBuffer) AfterInput() {
	b.E("div")
	b.E("div")
}

// Auswahlbox.
// title: Überschrift
// name: Feldname und ID für "select"
// labels: Überschriften für die Werte
// values: Werte
// selected: der ausgewählte Wert
func (b *HTMLBuffer) Combobox(title string, name string,
	labels []string, values [] string, selected string) {
	b.B("div", "class", "form-group")
	b.BE("label", "for", name, "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.B("select", "class", "form-control", "name", name, "id", name)
	for i := range labels {
		sel := ""
		if values[i] == selected {
			sel = "selected"
		}
		b.BE("option", "value", values[i], 
			sel, "selected", labels[i])
	}
	b.E("select")
	b.E("div")
	b.E("div")
}

// Statisches Feld.
func (b *HTMLBuffer) TextStatic(title string, value string) *HTMLBuffer {
	b.B("div", "class", "form-group")
	b.BE("label", "class", "col-sm-2 control-label", title)
	b.B("div", "class", "col-sm-10")
	b.BE("p", "class", "form-control-static", value)
	b.E("div"); // class=col-sm-10
	b.E("div")
	return b
}

// Krümmelpfad.
//
// titlesAndURLs: Titel und Links (ungerade Anzahl von Parametern): title0, url0, title1, url1, ..., titleN
func (b *HTMLBuffer) Breadcrumb(titlesAndURLs ...string) *HTMLBuffer {
	b.B("ol", "class", "breadcrumb")
	for i := 0; i < len(titlesAndURLs) - 1; i += 2 {
		b.B("li")
		b.B("a", "href", titlesAndURLs[i + 1])
		b.T(titlesAndURLs[i])
		b.E("a")
		b.E("li")
	}
	b.B("li", "class", "active")
	b.T(titlesAndURLs[len(titlesAndURLs) - 1])
	b.E("li")
	b.E("ol")
	return b
}

