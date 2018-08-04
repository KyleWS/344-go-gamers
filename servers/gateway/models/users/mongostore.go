package users

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStore is a structure for a mongoDB connection
type MongoStore struct {
	session *mgo.Session
	dbname  string
	colname string
}

//NewMongoStore constructs a new MongoStore
func NewMongoStore(sess *mgo.Session, dbName string, collectionName string) *MongoStore {
	if sess == nil {
		panic("nil pointer passed for session")
	}
	return &MongoStore{
		session: sess,
		dbname:  dbName,
		colname: collectionName,
	}
}

//GetByID returns the User with the given ID
func (s *MongoStore) GetByID(id bson.ObjectId) (*User, error) {
	col := s.session.DB(s.dbname).C(s.colname)
	u := &User{}
	if err := col.FindId(id).One(u); err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

//GetByEmail returns the User with the given email
func (s *MongoStore) GetByEmail(email string) (*User, error) {
	col := s.session.DB(s.dbname).C(s.colname)
	u := &User{}
	if err := col.Find(bson.M{"email": email}).One(u); err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

//GetByUserName returns the User with the given Username
func (s *MongoStore) GetByUserName(username string) (*User, error) {
	col := s.session.DB(s.dbname).C(s.colname)
	u := &User{}
	if err := col.Find(bson.M{"username": username}).One(u); err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

//Insert converts the NewUser to a User, inserts
//it into the database, and returns it
func (s *MongoStore) Insert(newUser *NewUser) (*User, error) {
	u, err := newUser.ToUser()
	if err != nil {
		return nil, err
	}
	col := s.session.DB(s.dbname).C(s.colname)
	if err := col.Insert(u); err != nil {
		return nil, fmt.Errorf("error inserting user: %v", err)
	}
	return u, nil
}

//Update applies UserUpdates to the given user ID
func (s *MongoStore) Update(userID bson.ObjectId, updates *Updates) error {
	// get user to apply updates to
	u, err := s.GetByID(userID)
	if err != nil {
		return err
	}
	// apply updates to user to ensure they are valid
	if err := u.ApplyUpdates(updates); err != nil {
		return err
	}

	col := s.session.DB(s.dbname).C(s.colname)
	// update user fields
	return col.Update(bson.M{"_id": userID}, bson.M{"$set": updates})
}

//Delete deletes the user with the given ID
func (s *MongoStore) Delete(userID bson.ObjectId) error {
	return s.session.DB(s.dbname).C(s.colname).Remove(bson.M{"_id": userID})
}

// IdsToUsers takes a slice of ids and returns a slice of user pointers
func (s *MongoStore) IdsToUsers(bids []bson.ObjectId) []*User {
	res := make([]*User, 0, len(bids))
	for _, bid := range bids {
		u, err := s.GetByID(bid)
		// no need to fail out if user lookup fails for some reason
		// just write it to log
		if err != nil {
			log.Printf("Got at error while getting user from bid slice: %v\n", err)
		}
		res = append(res, u)
	}
	return res
}
