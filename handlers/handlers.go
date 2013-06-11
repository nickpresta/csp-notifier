package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"labix.org/v2/mgo"
)

type ReportContent struct {
	DocumentURI       string `json:"document-uri" bson:"documentURI"`
	Referrer          string `json:"referrer" bson:"referrer"`
	BlockedURI        string `json:"blocked-uri" bson:"blockedURI"`
	ViolatedDirective string `json:"violated-directive" bson:"violatedDirective"`
	OriginalPolicy    string `json:"original-policy" bson:"originalPolicy"`
}

type Report struct {
	Content ReportContent `json:"csp-report" bson:"content"`
}

type BSONReportContent struct {
	DocumentURI       string `json:"documentURI" bson:"documentURI"`
	Referrer          string `json:"referrer" bson:"referrer"`
	BlockedURI        string `json:"blockedURI" bson:"blockedURI"`
	ViolatedDirective string `json:"violatedDirective" bson:"violatedDirective"`
	OriginalPolicy    string `json:"originalPolicy" bson:"originalPolicy"`
}

var templates = template.Must(template.ParseFiles("index.html"))

func MainHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", "")
}

func AddHandler(w http.ResponseWriter, r *http.Request, session *mgo.Session) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse POST body: %+v", err), http.StatusInternalServerError)
	}

	var report Report
	if err := json.Unmarshal(body, &report); err != nil {
		http.Error(w, fmt.Sprintf("Could not parse POST body: %+v", err), http.StatusInternalServerError)
	}

	c := session.DB("").C("reports")
	if err = c.Insert(&report.Content); err != nil {
		fmt.Printf("Failed to insert: %s\n", err)
	}
}

func ReportHandler(w http.ResponseWriter, r *http.Request, session *mgo.Session) {
	w.Header().Set("Content-Type", "application/json")
	c := session.DB("").C("reports")

	var reports []BSONReportContent
	err := c.Find(nil).Sort("-_id").Limit(100).All(&reports)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	out, err := json.Marshal(reports)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, string(out))
}
