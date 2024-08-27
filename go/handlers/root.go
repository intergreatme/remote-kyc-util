/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package handlers

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
// 	path := "./template/form.html"

// 	ts, err := template.ParseFiles(path) //parse the html file homepage.html
// 	if err != nil {                      // if there is an error
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}
// 	err = ts.Execute(w, nil) //execute the file form.html
// 	if err != nil {          // if there is an error
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 		return
// 	}
// }

func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	originTxID := uuid.New().String()
	ordernumber := fmt.Sprintf("%d%04d", time.Now().Unix(), rand.Intn(10000))

	data := map[string]interface{}{
		"OriginTxID":  originTxID,
		"OrderNumber": ordernumber,
	}
	tmpl := template.Must(template.ParseFiles("./template/form.html"))
	tmpl.Execute(w, data)
}
