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
	"github.com/labstack/echo/v4/middleware"
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

// request jsonの定義
type Request struct {
	Time     time.Time `json:"time"`
	Location string    `json:"location"`
	Event    string    `json:"event"`
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
	e.Use(middleware.CORS())

	// Routing
	// フロントエンドがまだないので全部GET。あとでデータ取得以外はPOSTに直す。
	e.GET("/event/:id", getEventById)
	e.GET("/events", getAllEvents)
	e.POST("/create", createEvent)
	e.POST("/update/:id", updateEventById)
	e.DELETE("/delete/:id", deleteEventById)

	// Port番号を関数から取得
	e.Logger.Fatal(e.Start(port()))
}

// GET : eventsテーブルのレコードをid指定で取得
func getEventById(c echo.Context) error {
	var event Event

	id := c.Param("id")
	db.Find(&event, id)
	// 取得したデータをJSONにして返却
	return c.JSON(http.StatusOK, event)
}

// GET : eventsテーブルのレコードを全件取得
func getAllEvents(c echo.Context) error {
	var event []Event

	db.Find(&event)
	// 取得したデータをJSONにして返却
	return c.JSON(http.StatusOK, event)
}

// POST : eventsテーブルにレコードを登録
func createEvent(c echo.Context) error {
	// リクエストをbind
	post := new(Request)
	c.Bind(post)

	var event Event

	// request json に応じて値を更新
	if post.Time.IsZero() != true {
		event.Time = post.Time
	} else {
		event.Time = time.Now() // 時間は入っていない場合は現在時刻を入れる
	}

	if len(post.Event) > 0 {
		event.Event = post.Event
	}

	if len(post.Location) > 0 {
		event.Location = post.Location
	}

	db.Create(&event)
	// 取得したデータをJSONにして返却
	return c.String(http.StatusOK, "record has been created")
}

// POST : eventsテーブルのレコードをid指定で更新
func updateEventById(c echo.Context) error {
	// リクエストをbind
	post := new(Request)
	c.Bind(post)

	var event Event
	// idで更新するレコードを取得
	id := c.Param("id")
	db.Find(&event, id)

	// request json に応じて値を更新
	// if post.Time.IsZero() != true {
	event.Time = post.Time
	// }

	if len(post.Event) > 0 {
		event.Event = post.Event
	}

	if len(post.Location) > 0 {
		event.Location = post.Location
	}

	event.UpdatedAt = time.Now()

	db.Save(&event)

	// 取得したデータをJSONにして返却
	return c.JSON(http.StatusOK, "record has been updated")
}

// DELETE : eventsテーブルのレコードをid指定で削除
func deleteEventById(c echo.Context) error {
	id := c.Param("id")

	db.Delete(&Event{}, id)

	// 取得したデータをJSONにして返却
	return c.JSON(http.StatusOK, "record has been deleted")
}
