# passwd

[![GoDoc](https://godoc.org/github.com/kelseyhightower/htbolt/passwd?status.svg)](https://godoc.org/github.com/kelseyhightower/htbolt/passwd)

## Usage

### Create a boltdb password database

```
htbolt -f .boltpasswd -u kelsey
```

### Write some code

```
package main

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/kelseyhightower/htbolt/passwd"
)

func basicAuthHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	db, err := bolt.Open(".boltpasswd", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = passwd.Verify(username, password, db)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
}

func main() {
	http.HandleFunc("/", basicAuthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Use basic auth with cURL.

```
curl -i http://kelsey:password@127.0.0.1:8080
```
