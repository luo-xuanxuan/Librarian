package user

import (
	"Librarian/utils"
	"database/sql"
	"encoding/json"
)

type User struct {
	ID   string
	Keys map[string]any
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./user/user.db")
	if err != nil {
		utils.Log.Error(err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        data TEXT NOT NULL
    );
    `
	_, err = db.Exec(createTable)
	if err != nil {
		utils.Log.Error(err)
	}
}

func (u *User) Save() error {
	jsonData, err := json.Marshal(u.Keys)
	if err != nil {
		return err
	}

	_, err = db.Exec("REPLACE INTO users (id, data) VALUES (?, ?)", u.ID, jsonData)
	return err
}

func (u *User) Load(ID string) error {
	var jsonData string
	err := db.QueryRow("SELECT data FROM users WHERE id = ?", ID).Scan(&jsonData)

	if err == sql.ErrNoRows {
		// No user found, initialize an empty map for Keys
		u.Keys = make(map[string]interface{})
		u.ID = ID
		return nil // Optionally, you could also return an error or a specific message indicating a new user was created
	} else if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(jsonData), &u.Keys)
	if err != nil {
		return err
	}

	u.ID = ID
	return nil
}

func (u *User) Get(key string) interface{} {
	return u.Keys[key]
}

func (u *User) Set(key string, value interface{}) {
	if u.Keys == nil {
		u.Keys = make(map[string]interface{})
	}
	u.Keys[key] = value
	u.Save()
}
