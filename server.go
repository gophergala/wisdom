package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	PORT         = os.Getenv("PORT")
	DATABASE_URL = os.Getenv("DATABASE_URL")
)

type Author struct {
	Id        int    `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Twitter   string `json:"twitter_username"`
}

type Tag struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

type Quote struct {
	Id         int    `json:"id"`
	PostId     string `json:"post_id"`
	Author     Author `json:"author"`
	Content    string `json:"content"`
	Permalink  string `json:"permalink"`
	PictureUrl string `json:"picture_url"`
	Tags       []Tag  `json:"tags"`
}

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
	StatementAuthors       *sql.Stmt
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
	http.Redirect(w, r, "http://gophergala.github.io/wisdom", 302)
	return nil
}

// redirect to github pages
func notFoundHandler(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError {
	return &apiError{
		"notFoundHandler",
		errors.New("Not Found"),
		"Not Found",
		http.StatusNotFound,
	}
}

// response random quotes
func randomHandler(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError {
	// get the quote
	var quote_id, quote_author_id int
	var post_id, content, permalink, picture_url string
	err := dbUtils.StatementRandom.QueryRow().Scan(&quote_id, &quote_author_id, &post_id, &content, &permalink, &picture_url)
	if err == sql.ErrNoRows {
		return &apiError{
			Tag:     "quote.noRows",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}
	if err != nil {
		return &apiError{
			Tag:     "quote.err!=nil",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}

	quote := &Quote{
		Id:         quote_id,
		PostId:     post_id,
		Content:    content,
		Permalink:  permalink,
		PictureUrl: picture_url,
	}

	// get the author
	var author Author
	var author_id int
	var avatar_url, name, company_name, twitter_username sql.NullString
	err = dbUtils.StatementAuthorById.QueryRow(quote_author_id).Scan(&author_id, &avatar_url, &name, &company_name, &twitter_username)
	if err == sql.ErrNoRows {
		return &apiError{
			Tag:     "author.noRows",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}
	if err != nil {
		return &apiError{
			Tag:     "author.err!=nil",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}

	author.Id = author_id
	if avatar_url.Valid {
		author.AvatarUrl = avatar_url.String
	} else {
		author.AvatarUrl = ""
	}
	if name.Valid {
		author.Name = name.String
	}
	if company_name.Valid {
		author.Company = company_name.String
	} else {
		author.Company = ""
	}
	if twitter_username.Valid {
		author.Twitter = twitter_username.String
	} else {
		author.Twitter = ""
	}

	quote.Author = author

	// get the tag ids
	var tag_ids []int
	rows, err := dbUtils.StatementTagsByQuoteId.Query(quote_id)
	if err != nil {
		return &apiError{
			Tag:     "tags.rowsErr",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}
	defer rows.Close()
	for rows.Next() {
		var tag_id int
		if err := rows.Scan(&tag_id); err != nil {
			return &apiError{
				Tag:     "tag.rows.Scan",
				Error:   err,
				Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				Code:    http.StatusNoContent,
			}
		}
		tag_ids = append(tag_ids, tag_id)
	}

	// get the tags
	var tags []Tag
	for _, tag_id := range tag_ids {
		var mtag_id int
		var mtag_label string
		var tag Tag
		err := dbUtils.StatementTagById.QueryRow(tag_id).Scan(&mtag_id, &mtag_label)
		if err == sql.ErrNoRows {
			return &apiError{
				Tag:     "mtag.noRows",
				Error:   err,
				Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				Code:    http.StatusNoContent,
			}
		}
		if err != nil {
			return &apiError{
				Tag:     "mtag.err!=nil",
				Error:   err,
				Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				Code:    http.StatusNoContent,
			}
		}
		tag.Id = mtag_id
		tag.Label = mtag_label
		tags = append(tags, tag)
	}
	quote.Tags = tags

	// write json to response
	// response JSON
	randomResp := json.NewEncoder(w)
	random_err_json := randomResp.Encode(quote)
	if random_err_json != nil {
		log.Println("Encode JSON for error response was failed.")
		return &apiError{
			Tag:     "random_err_json.err!=nil",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}
	return nil
}

// /v1/authors endpoint. return an array of authors
func authorsHandler(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError {
	// get the rows
	var authors []Author
	authorsRows, err := dbUtils.StatementAuthors.Query()
	if err != nil {
		return &apiError{
			Tag:     "authorRows.rowsErr",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}
	defer authorsRows.Close()
	for authorsRows.Next() {
		var author Author
		var author_id int
		var avatar_url, name, company_name, twitter_username sql.NullString
		if err := authorsRows.Scan(&author_id, &avatar_url, &name, &company_name, &twitter_username); err != nil {
			return &apiError{
				Tag:     "authorRows.Scan",
				Error:   err,
				Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				Code:    http.StatusNoContent,
			}
		}
		author.Id = author_id
		if avatar_url.Valid {
			author.AvatarUrl = avatar_url.String
		} else {
			author.AvatarUrl = ""
		}
		if name.Valid {
			author.Name = name.String
		}
		if company_name.Valid {
			author.Company = company_name.String
		} else {
			author.Company = ""
		}
		if twitter_username.Valid {
			author.Twitter = twitter_username.String
		} else {
			author.Twitter = ""
		}
		authors = append(authors, author)
	}

	// response json
	authorsResp := json.NewEncoder(w)
	authors_err_json := authorsResp.Encode(authors)
	if authors_err_json != nil {
		log.Println("Encode JSON for error response was failed.")
		return &apiError{
			Tag:     "authors_err_json.err!=nil",
			Error:   err,
			Message: "OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			Code:    http.StatusNoContent,
		}
	}
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

	r := mux.NewRouter()
	// index handler doesn't need database utils
	r.Handle("/", ApiHandler{Handler: indexHandler})

	// Random handler
	// prepare a statement
	stmtQueryRandomQuote, err := db.Prepare("SELECT * FROM quotes ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		log.Println(err)
	}

	stmtQueryAuthor, err := db.Prepare("SELECT * FROM authors WHERE id = $1;")
	if err != nil {
		log.Println(err)
	}

	stmtQueryTagsByQuoteId, err := db.Prepare("SELECT tag_id FROM quotes_tags WHERE quote_id = $1")
	if err != nil {
		log.Println(err)
	}

	stmtQueryTagById, err := db.Prepare("SELECT * FROM tags WHERE id = $1")
	if err != nil {
		log.Println(err)
	}

	randomDBUtils := &DatabaseUtils{
		StatementRandom:        stmtQueryRandomQuote,
		StatementAuthorById:    stmtQueryAuthor,
		StatementTagsByQuoteId: stmtQueryTagsByQuoteId,
		StatementTagById:       stmtQueryTagById,
	}
	r.Handle("/v1/random", ApiHandler{randomDBUtils, randomHandler})

	// Authors handler
	stmtQueryAuthors, err := db.Prepare("SELECT * FROM authors")
	if err != nil {
		log.Println(err)
	}

	authorsDBUtils := &DatabaseUtils{
		StatementAuthors: stmtQueryAuthors,
	}

	r.Handle("/v1/authors", ApiHandler{authorsDBUtils, authorsHandler})

	// not found handler
	r.NotFoundHandler = ApiHandler{Handler: notFoundHandler}
	// server listener
	http.Handle("/", r)
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
