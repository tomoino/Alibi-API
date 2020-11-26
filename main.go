package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// この関数を追加
func port() string {

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	return ":" + port
}

func connect() {
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
		databaseUrl := ""
		if err != nil {
			// .env読めなかった場合の処理
			databaseUrl = "ENVないよ！"
		} else {
			databaseUrl = os.Getenv("DATABASE_URL")
		}
		return c.String(http.StatusOK, "Hello, World!"+databaseUrl)
	})

	// Port番号を関数から取得
	e.Logger.Fatal(e.Start(port()))
}
