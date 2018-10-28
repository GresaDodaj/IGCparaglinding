package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

//AdminHandlerDelete is used
func AdminHandlerDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		del, err := collection.DeleteMany(context.Background(), nil)
		if err != nil {
			http.Error(w, "404", 400)
		}

		fmt.Fprint(w, del.DeletedCount)
	} else {
		http.Error(w, "", 400)
	}
}

//AdminHandlerGet is used
func AdminHandlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		length, err := collection.Count(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, length)
	} else {
		http.Error(w, "", 400)
	}
}
