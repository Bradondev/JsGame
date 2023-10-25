package main

import (
	"fmt"
	//"errors"
	"net/http"

	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


var CurrentPlayer PlayerClass
var CurrentEnemy Enemy
var battleAnnouncement string 


type Enemy struct {
	ID     string `json:"id"`

	Name  string `json:"name"`
	Attack int32 `json:"Attack"`
	Health int32 `json:"health"`
	Defence int32 `json:"defence"`
	Defeated bool `json:"defeated"`
	Speed int32 `json:"speed"`
	ExperincePoints int32 `json:"points"`
	
}

type PlayerClass struct {

	ID     string `json:"id"`
	Level int32 `json:"level"`
	CurrentExperincePoints int32 `json:"points"`
	ExperincePointsUntilNextLevel int32 `json:"neededpoints"`
	Name  string `json:"name"`
	Attack int32 `json:"Attack"`
	Health int32 `json:"health"`
	Defence int32 `json:"defence"`
	Speed int32 `json:"speed"`
}

func (player PlayerClass) GainExperincePoints(enemy Enemy) {
	player.ExperincePointsUntilNextLevel += player.Level * 5
	player.CurrentExperincePoints += enemy.ExperincePoints
	if player.ExperincePointsUntilNextLevel <= player.CurrentExperincePoints{
		player.CurrentExperincePoints = player.ExperincePointsUntilNextLevel - player.CurrentExperincePoints
		player.Level +=1
		player.ExperincePointsUntilNextLevel = 0 
		fmt.Println("Player Leveled up they are now level " ,player.Level)
	}
}




var PlayerCharacter = []PlayerClass{
	{ID: "1", Name: "Mage", Attack: 6, Defence: 3,Health: 14,Speed: 5 ,Level: 0},
	{ID: "2", Name: "Knight", Attack: 3, Defence: 5,Health: 20, Speed: 6},
	{ID: "3", Name: "Rouge", Attack: 3, Defence: 4,Health: 16, Speed: 10},
	{ID: "4", Name: "Tank", Attack: 3, Defence: 4,Health: 22,Speed: 2 },
}



var Enemies = []Enemy{
	{ID: "1", Name: "Slime", Attack: 1, Defence: 1,Health: 8, Defeated: false ,Speed: 3, ExperincePoints: 5},
	{ID: "2", Name: "Goblin", Attack: 3, Defence: 5,Health: 20, Defeated: false ,Speed: 6, ExperincePoints: 10},
	{ID: "3", Name: "Bat", Attack: 3, Defence: 2,Health: 10, Defeated: false ,Speed: 8,ExperincePoints: 7},
}



//func bookById(c *gin.Context){
	//id := c.Param("id")
	//book,err := getBooksById(id)

	//if err != nil {
	//	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		//return 
	//}
	//c.IndentedJSON(http.StatusOK, book)
//}


 
func BattleTurn(enemy Enemy ,player PlayerClass){
	battleAnnouncement:= enemy.Name + " and " + player.Name +" have enter batte "
	println(battleAnnouncement)
	if enemy.Speed >= player.Speed{
		println("enemy is faster")
	}else {
		println("Player is faster")
	}
	
}





//func getBooksById(id string) (*book,error){
	//for i, b:= range books {
		//if b.ID == id {
			//return &books[i],nil
		//}
	//}
	///return nil, errors.New("books not found")
//}

func createBooks(c *gin.Context){
	var NewEnemy Enemy

	if err := c.BindJSON(&NewEnemy); err != nil {
		return 
	}
	Enemies = append(Enemies,NewEnemy)
	c.IndentedJSON(http.StatusCreated, NewEnemy)
}

func main(){
	BattleTurn(Enemies[0], PlayerCharacter[0])
	PlayerCharacter[0].GainExperincePoints(Enemies[0])
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