package users

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Email     string        `json:"email"`
	PassHash  []byte        `json:"-"` //stored, but not encoded to clients
	UserName  string        `json:"userName"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	PhotoURL  string        `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	if len(nu.UserName) == 0 {
		return errors.New("username must be provided")
	}
	if len(nu.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return errors.New("passwords must match")
	}
	if _, err := mail.ParseAddress(nu.Email); err != nil {
		return fmt.Errorf("invalid email: %v", err)
	}

	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	e := strings.ToLower(strings.TrimSpace(nu.Email))
	h := md5.New()
	h.Write([]byte(e))
	hash := hex.EncodeToString(h.Sum(nil))
	u := &User{
		ID:        bson.NewObjectId(),
		Email:     nu.Email,
		UserName:  nu.UserName,
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		PhotoURL:  gravatarBasePhotoURL + hash,
	}

	if err := u.SetPassword(nu.Password); err != nil {
		return nil, fmt.Errorf("error setting password: %v", err)
	}

	return u, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put betweeen the names
func (u *User) FullName() string {
	fname := u.FirstName
	if len(fname) == 0 {
		return u.LastName
	}
	if len(u.LastName) > 0 {
		fname += " " + u.LastName
	}
	return fname
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	if len(updates.FirstName) == 0 {
		return errors.New("must provide a first name")
	}
	if len(updates.LastName) == 0 {
		return errors.New("must provide a last name")
	}
	u.FirstName = updates.FirstName
	u.LastName = updates.LastName
	return nil
}
