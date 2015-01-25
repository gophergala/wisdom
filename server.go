package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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
	DB                               *sql.DB
	StatementRandom                  *sql.Stmt
	StatementAuthorById              *sql.Stmt
	StatementTagIdsByQuoteId         *sql.Stmt
	StatementTagById                 *sql.Stmt
	StatementAuthors                 *sql.Stmt
	StatementAuthorByTwitterUsername *sql.Stmt
	StatementQuotesByAuthorId        *sql.Stmt
	StatementTags                    *sql.Stmt
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
			"randomHandler.StatementRandom.QueryRow.sql.ErrNoRows",
			err,
			"Quote not found",
			http.StatusNotFound,
		}
	}
	if err != nil {
		return &apiError{
			"randomHandler.StatementRandom.QueryRow.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
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
			"randomHandler.StatementAuthorById.QueryRow.sql.ErrNoRows",
			err,
			"Author not found",
			http.StatusNotFound,
		}
	}
	if err != nil {
		return &apiError{
			"randomHandler.StatementAuthorById.QueryRow.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}

	author.Id = author_id
	author.Name = name.String
	if avatar_url.Valid {
		author.AvatarUrl = avatar_url.String
	} else {
		author.AvatarUrl = ""
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
	tagIdsRows, err := dbUtils.StatementTagIdsByQuoteId.Query(quote_id)
	if err != nil {
		return &apiError{
			"randomHandler.tagIdsRows.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	defer tagIdsRows.Close()
	for tagIdsRows.Next() {
		var tag_id int
		if err := tagIdsRows.Scan(&tag_id); err != nil {
			return &apiError{
				"randomHandler.tagIdsRows.Scan.Err",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
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
				"randomHandler.StatementTagById.QueryRow.sql.ErrNoRows",
				err,
				"Tag not found",
				http.StatusNotFound,
			}
		}
		if err != nil {
			return &apiError{
				"randomHandler.StatementTagById.QueryRow.Err",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
			}
		}
		tag.Id = mtag_id
		tag.Label = mtag_label
		tags = append(tags, tag)
	}
	quote.Tags = tags

	// JSONP response
	query := r.URL.Query()
	jsonp := query.Get("jsonp")
	callback := query.Get("callback")
	if callback != "" || jsonp != "" {
		jsonResult, err := json.Marshal(quote)
		if err != nil {
			log.Println(err)
		}
		if callback != "" {
			fmt.Fprintf(w, "%s(%s)", callback, jsonResult)
			return nil
		}

		if jsonp != "" {
			fmt.Fprintf(w, "%s(%s)", jsonp, jsonResult)
			return nil
		}
	}

	// write json to response
	// response JSON
	randomResp := json.NewEncoder(w)
	err = randomResp.Encode(quote)
	if err != nil {
		return &apiError{
			"randomHandler.randomResp.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
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
			"authorsHandler.authorsRows.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	defer authorsRows.Close()
	for authorsRows.Next() {
		var author Author
		var author_id int
		var avatar_url, name, company_name, twitter_username sql.NullString
		if err := authorsRows.Scan(&author_id, &avatar_url, &name, &company_name, &twitter_username); err != nil {
			return &apiError{
				"authorsHandler.authorsRows.Scan.Err",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
			}
		}
		author.Id = author_id
		author.Name = name.String
		if avatar_url.Valid {
			author.AvatarUrl = avatar_url.String
		} else {
			author.AvatarUrl = ""
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

	// JSONP response
	query := r.URL.Query()
	jsonp := query.Get("jsonp")
	callback := query.Get("callback")
	if callback != "" || jsonp != "" {
		jsonResult, err := json.Marshal(authors)
		if err != nil {
			log.Println(err)
		}
		if callback != "" {
			fmt.Fprintf(w, "%s(%s)", callback, jsonResult)
			return nil
		}

		if jsonp != "" {
			fmt.Fprintf(w, "%s(%s)", jsonp, jsonResult)
			return nil
		}
	}

	// response json
	authorsResp := json.NewEncoder(w)
	err = authorsResp.Encode(authors)
	if err != nil {
		return &apiError{
			"authorsHandler.authorsResp.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	return nil
}

func authorTwitterHandler(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError {
	// get the parameter
	vars := mux.Vars(r)
	twitter_username := vars["twitter_username"]

	// get the author
	var author Author
	var author_id int
	var avatar_url, name, company_name, author_twitter_username sql.NullString
	err := dbUtils.StatementAuthorByTwitterUsername.QueryRow(twitter_username).Scan(&author_id, &avatar_url, &name, &company_name, &author_twitter_username)
	if err == sql.ErrNoRows {
		return &apiError{
			"authorTwitterHandler.sql.ErrNoRows",
			err,
			"Author not found",
			http.StatusNotFound,
		}
	}
	if err != nil {
		return &apiError{
			"authorTwitterHandler.err!=nil",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}

	author.Id = author_id
	author.Name = name.String
	if avatar_url.Valid {
		author.AvatarUrl = avatar_url.String
	} else {
		author.AvatarUrl = ""
	}
	if company_name.Valid {
		author.Company = company_name.String
	} else {
		author.Company = ""
	}
	if author_twitter_username.Valid {
		author.Twitter = author_twitter_username.String
	} else {
		author.Twitter = ""
	}

	// get the quotes
	var quotes []*Quote
	quotesRows, err := dbUtils.StatementQuotesByAuthorId.Query(author_id)
	if err != nil {
		return &apiError{
			"authorTwitterHandler.quotesRows.err!=nil",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	defer quotesRows.Close()
	for quotesRows.Next() {
		// get the quote
		var quote_id, quote_author_id int
		var post_id, content, permalink, picture_url string
		if err := quotesRows.Scan(&quote_id, &quote_author_id, &post_id, &content, &permalink, &picture_url); err != nil {
			return &apiError{
				"authorTwitterHandler.quotesRows.Scan",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
			}
		}

		quote := &Quote{
			Id:         quote_id,
			PostId:     post_id,
			Content:    content,
			Permalink:  permalink,
			PictureUrl: picture_url,
		}

		// get the tag ids
		var tag_ids []int
		tagIdsRows, err := dbUtils.StatementTagIdsByQuoteId.Query(quote.Id)
		if err != nil {
			return &apiError{
				"authorTwitterHandler.tagIdsRows.err!=nil",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
			}
		}
		defer tagIdsRows.Close()
		for tagIdsRows.Next() {
			var tag_id int
			if err := tagIdsRows.Scan(&tag_id); err != nil {
				return &apiError{
					"authorTwitterHandler.tagIdsRows.Scan",
					err,
					"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
					http.StatusInternalServerError,
				}
			}
			tag_ids = append(tag_ids, tag_id)
		}

		// get the tags
		var tags []Tag
		for _, tag_id := range tag_ids {
			var tag Tag
			var mtag_id int
			var mtag_label string
			err := dbUtils.StatementTagById.QueryRow(tag_id).Scan(&mtag_id, &mtag_label)
			if err == sql.ErrNoRows {
				return &apiError{
					"authorTwitterHandler.tagIdsRows.StatementTagById.sql.ErrNoRows",
					err,
					"Author not found",
					http.StatusNotFound,
				}
			}
			if err != nil {
				return &apiError{
					"authorTwitterHandler.tagIdsRows.StatementTagById.Err",
					err,
					"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
					http.StatusInternalServerError,
				}
			}
			tag.Id = mtag_id
			tag.Label = mtag_label
			tags = append(tags, tag)
		}

		quote.Tags = tags
		quote.Author = author
		quotes = append(quotes, quote)
	}

	// JSONP response
	query := r.URL.Query()
	jsonp := query.Get("jsonp")
	callback := query.Get("callback")
	if callback != "" || jsonp != "" {
		jsonResult, err := json.Marshal(quotes)
		if err != nil {
			log.Println(err)
		}
		if callback != "" {
			fmt.Fprintf(w, "%s(%s)", callback, jsonResult)
			return nil
		}

		if jsonp != "" {
			fmt.Fprintf(w, "%s(%s)", jsonp, jsonResult)
			return nil
		}
	}

	// response JSON
	quotesResp := json.NewEncoder(w)
	err = quotesResp.Encode(quotes)
	if err != nil {
		return &apiError{
			"authorTwitterHandler.quotesResp.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	return nil
}

func authorTwitterRandomHandler(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError {
	// get the parameter
	vars := mux.Vars(r)
	twitter_username := vars["twitter_username"]

	// get the author
	var author Author
	var author_id int
	var avatar_url, name, company_name, author_twitter_username sql.NullString
	err := dbUtils.StatementAuthorByTwitterUsername.QueryRow(twitter_username).Scan(&author_id, &avatar_url, &name, &company_name, &author_twitter_username)
	if err == sql.ErrNoRows {
		return &apiError{
			"authorTwitterRandomHandler.sql.ErrNoRows",
			err,
			"Author not found",
			http.StatusNotFound,
		}
	}
	if err != nil {
		return &apiError{
			"authorTwitterRandomHandler.err!=nil",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}

	author.Id = author_id
	author.Name = name.String
	if avatar_url.Valid {
		author.AvatarUrl = avatar_url.String
	} else {
		author.AvatarUrl = ""
	}
	if company_name.Valid {
		author.Company = company_name.String
	} else {
		author.Company = ""
	}
	if author_twitter_username.Valid {
		author.Twitter = author_twitter_username.String
	} else {
		author.Twitter = ""
	}

	// get the quotes
	var quotes []*Quote
	quotesRows, err := dbUtils.StatementQuotesByAuthorId.Query(author_id)
	if err != nil {
		return &apiError{
			"authorTwitterRandomHandler.quotesRows.err!=nil",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	defer quotesRows.Close()
	for quotesRows.Next() {
		// get the quote
		var quote_id, quote_author_id int
		var post_id, content, permalink, picture_url string
		if err := quotesRows.Scan(&quote_id, &quote_author_id, &post_id, &content, &permalink, &picture_url); err != nil {
			return &apiError{
				"authorTwitterRandomHandler.quotesRows.Scan",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
			}
		}

		quote := &Quote{
			Id:         quote_id,
			PostId:     post_id,
			Content:    content,
			Permalink:  permalink,
			PictureUrl: picture_url,
		}

		// get the tag ids
		var tag_ids []int
		tagIdsRows, err := dbUtils.StatementTagIdsByQuoteId.Query(quote.Id)
		if err != nil {
			return &apiError{
				"authorTwitterRandomHandler.tagIdsRows.err!=nil",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
			}
		}
		defer tagIdsRows.Close()
		for tagIdsRows.Next() {
			var tag_id int
			if err := tagIdsRows.Scan(&tag_id); err != nil {
				return &apiError{
					"authorTwitterRandomHandler.tagIdsRows.Scan",
					err,
					"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
					http.StatusInternalServerError,
				}
			}
			tag_ids = append(tag_ids, tag_id)
		}

		// get the tags
		var tags []Tag
		for _, tag_id := range tag_ids {
			var tag Tag
			var mtag_id int
			var mtag_label string
			err := dbUtils.StatementTagById.QueryRow(tag_id).Scan(&mtag_id, &mtag_label)
			if err == sql.ErrNoRows {
				return &apiError{
					"authorTwitterRandomHandler.tagIdsRows.StatementTagById.sql.ErrNoRows",
					err,
					"Author not found",
					http.StatusNotFound,
				}
			}
			if err != nil {
				return &apiError{
					"authorTwitterRandomHandler.tagIdsRows.StatementTagById.Err",
					err,
					"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
					http.StatusInternalServerError,
				}
			}
			tag.Id = mtag_id
			tag.Label = mtag_label
			tags = append(tags, tag)
		}

		quote.Tags = tags
		quote.Author = author
		quotes = append(quotes, quote)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	random := rand.Intn(len(quotes))

	// JSONP response
	query := r.URL.Query()
	jsonp := query.Get("jsonp")
	callback := query.Get("callback")
	if callback != "" || jsonp != "" {
		jsonResult, err := json.Marshal(quotes[random])
		if err != nil {
			log.Println(err)
		}
		if callback != "" {
			fmt.Fprintf(w, "%s(%s)", callback, jsonResult)
			return nil
		}

		if jsonp != "" {
			fmt.Fprintf(w, "%s(%s)", jsonp, jsonResult)
			return nil
		}
	}

	// response JSON
	quotesResp := json.NewEncoder(w)
	err = quotesResp.Encode(quotes[random])
	if err != nil {
		return &apiError{
			"authorTwitterRandomHandler.quotesResp.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	return nil
}

// tags handler
func tagsHandler(w http.ResponseWriter, r *http.Request, dbUtils *DatabaseUtils) *apiError {
	// get the tags
	var tags []Tag
	tagsRows, err := dbUtils.StatementTags.Query()
	if err != nil {
		return &apiError{
			"tagsHandler.tagsRows.err!=nil",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
		}
	}
	defer tagsRows.Close()
	for tagsRows.Next() {
		var tag Tag
		var tag_id int
		var tag_label string
		if err := tagsRows.Scan(&tag_id, &tag_label); err != nil {
			return &apiError{
				"tagsHandler.tagsRows.Scan",
				err,
				"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
				http.StatusInternalServerError,
			}
		}
		tag.Id = tag_id
		tag.Label = tag_label
		tags = append(tags, tag)
	}

	// JSONP response
	query := r.URL.Query()
	jsonp := query.Get("jsonp")
	callback := query.Get("callback")
	if callback != "" || jsonp != "" {
		jsonResult, err := json.Marshal(tags)
		if err != nil {
			log.Println(err)
		}
		if callback != "" {
			fmt.Fprintf(w, "%s(%s)", callback, jsonResult)
			return nil
		}

		if jsonp != "" {
			fmt.Fprintf(w, "%s(%s)", jsonp, jsonResult)
			return nil
		}
	}

	// response JSON
	tagsResp := json.NewEncoder(w)
	err = tagsResp.Encode(tags)
	if err != nil {
		return &apiError{
			"tagsHandler.tagsResp.Err",
			err,
			"OOOOOPPPSSSS! error happen. don't panic! we will be back soon :)",
			http.StatusInternalServerError,
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

	stmtQueryTagIdsByQuoteId, err := db.Prepare("SELECT tag_id FROM quotes_tags WHERE quote_id = $1")
	if err != nil {
		log.Println(err)
	}

	stmtQueryTagById, err := db.Prepare("SELECT * FROM tags WHERE id = $1")
	if err != nil {
		log.Println(err)
	}

	randomDBUtils := &DatabaseUtils{
		StatementRandom:          stmtQueryRandomQuote,
		StatementAuthorById:      stmtQueryAuthor,
		StatementTagIdsByQuoteId: stmtQueryTagIdsByQuoteId,
		StatementTagById:         stmtQueryTagById,
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

	// /v1/author/twitter_username handler
	stmtQueryAuthorByTwitterUsername, err := db.Prepare("SELECT * FROM authors WHERE twitter_username = $1")
	if err != nil {
		log.Println(err)
	}

	stmtQueryQuotesByAuthorId, err := db.Prepare("SELECT * FROM quotes WHERE author_id = $1")
	if err != nil {
		log.Println(err)
	}

	authorTwitterDBUtils := &DatabaseUtils{
		StatementAuthorByTwitterUsername: stmtQueryAuthorByTwitterUsername,
		StatementQuotesByAuthorId:        stmtQueryQuotesByAuthorId,
		StatementTagIdsByQuoteId:         stmtQueryTagIdsByQuoteId,
		StatementTagById:                 stmtQueryTagById,
	}
	r.Handle("/v1/author/{twitter_username}", ApiHandler{authorTwitterDBUtils, authorTwitterHandler})

	authorTwitterRandomDBUtils := &DatabaseUtils{
		StatementAuthorByTwitterUsername: stmtQueryAuthorByTwitterUsername,
		StatementQuotesByAuthorId:        stmtQueryQuotesByAuthorId,
		StatementTagIdsByQuoteId:         stmtQueryTagIdsByQuoteId,
		StatementTagById:                 stmtQueryTagById,
	}
	r.Handle("/v1/author/{twitter_username}/random", ApiHandler{authorTwitterRandomDBUtils, authorTwitterRandomHandler})

	// tags handler
	stmtQueryTags, err := db.Prepare("SELECT * FROM tags")
	if err != nil {
		log.Println(err)
	}

	tagsDBUtils := &DatabaseUtils{
		StatementTags: stmtQueryTags,
	}
	r.Handle("/v1/tags", ApiHandler{tagsDBUtils, tagsHandler})

	// not found handler
	r.NotFoundHandler = ApiHandler{Handler: notFoundHandler}
	// server listener
	http.Handle("/", r)
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
