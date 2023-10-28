package main

import (
	"fmt"
	//"go/printer"

	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//"google.golang.org/protobuf/internal/strs"
)


var CurrentPlayer PlayerClass
var CurrentEnemy Enemy
var battleAnnouncement string 



type message struct{
	Message string `json:"message"`
}

type Enemy struct {
	Name  string `json:"name"`
	ID     string `json:"id"`
	Attack int32 `json:"Attack"`
	Health int32 `json:"health"`
	Defence int32 `json:"defence"`
	Defeated bool `json:"defeated"`
	Speed int32 `json:"speed"`
	ExperincePoints int32 `json:"points"`

}

type PlayerClass struct {
	Name  string `json:"name"`
	ID     string `json:"id"`
	Level int32 `json:"level"`
	CurrentExperincePoints int32 `json:"points"`
	ExperincePointsUntilNextLevel int32 `json:"neededpoints"`
	Attack int32 `json:"Attack"`
	Health int32 `json:"health"`
	Defence int32 `json:"defence"`
	Speed int32 `json:"speed"`
}

func (e Enemy) EnemyAttack(enemy Enemy, player PlayerClass,FirstToAttack bool){
	SetBattleAnnouncement( player.Name+" : " +strconv.Itoa(int(player.Health))+"HP")
	if (player.Defence < enemy.Attack){
		player.Health -= enemy.Attack 
		SetBattleAnnouncement( player.Name +" Was attacked by " + enemy.Name)
		SetBattleAnnouncement( player.Name+" : " +strconv.Itoa(int(player.Health))+"HP")
	}
	
	if player.Health <= 0 {
		player.Health = 0
		EndCurrentBattle(enemy, player)
		return
	}
	if (FirstToAttack){
		player.SimpleAttack(enemy,player,false)	
	}
	
}

func (e Enemy) EnemyOnDeath( enemy Enemy ,player PlayerClass){
	fmt.Println("The enemy give the player exp")
	player.GainExperincePoints(enemy)
	EndCurrentBattle(enemy, player)
}

func (p PlayerClass) SimpleAttack( enemy Enemy, player PlayerClass,FirstToAttack bool) bool{	
	SetBattleAnnouncement(enemy.Name+" : " +strconv.Itoa(int(enemy.Health))+"HP")
	if (enemy.Defence < player.Attack){
		enemy.Health -= player.Attack
		SetBattleAnnouncement( enemy.Name +" Was attacked by " + player.Name)
		SetBattleAnnouncement(enemy.Name+" : " +strconv.Itoa(int(enemy.Health)) +"HP")
	}
	
	
	if enemy.Health <= 0 {
		enemy.Defeated = true
		enemy.EnemyOnDeath(enemy, player)
		
		return true
	}
	if (FirstToAttack){
		//Then the Enemy will attack play since they havent attacked yet
		enemy.EnemyAttack(enemy,player,false)	
	}
	return false
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
	{ID: "1", Name: "Mage", Attack: 6, Defence: 2,Health: 14,Speed: 5 ,Level: 1},
	{ID: "2", Name: "Knight", Attack: 3, Defence: 3,Health: 20, Speed: 6},
	{ID: "3", Name: "Rouge", Attack: 3, Defence: 2,Health: 16, Speed: 10},
	{ID: "4", Name: "Tank", Attack: 3, Defence: 4,Health: 22,Speed: 2 },
}



var Enemies = []Enemy{
	{ID: "1", Name: "Slime", Attack: 1, Defence: 1,Health: 6, Defeated: false ,Speed: 3, ExperincePoints: 5},
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
		enemy.EnemyAttack(enemy,player,true)
		
	}else {
		player.SimpleAttack(enemy, player,true)
	}
	CurrentEnemy.Health = enemy.Health
	CurrentPlayer.Health = player.Health
}


func EndCurrentBattle(enemy Enemy, player PlayerClass){
	if (enemy.Defeated == true){
		SetBattleAnnouncement("The current battle as End," +enemy.Name +" was defeated.")
	}else {
		SetBattleAnnouncement("The current battle as End,"+" the player has been defeated.")
	}
}


func SetBattleAnnouncement(NewAnnouncement string){
	battleAnnouncement = NewAnnouncement
	println(battleAnnouncement)
}



