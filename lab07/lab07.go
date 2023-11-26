package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Pages int    `json:"pages"`
}

var bookshelf = []*Book{
	{ID: 1, Name: "Blue Bird", Pages: 500},
}

var maxID = 1

func generateNextID() int {
	maxID++
	return maxID
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, bookshelf)
}

func getBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "book not found"})
		return
	}

	book, _, err := findBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func addBook(c *gin.Context) {
	var newBook Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "book not found"})
		return
	}

	for _, existingBook := range bookshelf {
		if existingBook.Name == newBook.Name {
			c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
			return
		}
	}

	newBook.ID = generateNextID()
	bookshelf = append(bookshelf, &newBook)

	c.JSON(http.StatusCreated, newBook)
}

func deleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "book not found"})
		return
	}

	_, index, err := findBookByID(id)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "book not found"})
		return
	}

	bookshelf = append(bookshelf[:index], bookshelf[index+1:]...)

	c.Status(http.StatusNoContent)
}

func updateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "book not found"})
		return
	}

	book, index, err := findBookByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	var updatedBook Book

	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "book not found"})
		return
	}

	for _, existingBook := range bookshelf {
		if existingBook.Name == updatedBook.Name && existingBook.ID != id {
			c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
			return
		}
	}

	book.Name = updatedBook.Name
	book.Pages = updatedBook.Pages

	c.JSON(http.StatusOK, book)
}

func findBookByID(id int) (*Book, int, error) {
	for i, b := range bookshelf {
		if b.ID == id {
			return b, i, nil
		}
	}
	return nil, -1, fmt.Errorf("book not found")
}

func main() {
	r := gin.Default()

	r.GET("/bookshelf", getBooks)
	r.GET("/bookshelf/:id", getBook)
	r.POST("/bookshelf", addBook)
	r.DELETE("/bookshelf/:id", deleteBook)
	r.PUT("/bookshelf/:id", updateBook)

	err := r.Run(":8087")
	if err != nil {
		fmt.Println(err)
		return
	}
}
