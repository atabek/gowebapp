package controller

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/atabek/gowebapp/model"
	"github.com/atabek/gowebapp/shared/session"
	"github.com/atabek/gowebapp/shared/view"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
)

// NotepadReadGET displays the notes in the notepad
// func NotepadReadGET(w http.ResponseWriter, r *http.Request) {
// 	// Get session
// 	sess := session.Instance(r)
//
// 	userID := fmt.Sprintf("%s", sess.Values["id"])
//
// 	notes, err := model.NotesByUserID(userID)
// 	if err != nil {
// 		log.Println(err)
// 		notes = []model.Note{}
// 	}
//
// 	// Display the view
// 	v := view.New(r)
// 	v.Name = "notepad/read"
// 	v.Vars["first_name"] = sess.Values["first_name"]
// 	v.Vars["notes"] = notes
// 	v.Render(w)
// }

// ClockinCreateGET displays the clockin creation page
func ClockinCreateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "clockin/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

// ClockinCreatePOST handles the note creation form submission
func ClockinCreatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"student_id"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		ClockinCreateGET(w, r)
		return
	}

	// Get form values
	student_id := r.FormValue("student_id")

	//studentID := fmt.Sprintf("%s", sess.Values["studentID"])

	// Get database result
	err := model.ClockinCreate(student_id)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		//sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
        sess.AddFlash(view.Flash{err.Error(), view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Student clockin/out successful!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/clockin/create", http.StatusFound)
		return
	}

	// Display the same page
	ClockinCreateGET(w, r)
}

// ClockinByStudentIdJsonGET displays the note update page
func ClockinByStudentIdJsonGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the note id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	studentID := params.ByName("id")

	// Get the clockins of a particular Student
	clockins, err := model.ClockinsByStudentID(studentID)

	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	// Marshal provided interface into JSON structure
	cj, _ := json.Marshal(clockins)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", cj)
}

/*
// NotepadUpdateGET displays the note update page
func NotepadUpdateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	var results []Clockin
	// Get the note id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	studentID := params.ByName("id")

	// Get the clockins
	results, err = ClockinsByStudentID(student_id)
	fmt.Println(results)

	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	// Marshal provided interface into JSON structure
	sj, _ := json.Marshal(results)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", sj)

	// Display the view
	v := view.New(r)
	v.Name = "clockins/update"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["note"] = note.Content
	v.Render(w)
}


// NotepadUpdatePOST handles the note update form submission
func NotepadUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"note"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		NotepadUpdateGET(w, r)
		return
	}

	// Get form values
	content := r.FormValue("note")

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	// Get database result
	err := model.NoteUpdate(content, userID, noteID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note updated!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/notepad", http.StatusFound)
		return
	}

	// Display the same page
	NotepadUpdateGET(w, r)
}

// NotepadDeleteGET handles the note deletion
func NotepadDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	noteID := params.ByName("id")

	// Get database result
	err := model.NoteDelete(userID, noteID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Note deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/notepad", http.StatusFound)
	return
}
*/
