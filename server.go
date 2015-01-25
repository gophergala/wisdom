package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// DatabaseUtils represent database utility that used by handler
type DatabaseUtils struct {
	DB                     *sql.DB
	StatementRandom        *sql.Stmt
	StatementAuthorById    *sql.Stmt
	StatementTagsByQuoteId *sql.Stmt
	StatementTagById       *sql.Stmt
}

// ApiHandler global API mux
type ApiHandler struct {
	DBUtils *DatabaseUtils
	Handler func(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError
}

func (api ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// add header on every response
	w.Header().Add("Server", "Wisdom powered by Gophergala")
	w.Header().Add("X-Wisdom-Media-Type", "wisdom.V1")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// if handler return an &apiError
	err := api.Handler(w, r, api.DBUtils)
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
func indexHandler(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError {

	// response "404 not found" on every undefined
	// URL pattern handler
	if r.URL.Path != "/" {
		return &apiError{
			"indexHandler url",
			errors.New("Not Found"),
			"Not Found",
			http.StatusNotFound,
		}
	}

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

	// index handler doesn't need database utils
	http.Handle("/", ApiHandler{Handler: indexHandler})

	// server listener
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
