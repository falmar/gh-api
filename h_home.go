// Copyright 2016 David Lavieri.  All rights reserved.
// Use of this source code is governed by a MIT License
// License that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("tpl/home.gohtml")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?scope=%s&client_id=%s",
		os.Getenv("SCOPE"),
		os.Getenv("CLIENT_ID"),
	)

	err = tpl.Execute(w, map[string]string{
		"url": url,
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
