package passwd

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

type Entries []Entry

type Entry struct {
	Comment      string
	PasswordHash []byte
	Username     string
}

// Add adds a new entry to the specified boltdb database.
// The supplied password is hashed using the bcrypt algorithm before storing
// in the database.
func Add(username, password, comment string, cost int, db *bolt.DB) error {
	entry, err := NewEntry(username, password, comment, cost)
	if err != nil {
		return err
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("htpasswd"))
		return b.Put([]byte(username), data)
	})
	return err
}

func NewEntry(username, password, comment string, cost int) (*Entry, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}
	return &Entry{comment, passwordHash, username}, nil
}

func Delete(username string, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("htpasswd")).Delete([]byte(username))
	})
	return err
}

func List(db *bolt.DB) (Entries, error) {
	entries := make(Entries, 0)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("htpasswd"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e Entry
			if err := json.Unmarshal(v, &e); err != nil {
				return err
			}
			entries = append(entries, e)
		}
		return nil
	})
	return entries, err
}

func Verify(username, password string, db *bolt.DB) error {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("htpasswd"))
		rawEntry := b.Get([]byte(username))

		var e Entry
		if err := json.Unmarshal(rawEntry, &e); err != nil {
			return err
		}

		return bcrypt.CompareHashAndPassword(e.PasswordHash, []byte(password))
	})
	return err
}
