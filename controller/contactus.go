package controller

import (
	"net/http"

	"github.com/atabek/gowebapp/shared/session"
	"github.com/atabek/gowebapp/shared/view"
)

// Displays the default home page
func Students(w http.ResponseWriter, r *http.Request) {
	// Get session
	session := session.Instance(r)

	if session.Values["id"] != nil {
		// Display the view
		v := view.New(r)
		v.Name = "students/auth"
		v.Vars["first_name"] = session.Values["first_name"]
		v.Render(w)
	} else {
		// Display the view
		v := view.New(r)
		v.Name = "students/anon"
		v.Render(w)
		return
	}
}
