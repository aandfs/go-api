package bookcontroller

import (
	"encoding/json"
	"errors"
	"go-api/config"
	"go-api/helper"
	"go-api/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var books []models.Book
	var booksResponse []models.BookResponse

	if err := config.DB.Joins("Author").Find(&books).Find(&booksResponse).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "List Books", booksResponse)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		helper.Response(w, 500, err.Error(), nil)
	}

	defer r.Body.Close()

	//check author
	var author models.Author
	if err := config.DB.First(&author, book.AuthorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Author not found", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if err := config.DB.Create(&book).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Success create book", nil)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var book models.Book
	var BookResponse models.BookResponse

	if err := config.DB.Joins("Author").First(&book, id).First(&BookResponse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Book not found", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Detail Book", BookResponse)
}

func Update(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var book models.Book

	if err := config.DB.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Book not found", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	var bookPayload models.Book
	if err := json.NewDecoder(r.Body).Decode(&bookPayload); err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	defer r.Body.Close()

	var author models.Author
	if bookPayload.AuthorID != 0 {
		if err := config.DB.First(&author, bookPayload.AuthorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				helper.Response(w, 404, "Author not found", nil)
				return
			}

			helper.Response(w, 500, err.Error(), nil)
			return
		}
	}

	if err := config.DB.Where("id = ?", id).Updates(&bookPayload).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "success update book", nil)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var book models.Book
	res := config.DB.Delete(&book, id)

	if res.Error != nil {
		helper.Response(w, 500, res.Error.Error(), nil)
		return
	}

	if res.RowsAffected == 0 {
		helper.Response(w, 404, "Book not found", nil)
		return
	}

	helper.Response(w, 200, "sucess delete book", nil)
}
