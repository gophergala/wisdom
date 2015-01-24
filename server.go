package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	PORT         = os.Getenv("PORT")
	DATABASE_URL = os.Getenv("DATABASE_URL")
)

// apiError define structure of API error
type apiError struct {
	Tag     string `json:"-"`
	Error   error  `json:"-"`
	Message string `json:"error"`
	Code    int    `json:"code"`
}

// ApiHandler global API mux
type ApiHandler struct {
	DB      *sql.DB
	Handler func(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError
}

func (api ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// add header on every response
	w.Header().Add("Server", "Automata/0.1")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// if handler return an &apiError
	err := api.Handler(w, r, api.DB)
	if err != nil {
		// http log
		log.Printf("%s %s %s [%s] %s", r.RemoteAddr, r.Method, r.URL, err.Tag, err.Error)

		// response proper http status code
		w.WriteHeader(err.Code)

		// response JSON
		resp := json.NewEncoder(w)
		err_json := resp.Encode(err)
		if err_json != nil {
			log.Println("Encode JSON for error response was failed.")

			return
		}

		return
	}

	// http log
	// TODO: print response
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

// redirect to github pages
func indexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError {
	http.Redirect(w, r, "http://gophergala.github.io/wisdom", 302)
	return nil
}

func main() {
	log.Println("Opening connection to database ... ")
	db, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	// Ping database connection to check connection are OK
	log.Println("Ping database connection ... ")
	err = db.Ping()
	if err != nil {
		log.Println("Ping database connection: failure :(")
		log.Fatal(err)
	}
	log.Println("Ping database connection: success!")

	// index handler
	http.Handle("/", ApiHandler{db, indexHandler})

	// server listener
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
