package main

import (
	"fmt"
	"log"
	//"time"

	"database/sql"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/basicauth"
	_ "github.com/mattn/go-sqlite3"
)

type AppConfig struct {
	DB *sql.DB
}

func createReading(ctx iris.Context) {
	db := ctx.Values().Get("db").(*sql.DB)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into readings(batch, json) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ctx.Params().Get("id"), ctx.FormValue("reading"))
	log.Printf(ctx.Params().Get("id"))
	log.Printf(ctx.FormValue("reading"))

	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	ctx.Text("ok")
}

func getReading(ctx iris.Context) {
	db := ctx.Values().Get("db").(*sql.DB)

	stmt, err := db.Prepare("select json from readings where batch = ?")
	if err != nil {
		// log.Fatal(err)
		ctx.Text(err.Error())
		return
	}
	defer stmt.Close()
	var json string
	err = stmt.QueryRow(ctx.Params().Get("id")).Scan(&json)
	if err != nil {
		// log.Fatal(err)
		ctx.Text(err.Error())
		return
	}

	fmt.Println(json)

	ctx.ViewData("json", json)
	// Render template file: ./views/hello.html
	ctx.View("show.html")

	// ctx.JSON(iris.Map{
	// 	"message": name,
	// })
}

func createDb(fileName string) *sql.DB {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	sqlStmt := `
	create table readings (batch VARCHAR(32) not null primary key, json text);
	delete from readings;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		// return
	}
	return db
}

func newIrisApp(appConfig *AppConfig) *iris.Application {
	// middleware
	authentication := basicauth.Default(map[string]string{"user": "password"})
	database := func(ctx iris.Context) { ctx.Values().Set("db", appConfig.DB); ctx.Next() }

	app := iris.Default()
	app.RegisterView(iris.HTML("./views", ".html"))
	app.Post("/readings/{id:string}", authentication, database, createReading)
	app.Get("/readings/{id:string}", database, getReading)
	// tm := time.Unix(1518328047, 0)
	// fmt.Println(tm)
	return app
}

func main() {
	appConfig := AppConfig{DB: createDb("./foo.db")}
	irisApp := newIrisApp(&appConfig)
	// listen and serve on http://0.0.0.0:8080.
	irisApp.Run(iris.Addr(":8080"))
}
