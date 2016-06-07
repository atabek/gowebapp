package controller

import (
	"log"
	"net/http"
	"html"
	"encoding/json"
	"fmt"
	"strconv"
	// "reflect"

	"github.com/atabek/gowebapp/model"
	"github.com/atabek/gowebapp/shared/recaptcha"
	"github.com/atabek/gowebapp/shared/session"
	"github.com/atabek/gowebapp/shared/view"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/context"
)

func RegisterStudentGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "students/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	// Refill any form fields
	view.Repopulate([]string{"first_name", "last_name", "grade", "student_id",
		"fivedays", "caretype", "freereduced", "balance"}, r.Form, v.Vars)
	v.Render(w)
}

func RegisterStudentPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r,
		[]string{"first_name", "last_name", "grade", "student_id",
			"fivedays", "caretype", "freereduced", "balance"}); !validate {
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
	first_name  := r.FormValue("first_name")
	last_name   := r.FormValue("last_name")
	grade       := r.FormValue("grade")
	student_id  := r.FormValue("student_id")

	fivedays, err    := strconv.ParseBool(r.FormValue("fivedays"))
	caretype, err    := strconv.ParseInt(r.FormValue("caretype"), 10, 64)
	freereduced, err := strconv.ParseBool(r.FormValue("freereduced"))
	balance, err     := strconv.ParseFloat(r.FormValue("balance"), 64)

	// caretype32 := uint8(caretype)
	// balance32  := float32(balance)

	if err != nil{
		log.Println(err)
	}
	// fmt.Println("fivedays: ", fivedays, " and typeof: ", reflect.TypeOf(fivedays))
	// fmt.Println("caretype: ", caretype, " and typeof: ", reflect.TypeOf(caretype))
	// fmt.Println("freereduced: ", freereduced, " and typeof: ", reflect.TypeOf(freereduced))
	// fmt.Println("balance: ", balance, " and typeof: ", reflect.TypeOf(balance))

	// Get database result
	escaped_student_id := html.EscapeString(student_id)
	_, err = model.StudentBySID(escaped_student_id)

	if err == model.ErrNoResult { // If success (no user exists with that email)
		ex := model.StudentCreate(first_name, last_name, grade, escaped_student_id,
			balance, caretype, fivedays, freereduced)
		// Will only error if there is a problem with the query
		if ex != nil {
			log.Println(ex)
			sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
			sess.Save(r, w)
		} else {
			sess.AddFlash(view.Flash{"Account created successfully for: " + student_id, view.FlashSuccess})
			sess.Save(r, w)
			http.Redirect(w, r, "/students/create", http.StatusFound)
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

// Update student

// NotepadUpdateGET displays the note update page
func StudentUpdateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Get the note id
	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	studentID := params.ByName("id")

	//userID := fmt.Sprintf("%s", sess.Values["id"])

	// Get the note
	student, err := model.StudentBySID(studentID)
	if err != nil { // If the note doesn't exist
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	// Display the view
	v := view.New(r)
	v.Name = "students/update"
	v.Vars["token"]       = csrfbanana.Token(w, r, sess)
	v.Vars["first_name"]  = student.First_name;
	v.Vars["last_name"]   = student.Last_name;
	v.Vars["grade"]       = student.Grade;
	v.Vars["fivedays"]    = student.FiveDays;
	v.Vars["caretype"]    = student.CareType;
	v.Vars["freereduced"] = student.FreeReduced;
	v.Vars["balance"]     = student.Balance;
	v.Render(w)
}

// StudentUpdatePOST handles the student update form submission
func StudentUpdatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{
		"first_name", "last_name", "grade", "student_id",
		"fivedays", "caretype", "freereduced", "balance"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		StudentUpdateGET(w, r)
		return
	}

	// Get form values
	firstName := r.FormValue("first_name")
	lastName  := r.FormValue("last_name")
	grade     := r.FormValue("grade")

	// userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	studentID := params.ByName("id")

	// Get database result
	err := model.StudentUpdate(html.EscapeString(studentID), firstName, lastName, grade)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"student updated!", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/list", http.StatusFound)
		return
	}

	// Display the same page
	List(w, r)
}

// StudentDeleteGET handles the note deletion
func StudentDeleteGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// userID := fmt.Sprintf("%s", sess.Values["id"])

	var params httprouter.Params
	params = context.Get(r, "params").(httprouter.Params)
	studentID := params.ByName("id")
	fmt.Println(studentID)

	// Get database result
	err := model.StudentDelete(studentID)
	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Student deleted!", view.FlashSuccess})
		sess.Save(r, w)
	}

	http.Redirect(w, r, "/list", http.StatusFound)
	return
}


// Return the JSON of all the students in the database
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