func GinGetEnemybyId(c *gin.Context) {
	id := c.Param("id")
	Enemy, err := GetEnemybyId(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, Enemy)
}



func GetEnemybyId(id string) (*Enemy,error){
	for i, b:= range Enemies {
		if b.ID == id {
			return &Enemies[i],nil
		}
	}
	return nil, errors.New("books not found")
}

func createBooks(c *gin.Context){
	var NewEnemy Enemy

	if err := c.BindJSON(&NewEnemy); err != nil {
		return 
	}
	Enemies = append(Enemies,NewEnemy)
	c.IndentedJSON(http.StatusCreated, NewEnemy)
}


func SetPlayer(player PlayerClass) {
	CurrentPlayer = player 
	fmt.Println("Current player is ", player.Name)
}

func SetEnemy(enemy Enemy) {
	CurrentEnemy = enemy
	fmt.Println("Current enemy is ", enemy.Name)
}



func GetListOfEnemy(c *gin.Context){
	c.IndentedJSON(http.StatusOK, Enemies)
}

func GetListOfPlayer(c *gin.Context){
	c.IndentedJSON(http.StatusOK, PlayerCharacter)
}

func CreateNewEnemy(c *gin.Context){
	var NewEnemy Enemy

	if err:= c.BindJSON(&NewEnemy); err != nil{
		return
	}
	
	Enemies = append(Enemies, NewEnemy)
	c.IndentedJSON(http.StatusCreated,NewEnemy)
}

func GinStartBattle(c *gin.Context){
	var player PlayerClass = CurrentPlayer
	var enemy Enemy = CurrentEnemy

	battleAnnouncement:= enemy.Name + " and " + player.Name +" have enter batte "
	SetBattleAnnouncement(battleAnnouncement)
	println("test")

	
	


	if enemy.Speed >= player.Speed {
		//enemy.EnemyAttack(enemy,player,true
		CurrentPlayer.Health -=CurrentEnemy.Attack
		
		if CurrentPlayer.Health >= 0 {
			CurrentEnemy.Health-= CurrentPlayer.Attack
			
		}else {
			println("Player is dead")
		}

	}else {
		//player.SimpleAttack(enemy, player,true)	
		CurrentEnemy.Health-= CurrentPlayer.Attack
		if CurrentEnemy.Health >= 0 {
			CurrentPlayer.Health-= CurrentEnemy.Attack
		
		}else{
			println("Enemy is dead")
		}
		c.IndentedJSON(http.StatusAccepted, CurrentPlayer)
	}
	






	c.IndentedJSON(http.StatusOK,CurrentPlayer)
}

func GinGetCurrentPlayerOrEnemy(c *gin.Context){
	EnemyOrPlayer := c.Param("UnitType")
	println(EnemyOrPlayer)
	if (EnemyOrPlayer == "Player"){
		c.IndentedJSON(http.StatusOK,CurrentPlayer)
	}else{
		c.IndentedJSON(http.StatusOK,CurrentEnemy)
	}
	
}
func StartGame(c *gin.Context){
	c.IndentedJSON(http.StatusOK,gin.H{"message": "The Game Will Start once the button is clicked"})
}


func main(){
	router := gin.Default()

	//SetPlayer(PlayerCharacter[0])

	SetEnemy(Enemies[1])
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"PUT","PACTH","POST","DELETE","GET"},
		AllowHeaders: []string{"Content-Type"},
		AllowCredentials: true,
		
	}))
	//BattleTurn(CurrentEnemy, CurrentPlayer)
	router.GET("Startgame", StartGame)
	router.GET("Enemy/:id", GinGetEnemybyId)
	router.GET("/En" ,GetListOfEnemy)
	router.GET("/Pr", GetListOfPlayer)
	router.PATCH("/StartBattle",GinStartBattle)
	router.GET("/CurrentUnits/:UnitType",GinGetCurrentPlayerOrEnemy)
	router.Run("localhost:8080")


	//router.PATCH("/return",returnBook)
	//router.PATCH("/checkout", checkoutBook)
	//router.GET("/books", getBooks)
	//router.GET("/books/:id", bookById)
	//router.POST("/books", createBooks)
	//router.Run("localhost:8080")
	
}
