package yack

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type User struct {
	id                int64
	email             string
	creationDate      time.Time
	authToken         string
	authTokenValidity time.Time
	isAdmin           bool
	rootPackId        int64
	rootPack          *Pack
}

func NewUser(email string) *User {

	// Add root pack
	var rootPack = NewPack("")

	result, err := db.driver.Exec("insert into user ('email', 'rootPack', 'authToken', 'authTokenValidity') values(?, ?, ?, ?);", email, rootPack.Id(), "", time.Now())
	if err != nil {
		fmt.Println(err)
		return nil
	}

    id, _ := result.LastInsertId()
        fmt.Println("LastInsertId user", id)
        
	var user = model.Users.GetById(id)
	
    rootPack.SetOwner(user)
	return user;
}

func LoadUser(row *sql.Rows) *User {
	fmt.Println("LoadUser()")
	var user User
	var authTokenValidity string
	err := row.Scan(&user.id, &user.email, &user.authToken, &authTokenValidity, &user.rootPackId)
	if err != nil {
		fmt.Println("LoadUser Scan Error: ",err)
		return nil
	}
	

	fmt.Println("authTokenValidity=", authTokenValidity)
	user.authTokenValidity, _ = time.Parse("2006-01-02 15:04:05", authTokenValidity)

	fmt.Println("loaded id=", user.id, " email=", user.email, " authToken=", user.authToken, " authTokenValidity=", user.authTokenValidity)
	return &user
}

func (this *User) Id() int64 {
	return this.id
}

func (this *User) DisplayName() string {
	fmt.Println("DisplayName() ", this.email)
	return this.email
}

func (this *User) RootPack() *Pack {
    if this.rootPack == nil {
        this.rootPack = model.Packs.GetById(this.rootPackId)
    }
	return this.rootPack
}

func (this *User) AuthToken() string {
	fmt.Println("AuthToken()")
	if this.authToken == "" || this.authTokenValidity.Before(time.Now()) {
		if this.authToken == "" {
			fmt.Println("generateAuthToken because empty auth token")
		}
		if this.authTokenValidity.Before(time.Now()) {
			fmt.Println("generateAuthToken because ", this.authTokenValidity, " is before ", time.Now())
		}

		this.generateAuthToken()
	}
	fmt.Println("this=", &this)
	fmt.Println("authToken=", this.authToken)
	return this.authToken
}

func (this *User) generateAuthToken() {
	fmt.Println("generateAuthToken to id ", this.id)
	alphabet := "abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	password := make([]byte, 32)
	for i := 0; i < len(password); i++ {
		password[i] = alphabet[rand.Int()%len(alphabet)]
	}

	this.authToken = string(password)
	this.authTokenValidity = time.Now().Add(time.Hour * 24 * 15) // 15 days
	_, err := db.driver.Exec("UPDATE user SET authToken=?, authTokenValidity=? WHERE id=?;", this.authToken, this.authTokenValidity, this.id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("this=", &this)
	fmt.Println("token generated ", this.authToken, " valid until ", this.authTokenValidity)

}

func (this *User) GetInterruptedFiles() []*File {
	rows, err := db.driver.Query("SELECT * FROM file WHERE owner=? AND uploadState!=? ", this.id, UPLOADED)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	
	files := []*File{}
	
	for rows.Next() {
        files = append(files, LoadFile(rows))
	}

	return files
}

/////////////////
// Users list
/////////////////

type users struct {
}

func (this users) GetByAuthToken(authToken string, id int64) *User {
	fmt.Println("users: Get authToken=", authToken, " id=", id)

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

func (this users) GetById(id int64) *User {
	fmt.Println("users: Get id=", id)

	rows, err := db.driver.Query("select * from user where id=?", id)
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

func (this users) Equal(user1 *User, user2 *User) bool {
    if user1 == nil || user2 == nil {
        return user1 == user2
    }     
    return user1.Id() == user2.Id()
}


func (this users) GetByEmail(email string) *User {
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


