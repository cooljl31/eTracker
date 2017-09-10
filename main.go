package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

var tmpl *template.Template

// Categories fields from the exasol
type Categories struct {
	XMLName    xml.Name `xml:"category"`
	CustomerID string   `xml:"customer_id,attr"`
	SiteCode   string   `xml:"site_code"`
	Language   string   `xml:"language"`
	Name       string   `xml:"name"`
	From       string   `xml:"time>from"`
	Until      string   `xml:"time>to"`
	Date       string   `xml:"date"`
	Position   Tag      `xml:"positions"`
}

// Tag attr for positions
type Tag struct {
	CategoryID int    `xml:"category_id,attr"`
	Value      string `xml:",chardata"`
}

func (cat Categories) sConv(n string) string {
	return strings.ToLower(n)
}

func init() {
	tmpl = template.Must(template.ParseFiles("tpl.goxml"))
}

func main() {
	router := httprouter.New()
	router.GET(`/etracker/:site_code/:language/:customer_id/:from/:date`, category)
	http.ListenAndServe(`:3000`, router)
}

func category(res http.ResponseWriter, req *http.Request, prs httprouter.Params) {
	fmt.Printf("Thats cool %s\n", prs)
	pxml := Categories{
		CustomerID: prs.ByName("customer_id"),
		SiteCode:   prs.ByName("site_code"),
		Language:   prs.ByName("language"),
		From:       prs.ByName("from"),
		Name:       "hello",
		Until:      prs.ByName("from"),
		Date:       prs.ByName("date"),
		Position:   Tag{0, "heelo"},
	}

	out, err := xml.MarshalIndent(pxml, "\t", "\t")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(res, string(out)); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/xml")
}
