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

type userList struct {
}

func (this userList) GetByAuthToken(authToken string, id int) *User {
	fmt.Println("userList: Get authToken=", authToken, " id=", id)

	db.driver.Query("select id, email from user")

	rows, err := db.driver.Query("select * from user where id=? AND authToken=?", id, authToken)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		return LoadUser(rows)
	}
	return nil

}

func (this userList) GetByEmail(email string) *User {
	fmt.Println("userList: Get email=", email)

	db.driver.Query("select id, email from user")

	rows, err := db.driver.Query("select * from user where email=?", email)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		return LoadUser(rows)
	}
	return nil
}

type Model struct {
	Users userList
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
	//defer db.driver.Close()

	db.createDatabase()

	rows, err := db.driver.Query("select id, email from user")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		println(id, name)
	}
	rows.Close()

}

func (this Database) createDatabase() {
	fmt.Println("Create database")
	this.execQuery("CREATE TABLE user (id INTEGER PRIMARY KEY AUTOINCREMENT, email TEXT, authToken TEXT, authTokenValidity DATETIME);")
	this.execQuery("CREATE TABLE file (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, creationDate DATETIME, size INTEGER, sha TEXT, uploadState TEXT, file TEXT, owner INTEGER, description TEXT, mime TEXT, autoMime BOOL);")

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
