package handlers

import (
	"BookManager/db"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type author struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Birthday    string `json:"birthday"`
	Nationality string `json:"nationality"`
}

type bookBody struct {
	Name            string   `json:"name"`
	Author          author   `json:"author"`
	Category        string   `json:"category"`
	Volume          int      `json:"volume"`
	PublishedAt     string   `json:"published_at"`
	Summary         string   `json:"summary"`
	TableOfContents []string `json:"table_of_contents"`
	Publisher       string   `json:"publisher"`
}

type getAllBooksInfo struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Volume      int    `json:"volume"`
	PublishedAt string `json:"published_at"`
	Summary     string `json:"summary"`
	Publisher   string `json:"publisher"`
}

// HandleCreateGetAllBooks handles creating a book and getting all books
func (bm *BookManagerServer) HandleCreateGetAllBooks(w http.ResponseWriter, r *http.Request) {
	// Grab Authorization Header
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve the related account to the token
	username, err := bm.Authenticate.GetUsernameByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Handle create book & get all books based on given method
	if r.Method == http.MethodPost {
		// Handle create book
		// Parse the request body for the new book
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			bm.Logger.WithError(err).Warn("can not read the request data")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var BD bookBody
		err = json.Unmarshal(reqData, &BD)
		if err != nil {
			bm.Logger.WithError(err).Warn("can not unmarshal the new book request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Add the book to the database
		book := &db.Book{
			Name:              BD.Name,
			AuthorFirstName:   BD.Author.FirstName,
			AuthorLastName:    BD.Author.LastName,
			AuthorBirthday:    BD.Author.Birthday,
			AuthorNationality: BD.Author.Nationality,
			Category:          BD.Category,
			Volume:            BD.Volume,
			PublishedAt:       BD.PublishedAt,
			Summary:           BD.Summary,
			TableOfContents:   BD.TableOfContents,
			Publisher:         BD.Publisher,
			UserOwner:         *username,
		}

		err = bm.Db.AddNewBook(book)
		if err != nil {
			bm.Logger.WithError(err).Warn("can not add the new book")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response := map[string]interface{}{
			"message": "book has been added successfully",
		}

		resBody, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resBody)

	} else if r.Method == http.MethodGet {
		// Handle get all books
		// Retrieve books from database
		books, err := bm.Db.GetBooksByUsername(*username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Keep just necessary data
		var responseBooks []getAllBooksInfo

		for _, book := range *books {
			author := book.AuthorFirstName + " " + book.AuthorLastName

			// Create an getAllBooksInfo object
			BI := getAllBooksInfo{
				Name:        book.Name,
				Author:      author,
				Category:    book.Category,
				Volume:      book.Volume,
				PublishedAt: book.PublishedAt,
				Summary:     book.Summary,
				Publisher:   book.Publisher,
			}

			// Append the created book object to the responseBooks slice
			responseBooks = append(responseBooks, BI)
		}

		// Create the response body
		res, err := json.Marshal(responseBooks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// HandleUpdateRetrieveDeleteBook is responsible for Update, Retrieve (get book by id) and Delete
func (bm *BookManagerServer) HandleUpdateRetrieveDeleteBook(w http.ResponseWriter, r *http.Request) {
	// Grab Authorization Header
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve the related account to the token
	username, err := bm.Authenticate.GetUsernameByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	bookID := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(bookID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		// Handle get book by id
		// Retrieve book by id
		book, err := bm.Db.GetBookByID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Assure that just the book owner can access required book
		if book.UserOwner != *username {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Create the response body
		res, err := json.Marshal(&bookBody{
			Name: book.Name,
			Author: author{
				FirstName:   book.AuthorFirstName,
				LastName:    book.AuthorLastName,
				Birthday:    book.AuthorBirthday,
				Nationality: book.AuthorNationality,
			},
			Category:        book.Category,
			Volume:          book.Volume,
			PublishedAt:     book.PublishedAt,
			Summary:         book.Summary,
			TableOfContents: book.TableOfContents,
			Publisher:       book.Publisher,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	} else if r.Method == http.MethodDelete {
		// Handle Delete book by id
		err := bm.Db.DeleteBookByID(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create the response body
		res := map[string]interface{}{
			"message": "book has been deleted successfully",
		}
		resBody, err := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)

	} else if r.Method == http.MethodPut {
		// Handle Update book by id
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			bm.Logger.WithError(err).Warn("can not read the request data")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var BD bookBody
		err = json.Unmarshal(reqData, &BD)
		if err != nil {
			bm.Logger.WithError(err).Warn("can not unmarshal the update book request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Fetch the book by ID
		book, err := bm.Db.GetBookByID(id)
		if err != nil {
			bm.Logger.WithError(err).Warn("failed to fetch the book")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Assure that just the book owner can access required book
		if book.UserOwner != *username {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Find parameters that are requested to update
		if BD.Name != "" {
			book.Name = BD.Name
		}

		if BD.Author.FirstName != "" {
			book.AuthorFirstName = BD.Author.FirstName
		}

		if BD.Author.LastName != "" {
			book.AuthorLastName = BD.Author.LastName
		}

		if BD.Author.Birthday != "" {
			book.AuthorBirthday = BD.Author.Birthday
		}

		if BD.Author.Nationality != "" {
			book.AuthorNationality = BD.Author.Nationality
		}

		if BD.Category != "" {
			book.Category = BD.Category
		}

		if BD.Volume != 0 {
			book.Volume = BD.Volume
		}

		if BD.PublishedAt != "" {
			book.PublishedAt = BD.PublishedAt
		}

		if BD.Summary != "" {
			book.Summary = BD.Summary
		}

		if len(BD.TableOfContents) > 0 {
			book.TableOfContents = BD.TableOfContents
		}

		if BD.Publisher != "" {
			book.Publisher = BD.Publisher
		}

		err = bm.Db.UpdateBook(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create the response body
		resBody := map[string]interface{}{
			"message": "book has been updated successfully",
		}
		res, err := json.Marshal(resBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
