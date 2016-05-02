package model

import (
	// "fmt"
	"time"

	"github.com/atabek/gowebapp/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Note
// *****************************************************************************

// Note table contains the information for each note
type Clockin struct {
	ObjectID     bson.ObjectId `bson:"_id"`
	StudentID    string        `bson:"user_id"`
	InAt         time.Time     `db:"in_at"  bson:"in_at"`
	OutAt        time.Time     `db:"out_at" bson:"out_at"`
	Deleted      uint8         `db:"deleted" bson:"deleted"`
}

// ClockinID returns the note id
func (u *Clockin) ClockinID() string {
	r := ""
    r = u.ObjectID.Hex()
	return r
}
/*
// ClockinByID gets note by ID
func ClockinByID(studentID string, clockinID string) (Clockin, error) {
	var err error

	result := Clockin{}

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")

		// Validate the object id
		if bson.IsObjectIdHex(clockinID) {
			err = c.FindId(bson.ObjectIdHex(clockinID)).One(&result)
			if result.StudentID != bson.ObjectIdHex(userID) {
				result = Note{}
				err = ErrUnauthorized
			}
		} else {
			err = ErrNoResult
		}
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}
*/

// ClockinByStudentID gets all clockins for a student
func ClockinByStudentID(studentID string) ([]Clockin, error) {
	var err error

	var result []Clockin

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")

		// Validate the object id
		if bson.IsObjectIdHex(studentID) {
			err = c.Find(bson.M{"student_id": bson.ObjectIdHex(studentID)}).All(&result)
		} else {
			err = ErrNoResult
		}
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// ClockinCreate creates a note
func ClockinCreate(studentID string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")

		clockin := &Clockin{
			ObjectID:  bson.NewObjectId(),
			StudentID: studentID,
			InAt:      now,
			OutAt:     time.Time{},
			Deleted:   0,
		}

		err = c.Insert(clockin)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

/*
// ClockinUpdate updates a note
func ClockinUpdate(studentID string, clockinID string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")
		var clockin Clockin
		clockin, err = ClockinByID(studentID, clockinID)
		if err == nil {
			clockin.outAt = now
            err = c.UpdateId(bson.ObjectIdHex(clockinID), &clockin)
		} else {
			err = ErrUnauthorized
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
*/

/*
// NoteDelete deletes a note
func NoteDelete(userID string, noteID string) error {
	var err error

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

		var note Note
		note, err = NoteByID(userID, noteID)
		if err == nil {
			// Confirm the owner is attempting to modify the note
			if note.UserID.Hex() == userID {
				err = c.RemoveId(bson.ObjectIdHex(noteID))
			} else {
				err = ErrUnauthorized
			}
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
*/
