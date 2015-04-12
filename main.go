package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/htbolt/passwd"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	comment  string
	cost     int
	database string
	delete   bool
	dryrun   bool
	list     bool
	password string
	username string
	verify   bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&comment, "c", "", "the comment field")
	flag.IntVar(&cost, "C", 10, "computing time used for the bcrypt algorithm")
	flag.StringVar(&database, "f", "", "The path to the boltdb file")
	flag.BoolVar(&delete, "x", false, "If the username exists in the specified boltdb file, it will be deleted.")
	flag.BoolVar(&dryrun, "n", false, "Display the results on standard output rather than updating a database.")
	flag.BoolVar(&list, "l", false, "Print each of the usernames and comments from the database on stdout.")
	flag.StringVar(&password, "p", "", "The password")
	flag.StringVar(&username, "u", "", "The username")
	flag.BoolVar(&verify, "v", false, "Verify the username and password.")
}

func main() {
	flag.Parse()

	var err error

	if !list {
		if username == "" {
			log.Fatal("non-empty username is required")
		}
	}

	if !list && !delete {
		if password == "" {
			password, err = getPassword()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if database == "" {
		log.Fatal("a boltdb database is required.")
	}

	if dryrun {
		entry, err := passwd.NewEntry(username, password, comment, cost)
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(data))
		os.Exit(0)
	}

	db, err := bolt.Open(database, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("htpasswd"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if list {
		entries, err := passwd.List(db)
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range entries {
			log.Printf("%s # %s\n", e.Username, e.Comment)
		}
		os.Exit(0)
	}

	if delete {
		if err := passwd.Delete(username, db); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if verify {
		if err := passwd.Verify(username, password, db); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if err := passwd.Add(username, password, comment, cost, db); err != nil {
		log.Fatal(err)
	}
}

func getPassword() (string, error) {
	fmt.Printf("password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println("")
	return string(password), nil
}
