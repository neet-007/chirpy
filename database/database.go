package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/fs"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int
	Body string
}

type User struct {
	Id       int
	Email    string
	Password string
}

type ReturnedUser struct {
	Id    int
	Email string
}

type DBStructure struct {
	Chirps map[int]Chirp   `json:"chirps"`
	Users  map[string]User `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		if err != fs.ErrExist {
			return db, err
		}
	}

	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		Id:   len(dbStructure.Chirps) + 1,
		Body: body,
	}

	dbStructure.Chirps[chirp.Id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, fmt.Errorf("writing db error %w", err)
	}

	return chirp, nil
}

func (db *DB) CreateUser(email string, password string) (ReturnedUser, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return ReturnedUser{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return ReturnedUser{}, err
	}

	user := User{
		Id:       len(dbStructure.Users) + 1,
		Email:    email,
		Password: string(hashedPassword),
	}

	dbStructure.Users[user.Email] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return ReturnedUser{}, fmt.Errorf("writing db error %w", err)
	}

	return ReturnedUser{
		Id:    user.Id,
		Email: user.Email,
	}, nil
}

func (db *DB) GetUser(email string, password string) (ReturnedUser, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return ReturnedUser{}, err
	}

	returnUser, ok := dbStructure.Users[email]

	if !ok {
		return ReturnedUser{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(returnUser.Password), []byte(password))
	if err != nil {
		return ReturnedUser{}, fmt.Errorf("passwords don't match: %v", err)
	}

	return ReturnedUser{
		Id:    returnUser.Id,
		Email: returnUser.Email,
	}, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	returnChirps := []Chirp{}

	for _, chirp := range dbStructure.Chirps {
		returnChirps = append(returnChirps, chirp)
	}

	sort.Slice(returnChirps, func(i, j int) bool { return returnChirps[i].Id < returnChirps[j].Id })

	return returnChirps, nil
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	returnChirp, ok := dbStructure.Chirps[id]

	if !ok {
		return Chirp{}, nil
	}

	return returnChirp, nil

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	return os.WriteFile(db.path, []byte{}, 0666)
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	newData := DBStructure{}

	if len(data) == 0 {
		newData.Chirps = map[int]Chirp{}
		newData.Users = map[string]User{}
		return newData, nil
	}

	err = json.Unmarshal(data, &newData)
	if err != nil {
		return DBStructure{}, err
	}

	return newData, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	fmt.Println(dbStructure)
	json_, err := json.Marshal(dbStructure)
	if err != nil {
		return fmt.Errorf("writing this %v error %w", dbStructure.Chirps, err)
	}

	fmt.Println(json_)
	err = os.WriteFile(db.path, json_, 0666)

	if err != nil {
		return err
	}

	return nil
}
