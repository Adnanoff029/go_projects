package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Adnanoff029/bookstore/pkg/models"
	"github.com/Adnanoff029/bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

var NewBook models.Book

func GetBook(w http.ResponseWriter, r *http.Request) {
	AllBooks := models.GetAllBooks()
	res, _ := json.Marshal(AllBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	Id, err := strconv.ParseInt(params["bookId"], 0, 0)
	if err != nil {
		fmt.Println("Error getting the ID.")
	}
	FetchBook, _ := models.GetBookById(Id)
	res, _ := json.Marshal(FetchBook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	BookModel := &models.Book{}
	utils.ParseBody(r, BookModel)
	CreatedBook := BookModel.CreateBook()
	res, _ := json.Marshal(CreatedBook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	Id, err := strconv.ParseInt(params["bookId"], 0, 0)
	if err != nil {
		fmt.Println("Error in getting the Id.")
	}
	DeletedBook := models.DeleteBookByID(Id)
	res, _ := json.Marshal(DeletedBook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var UpdateBook = &models.Book{}
	utils.ParseBody(r, UpdateBook)
	params := mux.Vars(r)
	Id, err := strconv.ParseInt(params["bookId"], 0, 0)
	if err != nil {
		fmt.Println("Error in getting the Id.")
	}
	BookDetails, db := models.GetBookById(Id)
	if UpdateBook.Name != "" {
		BookDetails.Name = UpdateBook.Name
	}
	if UpdateBook.Author != "" {
		BookDetails.Author = UpdateBook.Author
	}
	if UpdateBook.Publication != "" {
		BookDetails.Publication = UpdateBook.Publication
	}
	db.Save(&BookDetails)

	res, _ := json.Marshal(BookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
