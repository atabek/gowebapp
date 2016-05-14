package controller

import (
	"log"
	"net/http"
	"html"

	"github.com/atabek/gowebapp/model"
	"github.com/atabek/gowebapp/shared/recaptcha"
	"github.com/atabek/gowebapp/shared/session"
	"github.com/atabek/gowebapp/shared/view"

	"github.com/josephspurrier/csrfbanana"
)

func RegisterStudentGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "registerstudent/registerstudent"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"first_name", "last_name", "grade", "student_id"}, r.Form, v.Vars)
	v.Render(w)
}

func RegisterStudentPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// // Prevent brute force login attempts by not hitting MySQL and pretending like it was invalid :-)
	// if sess.Values["register_attempt"] != nil && sess.Values["register_attempt"].(int) >= 5 {
	// 	log.Println("Brute force register prevented")
	// 	http.Redirect(w, r, "/registerstudent", http.StatusFound)
	// 	return
	// }

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"first_name", "last_name", "grade", "student_id"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		RegisterStudentGET(w, r)
		return
	}

	// Validate with Google reCAPTCHA
	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		RegisterStudentGET(w, r)
		return
	}

	// Get form values
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	grade := r.FormValue("grade")
	student_id := r.FormValue("student_id")
	//password, errp := passhash.HashString(r.FormValue("password"))

	// // If password hashing failed
	// if errp != nil {
	// 	log.Println(errp)
	// 	sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
	// 	sess.Save(r, w)
	// 	http.Redirect(w, r, "/register", http.StatusFound)
	// 	return
	// }

	// Get database result
	_, err := model.StudentBySID(html.EscapeString(student_id))

	if err == model.ErrNoResult { // If success (no user exists with that email)
		ex := model.StudentCreate(first_name, last_name, grade, html.EscapeString(student_id))
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
		} else {
			sess.AddFlash(view.Flash{"Account created successfully for: " + student_id, view.FlashSuccess})
			sess.Save(r, w)
			http.Redirect(w, r, "/registerstudent", http.StatusFound)
			return
		}
	} else if err != nil { // Catch all other errors
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else { // Else the user already exists
		sess.AddFlash(view.Flash{"Account already exists for: " + student_id, view.FlashError})
		sess.Save(r, w)
	}

	// Display the page
	RegisterStudentGET(w, r)
}
