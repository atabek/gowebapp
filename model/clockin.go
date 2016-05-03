package model

import (
	"fmt"
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
	StudentID    string        `bson:"student_id"`
	InAt         int64     `db:"in_at"  bson:"in_at"`
	OutAt        int64     `db:"out_at" bson:"out_at"`
	TotalTime    int64       `db:"total_time" bson:"total_time"`
	IsOut        bool          `db:"is_out" bson:"is_out"`
}

// ClockinID returns the clockin id
func (clockin *Clockin) ClockinID() string {
	r := ""
    r = clockin.ObjectID.Hex()
	return r
}

// ClockinByID gets note by ID
func ClockinByID(clockinID string) (Clockin, error) {
	var err error

	result := Clockin{}

	if !database.CheckConnection() {
		err = ErrUnavailable
		return result, err
	}

	// Create a copy of mongo
	session := database.Mongo.Copy()
	defer session.Close()
	c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")

	// Validate the object id
	if bson.IsObjectIdHex(clockinID) {
		err = ErrNoResult
		return result, err
	}

	err = c.FindId(bson.ObjectIdHex(clockinID)).One(&result)
	if err != nil {
		err = ErrUnauthorized
		return result, err
	}

	return result, standardizeError(err)
}

// ClockinByStudentID gets the last clockin for a student
func LastClockinByStudentID(student_id string) (Clockin, error) {
	var err error

	var result Clockin
    timeZero := 0

	if !database.CheckConnection() {
		err = ErrUnavailable
		return result, err
	}
	// Create a copy of mongo
	session := database.Mongo.Copy()
	defer session.Close()
	c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")

	err = c.Find(bson.M{"student_id": student_id, "out_at": timeZero}).One(&result)

	return result, standardizeError(err)
}

// ClockinCreate creates a clockin
func ClockinCreate(student_id string) error {
	var err error

	now := time.Now().Unix()
	if !database.CheckConnection() {
        err = ErrUnavailable
        return err
    }
	// Create a copy of mongo
	session := database.Mongo.Copy()
	defer session.Close()
	c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")
    s := session.DB(database.ReadConfig().MongoDB.Database).C("student")

    result := Student{}
    zeroTime := 0

    err = s.Find(bson.M{"student_id": student_id}).One(&result)
    if err != nil{
        err = ErrNoSuchStudent
        return err
    }

    colQuerier := bson.M{"student_id": student_id,
		"out_at": zeroTime, "is_out": false}
    count, err := c.Find(colQuerier).Count()

    if count == 0 {
		clockin := &Clockin{
			ObjectID:  bson.NewObjectId(),
			StudentID: student_id,
			InAt:      now,
		}
        err = c.Insert(clockin)
    }
    if count == 1 {
		// Error clockin by student id retrieving
        if err != nil{
            return err
        }

		// Update
		clockin, err := LastClockinByStudentID(student_id)
		if err != nil{
			return err
		}

		total_time := now - clockin.InAt
        change := bson.M{"$set":
			bson.M{"out_at": now, "total_time": total_time, "is_out": true}}
		err = c.Update(colQuerier, change)

        if err != nil{
            // Error updating the clockin
            return err
        }
    }
    if count > 1 {
        fmt.Println("You need to fix the clockin count for this student.")
    }

	return standardizeError(err)
}


// ClockinDelete deletes a note
// Also add an admin checking for delete functionality
func ClockinDeleteByID(clockinID string) error {
	var err error

	if !database.CheckConnection() {
		err = ErrUnavailable
		return err
	}
	// Create a copy of mongo
	session := database.Mongo.Copy()
	defer session.Close()
	c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")

	//var clockin Clockin
	_, err = ClockinByID(clockinID)
	if err != nil {
		return ErrNoResult
	}
	// Confirm the owner is attempting to modify the note
	// if note.UserID.Hex() == userID {
	// 	err = c.RemoveId(bson.ObjectIdHex(noteID))
	// }
	err = c.RemoveId(bson.ObjectIdHex(clockinID))
	if err != nil {
		return err
	}
	return standardizeError(err)
}
