package main

import (
	"errors"
	"net/http"

	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


type book struct {
	ID     string `json:"id"`
	Name  string `json:"name"`
	Attack int32 `json:"Attack"`
	Health int32 `json:"health"`
	Defence int32 `json:"defence"`
	
}

var books = []book{
	{ID: "1", Name: "In Search of Lost Time", Attack: 1, Defence: 3,Health: 20},
	{ID: "2", Name: "The Great gatsby", Attack: 2, Defence: 3,Health: 20},
	{ID: "3", Name: "War and Peace", Attack: 2, Defence: 2,Health: 20},
}


func getBooks(c *gin.Context){
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context){
	id := c.Param("id")
	book,err := getBooksById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return 
	}
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}
	book,err:= getBooksById(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)


}



func checkoutBook(c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}
	book,err:= getBooksById(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no books left"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}






func getBooksById(id string) (*book,error){
	for i, b:= range books {
		if b.ID == id {
			return &books[i],nil
		}
	}
	return nil, errors.New("books not found")
}

func createBooks(c *gin.Context){
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return 
	}
	books = append(books,newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main(){
	//router := gin.Default()
	//router.Use(cors.New(cors.Config{
		
	//}))
	//router.PATCH("/return",returnBook)
	//router.PATCH("/checkout", checkoutBook)
	//router.GET("/books", getBooks)
	//router.GET("/books/:id", bookById)
	//router.POST("/books", createBooks)
	//router.Run("localhost:8080")
	
}