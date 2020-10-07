package main

import (
	"net/http"
	"context"
)

// Hauptseite
type RootPage struct {
	BasicPage
}

// c: Kontext
// r: Anfrage
// w: Ausgabe
func NewRootPage(c context.Context, r *http.Request,
	w http.ResponseWriter) *RootPage {
	page := new(RootPage)
	page.c = c
	page.r = r
	page.w = w
	
	page.b = new(HTMLBuffer)

	return page
}

func (page *RootPage) paint() {
	page.w.Header().Set("Content-Type", "text/html")

	page.beforeBody("Npackd")

	page.b.BE("div", "Npackd")

	carouselTemplate.Execute(page.w, "test")
	
	page.afterBody()

	page.w.Write([]byte(page.b.buf.String()))
}


