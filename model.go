package yack

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"time"
)

type Database struct {
	driver *sql.DB
}

type Model struct {
	Users users
	Packs packs
	Files files
}

var model Model
var db Database

func Init() {

	rand.Seed(time.Now().UnixNano())

	fmt.Println("Init database")
	var err error
	db.driver, err = sql.Open("sqlite3", "/home/fred/.local/share/yack/yack.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	//TODO close later defer db.driver.Close()

	db.createDatabase()
}

func (this Database) createDatabase() {
	fmt.Println("Create database")
	this.execQuery("CREATE TABLE user (id INTEGER PRIMARY KEY AUTOINCREMENT, email TEXT, authToken TEXT, authTokenValidity DATETIME, rootPack INTEGER);")

	this.execQuery("CREATE TABLE file (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, creationDate DATETIME, size INTEGER, uploadedSize INTEGER, sha TEXT, uploadState TEXT, file TEXT, owner INTEGER, description TEXT, mime TEXT, autoMime BOOL);")

	this.execQuery("CREATE TABLE pack (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, creationDate DATETIME, owner INTEGER, parentPack INTEGER, isPublic BOOL);")
	this.execQuery("CREATE TABLE pack_file (pack INTEGER , file INTEGER);")

	this.execQuery("CREATE TABLE part (id INTEGER PRIMARY KEY AUTOINCREMENT, file INTEGER,  offset INTEGER, size INTEGER);")

}

func GetModel() Model {
	return model
}

func (this Database) execQuery(query string) {
	_, err := this.driver.Exec(query)

	if err != nil {
		fmt.Println(err)
		return
	}
}
