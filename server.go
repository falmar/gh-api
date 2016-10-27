// Copyright 2016 David Lavieri.  All rights reserved.
// Use of this source code is governed by a MIT License
// License that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	port := ":" + os.Getenv("PORT")

	router := httprouter.New()

	setRoutes(router)

	log.Println(fmt.Sprintf("Listening on %s", port))

	http.ListenAndServe(port, router)
}
