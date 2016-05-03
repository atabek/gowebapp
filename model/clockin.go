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
	InAt         time.Time     `db:"in_at"  bson:"in_at"`
	OutAt        time.Time     `db:"out_at" bson:"out_at"`
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

// ClockinByStudentID gets all clockins for a student
// func ClockinByStudentID(student_id string) (Clockin, error) {
// 	var err error
//
// 	var result Clockin
//     timeZero := time.Time{}
//
// 	if database.CheckConnection() {
// 		// Create a copy of mongo
// 		session := database.Mongo.Copy()
// 		defer session.Close()
// 		c := session.DB(database.ReadConfig().MongoDB.Database).C("clockin")
//
// 		// Validate the object id
//         // if bson.IsObjectIdHex(student_id) {
// 		err = c.Find(bson.M{"student_id": student_id, "out_at": timeZero}).One(&result)
//         fmt.Println("Got the following result: ", result)
// 		// } else {
// 		// 	err = ErrNoResult
// 		// }
// 	} else {
// 		err = ErrUnavailable
// 	}
//
// 	return result, standardizeError(err)
// }

// ClockinCreate creates a clockin
func ClockinCreate(student_id string) error {
	var err error
	now := time.Now()

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
    zeroTime := time.Time{}

    err = s.Find(bson.M{"student_id": student_id}).One(&result)
    if err != nil{
        err = ErrNoSuchStudent
        return err
    }
    colQuerier := bson.M{"student_id": student_id, "out_at": zeroTime}
    count, err := c.Find(colQuerier).Count()
	//fmt.Println("Count: ", count)

    if count == 0 {
        //fmt.Println("This student needs to clock out.")
        clockin := &Clockin{
            ObjectID:  bson.NewObjectId(),
            StudentID: student_id,
            InAt:      now,
            OutAt:     time.Time{},
        }
		//fmt.Println(clockin.ClockinID())
        err = c.Insert(clockin)
    }
    if count == 1 {
		// Error clockin by student id retrieving
        if err != nil{
            return err
        }

		// Update
        change := bson.M{"$set": bson.M{"out_at": time.Now(), "is_out": true}}
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
