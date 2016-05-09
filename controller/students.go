package controller

import (
	"net/http"
	"encoding/json"
	"fmt"
	"log"

	"github.com/atabek/gowebapp/model"
	"github.com/atabek/gowebapp/shared/session"
	"github.com/atabek/gowebapp/shared/view"
)

// Displays the About page
func StudentsJsonGET(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	students, err := model.StudentsGet()
	if err == nil{
		// Marshal provided interface into JSON structure
		sj, _ := json.Marshal(students)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", sj)
	} else {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}
}
