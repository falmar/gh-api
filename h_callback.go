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
	var result map[string]interface{}
	var name string
	var email string
	var jwtToken string

	code := r.URL.Query().Get("code")

	res, err := getAccessToken(code)

	if err != nil {
		log.Println(err)
		w.WriteHeader(res.StatusCode)
		return
	}

	json.NewDecoder(res.Body).Decode(&result)

	res.Body.Close()

	if t, ok := result["access_token"].(string); ok {
		jwtToken = t
	}

	res, err = getUserInformation(jwtToken)

	if err != nil {
		log.Println(err)
		w.WriteHeader(res.StatusCode)
		return
	}

	json.NewDecoder(res.Body).Decode(&result)

	res.Body.Close()

	tpl, err := template.ParseFiles("tpl/callback.gohtml")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if e, ok := result["email"].(string); ok {
		email = e
	}

	if n, ok := result["name"].(string); ok {
		name = n
	}

	err = tpl.Execute(w, map[string]interface{}{
		"name":  name,
		"email": email,
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func getAccessToken(code string) (*http.Response, error) {
	client := http.Client{}

	reqBody, err := json.Marshal(map[string]string{
		"client_id":     os.Getenv("CLIENT_ID"),
		"client_secret": os.Getenv("CLIENT_SECRET"),
		"code":          code,
		"scope":         os.Getenv("SCOPE"),
	})

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", buf)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func getUserInformation(token string) (*http.Response, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return client.Do(req)
}
