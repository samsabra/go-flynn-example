package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flynn-examples/go-flynn-example/Godeps/_workspace/src/github.com/flynn/flynn/pkg/postgres"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	db, err := postgres.Open("", "")
	if err != nil {
		log.Fatal(err)
	}

	m := postgres.NewMigrations()
	m.Add(1, "CREATE SEQUENCE hits")
	if err := m.Migrate(db.DB); err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("SELECT nextval('hits')")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var count int
		if err := stmt.QueryRow().Scan(&count); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Fprintf(w, "Hello from Flynn on port %s from container %s\nHits = %d!\n", port, os.Getenv("HOSTNAME"), count)
	})
	fmt.Println("hitcounter listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
