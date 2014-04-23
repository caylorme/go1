package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Default:  %s", r.URL.Path)
}
func APIHandler(w http.ResponseWriter, r *http.Request) {
	Test := strings.Split(r.URL.Path, "/")
	switch strings.ToLower(Test[2]) {
	case "wow":
		fmt.Fprintf(w, "WOW!")
	default:
		Test2 := r.URL.Query()
		fmt.Printf("Test: %v\n", Test2.Get("test"))
		js, err := json.Marshal(Test2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/API/", APIHandler)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
