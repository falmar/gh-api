// Copyright 2016 David Lavieri.  All rights reserved.
// Use of this source code is governed by a MIT License
// License that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// declare variables
	var result map[string]interface{}
	var name string
	var email string
	var jwtToken string

	// take the code from request
	code := r.URL.Query().Get("code")

	// makes a request to get access token
	res, err := getAccessToken(code)

	// check for request error
	if err != nil {
		log.Println(err)
		w.WriteHeader(res.StatusCode)
		return
	}

	// decode json response
	json.NewDecoder(res.Body).Decode(&result)
	// close response body buffer
	res.Body.Close()

	// type casting interface to string
	if t, ok := result["access_token"].(string); ok {
		jwtToken = t
	}

	// make request to get user information
	res, err = getUserInformation(jwtToken)

	// check for request error
	if err != nil {
		log.Println(err)
		w.WriteHeader(res.StatusCode)
		return
	}

	// decode json response
	json.NewDecoder(res.Body).Decode(&result)
	// close response body buffer
	res.Body.Close()

	// parse template files
	tpl, err := template.ParseFiles("tpl/callback.gohtml")

	// check for parsing errors
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// type casting interface to string
	if e, ok := result["email"].(string); ok {
		email = e
	}

	// type casting interface to string
	if n, ok := result["name"].(string); ok {
		name = n
	}

	// execute template
	err = tpl.Execute(w, map[string]interface{}{
		"name":  name,
		"email": email,
	})

	// check if error any occured during template execution
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func getAccessToken(code string) (*http.Response, error) {
	// create http client
	client := http.Client{}

	// parse map into json
	reqBody, err := json.Marshal(map[string]string{
		"client_id":     os.Getenv("CLIENT_ID"),
		"client_secret": os.Getenv("CLIENT_SECRET"),
		"code":          code,
		"scope":         os.Getenv("SCOPE"),
	})

	// check for error and return if any
	if err != nil {
		return nil, err
	}

	// creates a buffer
	buf := bytes.NewBuffer(reqBody)

	// instance a new request
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", buf)

	// check for error and return if any
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// execute and return client response
	return client.Do(req)
}

func getUserInformation(token string) (*http.Response, error) {
	// create http client
	client := http.Client{}

	// instance a new request
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)

	// check for error and return if any
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// execute and return client response
	return client.Do(req)
}
