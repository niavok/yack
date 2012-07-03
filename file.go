package yack

import (
	// "database/sql"
	"time"
)

type File struct {
	id           int
	name         string
	creationDate time.Time
	size         int
	sha          string
	uploadState  string
	file         string
	owner        *User
	isPublic     bool
	description  string
	mime         string
	autoMime     bool
}

/*
func NewFile(name string) *User {
	_, err := db.driver.Exec("insert into user ('email') values(?);", email)
    if err != nil {
        fmt.Println(err)
        return nil
    }

	return model.Users.GetByEmail(email)
}


func LoadUser(row *sql.Rows) *User {
    fmt.Println("LoadUser()")
	var user User
	var authTokenValidity string
    row.Scan(&user.id, &user.email, &user.authToken, &authTokenValidity)

    fmt.Println("authTokenValidity=", authTokenValidity)
    user.authTokenValidity, _ = time.Parse("2006-01-02 15:04:05", authTokenValidity)

    fmt.Println("loaded id=",user.id, " email=", user.email, " authToken=", user.authToken, " authTokenValidity=", user.authTokenValidity)
    return &user
}
*/

func (this *File) Id() int {
	return this.id
}

func (this *File) Name() string {
	return this.name
}

func (this *File) Size() int {
	return this.size
}

func (this *File) Sha() string {
	return this.sha
}

func (this *File) CanWrite(user *User) bool {
	if this.owner == user {
		return true
	}
	return false
}

/*

func (this *User) AuthToken() string{
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

}*/
