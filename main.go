package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// モデルの定義
type Event struct {
	Id        uint `gorm:"primaryKey"`
	Time      time.Time
	Location  string
	Event     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DBのインスタンスをグローバル変数に格納
var (
	db *gorm.DB
)

// この関数を追加
func port() string {

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	return ":" + port
}

// func connect() {
// 	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	databaseUrl := os.Getenv("DATABASE_URL")
// 	db, err := gorm.Open("postgres", databaseUrl)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// }

func main() {
	var err error
	err = godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		log.Fatal(err)
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	db, err = gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	connect()
	// 	// err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// 	//databaseUrl := os.Getenv("DATABASE_URL")
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	// Routing
	// フロントエンドがまだないので全部GET。あとでデータ取得以外はPOSTに直す。
	e.GET("/event/:id", getEventById)
	e.GET("/create", createEvent)
	e.GET("/update/:id", updateEventById)
	e.GET("/delete/:id", deleteEventById)

	// Port番号を関数から取得
	e.Logger.Fatal(e.Start(port()))
}

// eventsテーブルのレコードをid指定で取得
func getEventById(c echo.Context) error {
	var event Event

	id := c.Param("id")
	db.Find(&event, id)
	// 取得したデータをJSONにして返却
	return c.JSON(http.StatusOK, event)
}

// eventsテーブルにレコードを登録
func createEvent(c echo.Context) error {
	event := Event{Time: time.Now(), Location: "リビング", Event: "朝食"}
	db.Create(&event)
	// 取得したデータをJSONにして返却
	return c.String(http.StatusOK, "record has been created")
}

// eventsテーブルのレコードをid指定で更新
func updateEventById(c echo.Context) error {
	var event Event
	// var event := Event{Time: time.Now(), Location: "リビング", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	id := c.Param("id")
	db.Find(&event, id)

	event.Event = "夕食"

	db.Save(&event)

	// 取得したデータをJSONにして返却
	return c.JSON(http.StatusOK, "record has been updated")
}

// eventsテーブルのレコードをid指定で削除
func deleteEventById(c echo.Context) error {
	// var event Event

	id := c.Param("id")
	// db.Find(&event, id)

	db.Delete(&Event{}, id)

	// 取得したデータをJSONにして返却
	return c.JSON(http.StatusOK, "record has been deleted")
}
