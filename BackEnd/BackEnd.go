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
var DialougeTemp int = 0
var counterForArea int = 0
var MonsterEncountered int 
var CurrentArea ExploreAbleArea
var CurrentPlayer PlayerClass
var CurrentEnemy Enemy
var battleAnnouncement string 
var CurrentGame = GameState{
	InBattle: false,
	CanExplore: false,
}
type ListOfDialouge struct{
	NameOfDialouge string
	Dialouge []string `json:"Dialouge"`
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
	OrderOfRooms []int32 // 0 = TreasureRoom ,1 = monster , 2 = emtpy
}


func (player PlayerClass) TakeDamage(DamageAmount int32){
	CurrentPlayer.Health-= DamageAmount
	if CurrentPlayer.Health < 0 {
		CurrentPlayer.Health =0
	}
}


func (enemy Enemy) TakeDamage(DamageAmount int32){
	
	CurrentEnemy.Health -= DamageAmount
	if CurrentEnemy.Health < 0 {
		CurrentEnemy.Health = 0
	}
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

var GameDialouge =[]ListOfDialouge{
	{NameOfDialouge: "Starting Dialouge", Dialouge: []string{"Welcome to RPG.JS in this text based game you will go on an adventure to defeat the Demon King , Please enter : Go, to Continue","In order to get strong enough to fight the Demon King ,you'll need to level up and get better equipment, Please enter : Go, to Continue", " You'll do this by exploring the 4 different Areas, looting these areas ,and defeating monsters, Please enter : Go, to Continue", " Your adventure will begin which you choosen a class by entering: SetClass/Mage or Knight or Tank or Rouge" }},
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
var DungeonLayout = []int32{1,2,0,1,2,1,2}
var ForestLayout =[]int32{2,1,0,2,1}
var CaveLayout =[]int32{1,0,2,1,2,1}

var Areas = []ExploreAbleArea{
	{ID: "1", AreasName: "Forest",AreasRecommenedLevel: "1",MonstersFoundInArea: ForestMonsters, AmountOfRoomsInArea: 5,AmountOfMonstersInArea: 2,AmountOfTreasureRooms: 1,AreaCleared: false ,OrderOfRooms: ForestLayout},
	{ID: "2", AreasName: "Cave",AreasRecommenedLevel: "3", MonstersFoundInArea:CaveMonsters,AmountOfRoomsInArea: 6,AmountOfMonstersInArea: 3,AmountOfTreasureRooms: 2 ,AreaCleared: false, OrderOfRooms:CaveLayout },
	{ID:"3",AreasName: "Dungeon",AreasRecommenedLevel: "5", MonstersFoundInArea:DungeonMonsters,AmountOfRoomsInArea: 7,AmountOfMonstersInArea: 3 ,AmountOfTreasureRooms: 3 ,AreaCleared: false, OrderOfRooms:DungeonLayout },
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


 
// func BattleTurn(enemy Enemy ,player PlayerClass){
// 	battleAnnouncement:= enemy.Name + " and " + player.Name +" have enter batte "
// 	println(battleAnnouncement)
// 	if enemy.Speed >= player.Speed{
// 		enemy.EnemyAttack(enemy,player,true)
		
// 	}else {
// 		player.SimpleAttack(enemy, player,true)
// 	}
 
// }

func SetHpToOfPlayer(NewHealth int32) {
	CurrentPlayer.Health = NewHealth
}

 
func SetHpTOfCurrentEnemy(NewHealth int32) {
	CurrentEnemy.Health = NewHealth
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
	ProgressOneTurnInBattle()
	c.IndentedJSON(http.StatusOK,GameMessages[0])
	}

func ProgressOneTurnInBattle(){
 if CurrentGame.InBattle{
	if CurrentPlayer.Speed > CurrentEnemy.Speed{
		CurrentEnemy.TakeDamage(CurrentPlayer.Attack) 
		GameMessages[0].Announcement = CurrentPlayer.Name + " attacked " + CurrentEnemy.Name + " You dealt: " + strconv.Itoa(int(CurrentPlayer.Attack)) +" Points of damage. "+ CurrentEnemy.Name  + "'s HP  decreased to: "+ strconv.Itoa(int(CurrentEnemy.Health))
		if CurrentEnemy.Health > 0{
			CurrentPlayer.TakeDamage(CurrentEnemy.Attack) 
			GameMessages[0].Announcement += CurrentEnemy.Name + "  attacked you and dealt: " + strconv.Itoa(int(CurrentEnemy.Attack)) + " Points of damage your HP decreased to: "+ strconv.Itoa(int(CurrentPlayer.Health)) + " "
			if CurrentPlayer.Health <= 0{
				GameMessages[0].Announcement += " You have been defeated by " + CurrentEnemy.Name +". You Can no longer progress. Please Restart to try again"
				CurrentGame.InBattle = false
			}
		}else{
			GameMessages[0].Announcement += " You have defeated " + CurrentEnemy.Name +". You may progress in this area now"
			CurrentGame.InBattle = false
		}
	} else {
		CurrentPlayer.TakeDamage(CurrentEnemy.Attack) 
		GameMessages[0].Announcement = CurrentEnemy.Name + "  attacked you and dealt: " + strconv.Itoa(int(CurrentEnemy.Attack)) + " Points of damage your HP decreased to: "+ strconv.Itoa(int(CurrentPlayer.Health)) + " "
		if CurrentPlayer.Health > 0{
		CurrentEnemy.TakeDamage(CurrentPlayer.Attack)
		GameMessages[0].Announcement += CurrentPlayer.Name + " attacked " + CurrentEnemy.Name  + " You dealt: " + strconv.Itoa(int(CurrentPlayer.Attack))+ "Points of damage. " + CurrentEnemy.Name  +"s HP + decreased to: "+ strconv.Itoa(int(CurrentEnemy.Health))
			if CurrentEnemy.Health <= 0{

				GameMessages[0].Announcement += " You have defeated " + CurrentEnemy.Name +". You may progress in the area now"
				CurrentGame.InBattle = false
			}
		}else{
			GameMessages[0].Announcement += " You have been defeated by " + CurrentEnemy.Name +". You Can no longer progress. Please Restart to try again"
			CurrentGame.InBattle = false
		}
		
	}
	
 }else{
	GameMessages[0].Announcement = "You must be in battle to be able to progress in battle "
 }
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
	

	c.IndentedJSON(http.StatusOK,gin.H{"message": "Please enter : Go"})
}
func ProgressDialouge( c *gin.Context){
	if DialougeTemp < len(GameDialouge[0].Dialouge){
		GameMessages[0].Announcement = GameDialouge[0].Dialouge[DialougeTemp]
		DialougeTemp++
		
	}else{
		GameMessages[0].Announcement = "This Dialouge has Enned"
	}
	c.IndentedJSON(http.StatusOK,GameMessages[0])
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

func GetAreaByName(Name string) (*ExploreAbleArea,error){
	for i, b:= range Areas{
		if b.AreasName == Name {
			 CurrentArea=  Areas[i]
			return &Areas[i],nil
		}
	}
	return nil, errors.New("no Area found")
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

func ExploreArea(c *gin.Context){
	if CurrentGame.CanExplore{
	AreaName := c.Param("AreaName")
	Areas , err := GetAreaByName(AreaName)
	if err != nil {
		GameMessages[0].Announcement ="No Area found with that Name found. The Areas that can be picked are Forest, Cave, Or Dungon"
		c.IndentedJSON(http.StatusNotFound, GameMessages[0])
		return
	}
	GameMessages[0].Announcement = "You have entered the " + Areas.AreasName +" To progress to the next room  please enter : Progress"
	CurrentGame.CanExplore = false
	CurrentArea = *Areas
	c.IndentedJSON(http.StatusOK, GameMessages[0]) 
}else{
	GameMessages[0].Announcement = " You are un able to change the current area at this time"
	c.IndentedJSON(http.StatusOK, GameMessages[0]) 
}
	
}




func progress(c *gin.Context) {
	if !CurrentGame.InBattle {
	
    if counterForArea >= len(CurrentArea.OrderOfRooms) {
        // Reset the counters and respond when you reach the end of the area
        MonsterEncountered = 0
        counterForArea = 0
        GameMessages[0].Announcement = "You have reached the end of the Area"
		CurrentPlayer.Health = 20
        c.IndentedJSON(http.StatusOK, GameMessages[0])
        return
    }

    roomType := CurrentArea.OrderOfRooms[counterForArea]
    counterForArea++

    switch roomType {
    case 1:
        // Monster room
        if MonsterEncountered < len(CurrentArea.MonstersFoundInArea) {
            GameMessages[0].Announcement = "You have Encountered a " + CurrentArea.MonstersFoundInArea[MonsterEncountered].Name + ". You must fight to be able to progress" + " Please enter: Battle , to start the fight"
			CurrentEnemy = CurrentArea.MonstersFoundInArea[MonsterEncountered]
			CurrentGame.InBattle = true
            MonsterEncountered++
        } else {
            // Handle the case where there are no more monsters to encounter
            GameMessages[0].Announcement = "You have Encountered an empty room. To progress to the next room, please enter: Progress"
        }
    case 0:
        // Treasure room
        GameMessages[0].Announcement = "You have Encountered a Treasure Room. You found an item. To progress to the next room, please enter: Progress"
    case 2:
        // Empty room
        GameMessages[0].Announcement = "You have Encountered an Empty Room. To progress to the next room, please enter: Progress"
    default:
        // Handle unexpected room types
        GameMessages[0].Announcement = "Unexpected room type encountered"
    }

    c.IndentedJSON(http.StatusOK, GameMessages[0])
}else{
	GameMessages[0].Announcement = "You can not progress unless you defeat the current Enemy"
	c.IndentedJSON(http.StatusOK, GameMessages[0])
}
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
	SetPlayer(PlayerCharacter[1])
	CurrentGame.CanExplore = true
	router.GET("Explore/:AreaName" ,ExploreArea)
	router.GET("SetClass/:ClassName" ,SetClass )
	router.GET("Progress",progress)
	router.GET("GetAreas", GetListOfAreas)
	router.GET("Go",ProgressDialouge)
	router.GET("Startgame", StartGame)
	router.GET("Enemy/:id", GinGetEnemybyId)
	router.GET("/En" ,GetListOfEnemy)
	router.GET("/Pr", GetListOfPlayer)
	router.GET("Battle",GinStartBattle)
	router.GET("/CurrentUnits/:UnitType",GinGetCurrentPlayerOrEnemy)
	router.Run("localhost:8080")


	//router.PATCH("/return",returnBook)
	//router.PATCH("/checkout", checkoutBook)
	//router.GET("/books", getBooks)
	//router.GET("/books/:id", bookById)
	//router.POST("/books", createBooks)
	//router.Run("localhost:8080")
	
}
