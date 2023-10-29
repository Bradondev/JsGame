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


var CurrentArea ExploreAbleArea
var CurrentPlayer PlayerClass
var CurrentEnemy Enemy
var battleAnnouncement string 
var CurrentGame = GameState{
	InBattle: false,
	CanExplore: false,
}

type GameMessage struct{
	Announcement string `json:"Announcement"`
}


type GameState struct{
	InBattle bool
	CanExplore bool
	
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
type ExploreAbleArea struct{
	ID     string `json:"id"`
	AreasName string `json:"name"`
	AreasRecommenedLevel string `json:"recomLevel"`
	MonstersFoundInArea []Enemy `json:"monstersInArea"`
	AmountOfRoomsInArea int32 `json:"amountOfRooms"`
	AmountOfMonstersInArea int32 `json:"amountOfMonstersInArea"`
	AmountOfTreasureRooms int32 `json:"amountOfTreasureRooms"`
	AreaCleared bool 
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


///data used for the game 

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
	{ID: "4", Name: "Krugs",Attack: 2,Defence: 5,Health: 10,Defeated: false,Speed: 2,ExperincePoints: 10},
	{ID: "5", Name: "Bear", Attack: 4, Defence:  4,Health: 16, Defeated: false,Speed: 5, ExperincePoints: 16},
	{ID: "6", Name: "Skeleton", Attack: 5, Defence: 2, Health: 15, Defeated: false, Speed: 7, ExperincePoints: 17},
	{ID: "7", Name: "Warerat", Attack: 6, Defence: 4, Health: 20,Defeated: false ,Speed: 5,ExperincePoints: 20},
	{ID: "8", Name: "Undead Knight", Attack: 7, Defence: 5, Health:  23 ,Defeated:  false,Speed: 5,ExperincePoints: 25},
}


var ForestMonsters = []Enemy{Enemies[0],Enemies[1],Enemies[2],}
var CaveMonsters = []Enemy{Enemies[1],Enemies[3],Enemies[4]}
var DungeonMonsters = []Enemy{Enemies[7],Enemies[6],Enemies[5]}


var Areas = []ExploreAbleArea{
	{ID: "1", AreasName: "Forest",AreasRecommenedLevel: "1",MonstersFoundInArea: ForestMonsters, AmountOfRoomsInArea: 5,AmountOfMonstersInArea: 2,AmountOfTreasureRooms: 1,AreaCleared: false },
	{ID: "2", AreasName: "Cave",AreasRecommenedLevel: "3", MonstersFoundInArea:CaveMonsters,AmountOfRoomsInArea: 6,AmountOfMonstersInArea: 3,AmountOfTreasureRooms: 2 ,AreaCleared: false},
	{ID:"3",AreasName: "Dungeon",AreasRecommenedLevel: "5", MonstersFoundInArea:DungeonMonsters,AmountOfRoomsInArea: 7,AmountOfMonstersInArea: 4 ,AmountOfTreasureRooms: 3 ,AreaCleared: false },
}

var GameMessages = []GameMessage{{
	Announcement: ""},
}



//////////////

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
	c.IndentedJSON(http.StatusOK,gin.H{"message": "Please choose a clas by entering : SetClass/Mage or Knight or Tank or Rouge"  })
}
func ReturnCurrentPlayer(c *gin.Context){
	c.IndentedJSON(http.StatusOK,CurrentPlayer)
}


func GetPlayerByName(name string) (*PlayerClass,error){
	for i, b:= range PlayerCharacter {
		if b.Name == name {
			CurrentPlayer =  PlayerCharacter[i]
			return &PlayerCharacter[i],nil
		}
	}
	return nil, errors.New("no player found")
}

func GetAreaById(id string) (*ExploreAbleArea,error){
	for i, b:= range Areas{
		if b.ID == id {
			 CurrentArea=  Areas[i]
			return &Areas[i],nil
		}
	}
	return nil, errors.New("no player found")
}
func GetListOfAreas(c *gin.Context){
	c.IndentedJSON(http.StatusOK, Areas)

}

func SetClass(c *gin.Context){
	ClassName := c.Param("ClassName")
	PlayerClass, err := GetPlayerByName(ClassName)

	if err != nil {
		GameMessages[0].Announcement ="No Class found with that Name found. The Classes that can be picked are Mage, Knight, Tank ,and Rouge"
		c.IndentedJSON(http.StatusNotFound, GameMessages[0])
		return
	}
	GameMessages[0].Announcement = "Your Current Class is " + PlayerClass.Name +".Now you are able to explore please enter : Explore/Forest or Cave or Dungon"
	c.IndentedJSON(http.StatusOK, GameMessages[0]) 
	CurrentPlayer = *PlayerClass
	
}


// func ExploreArea(c *gin.Context){
// 	CurrentArea.

	
// }






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
	SetPlayer(PlayerCharacter[1])
	//BattleTurn(CurrentEnemy, CurrentPlayer)

	router.GET("Explore/:AreaName")
	router.GET("SetClass/:ClassName" ,SetClass )



	router.GET("GetAreas", GetListOfAreas)
	router.GET("GetPlayer",ReturnCurrentPlayer)
	//router.GET("SetPlayer/:id", GinChoosePlayerById)
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
