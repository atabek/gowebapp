package controller

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/atabek/gowebapp/model"
)

// Displays the About page
func StudentsJSONGet(w http.ResponseWriter, r *http.Request) {
	// // Display the view
	// v := view.New(r)
	// v.Name = "about/about"
	// v.Render(w)
	students, err := model.StudentsGet()
	if err == nil{
		// Marshal provided interface into JSON structure
		sj, _ := json.Marshal(students)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", sj)
	} else {
		fmt.Println(err)
	}

}
