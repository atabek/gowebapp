package model

import (
	"time"

	"github.com/atabek/gowebapp/shared/database"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Student
// *****************************************************************************

// Student table contains the information for each student
type Student struct {
	ObjectId   bson.ObjectId `bson:"_id"`
	First_name string        `db:"first_name" bson:"first_name"`
	Last_name  string        `db:"last_name"  bson:"last_name"`
	Grade      string        `db:"grade"      bson:"grade"`
	Student_id string        `db:"student_id" bson:"student_id"`
	Created_at time.Time     `db:"created_at" bson:"created_at"`
	Updated_at time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted    uint8         `db:"deleted"    bson:"deleted"`
}

// Id returns the student id
func (s *Student) ID() string {
	return s.ObjectId.Hex()
}

// SID returns the student school id
func (s *Student) SID() string {
	return s.Student_id
}

// StudentBySID gets student information from student school id
func StudentBySID(sid string) (Student, error) {
	var err error

	result := Student{}

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("student")
		err = c.Find(bson.M{"student_id": sid}).One(&result)
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// StudentCreate creates student
func StudentCreate(first_name, last_name, grade, student_id string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("student")

		student := &Student{
			ObjectId:   bson.NewObjectId(),
			First_name: first_name,
			Last_name:  last_name,
			Grade:      grade,
			Student_id: student_id,
			Created_at: now,
			Updated_at: now,
			Deleted:    0,
		}
		err = c.Insert(student)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// StudentUpdate updates a student
func StudentUpdate(studentID, firstName, lastName, grade string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("student")
		var student Student
		student, err = StudentBySID(studentID)
		if err == nil {
			student.Updated_at = now
			student.First_name = firstName
			student.Last_name  = lastName
			student.Grade      = grade
			err = c.UpdateId(bson.ObjectIdHex(student.ID()), &student)
		} else {
			err = ErrUnauthorized
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// StudentsGet gets students
func StudentsGet() ([]Student, error) {
	var err error

	// List all students
	var students = make([]Student, 0)

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("student")
		err = c.Find(nil).All(&students)
	} else {
		err = ErrUnavailable
	}

	return students, standardizeError(err)
}
