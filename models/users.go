package models

import (
	"errors"

	"github.com/flapan/lenslocked.com/hash"
	"github.com/flapan/lenslocked.com/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned if an invalid id is supplied
	ErrInvalidID = errors.New("models: ID must be > 0")

	// ErrInvalidPassword is returned when a provided password is invalid
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "sGfwegCagsdl3qwY"
const hmacSecretKey = "secret-hmac-key"

// User represents the user model stored in DB
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gotm:"nit null;unique_index"`
}

// UserDB is used to interact with the users database
//
// ByID will look up a user by ID given as input parameter
// 1 - User found -> User, nil
// 2 - User not found -> nil, ErrNotFound
// 3 - Another error -> nil, otherError
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	// Authenticate will verify the provided email and password are correct, if they are correct
	// the user corresponding to those will be returned, otherwise you will receive either
	// ErrNotFound, ErrInvalidPassword or another error if something goes wrong
	Authenticate(email, password string) (*User, error)
	UserDB
}

// Sets up UserService with database connection
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	// db, err := gorm.Open("postgres", connectionInfo)
	// if err != nil {
	// 	return nil, err
	// }
	// db.LogMode(true)
	// hmac := hash.NewHMAC(hmacSecretKey)
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

// This ensures that the type always matches the interface (saves a lot of test lines)
var _ UserService = &userService{}

type userService struct {
	UserDB
}

// Authenticate is used to authenticate a user with the provided email and
// password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

// This ensures that the type always matches the interface (saves a lot of test lines)
var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
}

// This ensures that the type always matches the interface (saves a lot of test lines)
var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail finds a user by email and returns it
// 1 - User, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user by the given token and returns that user.
// This method will handle hashing the token.
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Creates the provided user and backfill data like ID, created_at etc.
// Naive hashing of password, i.e no checking for length etc
func (ug *userGorm) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
}

// Updates a user
func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

//Closes the UserService database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

//DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automaticcaly migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// first will query using the provided gorm.db and it will get the first item returned
// and place it into dst, if nothing is returned is found in the query
// it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
