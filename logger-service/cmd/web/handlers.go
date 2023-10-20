package main

import (
	"fmt"
	"log"
	"net/http"
)

// authServiceURL is the url to the authentication service. Since we're using
// Docker, we specify the appropriate entry from docker-compose.yml
const authServiceURL = "http://authentication-service/authenticate"

// JSONPayload is the type for JSON posted to this API
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog is the handler to accept a post request consisting of json payload,
// and then write it to Mongo
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	log.Println("WriteLog")
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert the data
	err := app.logEvent(requestPayload.Name, requestPayload.Data)
	if err != nil {
		log.Println(err)
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// create the response we'll send back as JSON
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	// write the response back as JSON
	_ = app.writeJSON(w, http.StatusAccepted, resp)
}

// Logout logs the user out and redirects them to the login page
func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// log the event
	_ = app.logEvent("authentication", fmt.Sprintf("%s logged out of the logger service", "email"))

	// redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
