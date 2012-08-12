package yack

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	UPLOADED = "uploaded"
	UPLOADING = "uploading"
)

type File struct {
	id           int
	name         string
	creationDate time.Time
	size         int
	uploadedSize int
	sha          string
	uploadState  string
	file         string
	owner        *User
	isPublic     bool
	description  string
	mime         string
	autoMime     bool
}

func NewFile(user *User, name string, sha string, size int) *File {
	_, err := db.driver.Exec("INSERT INTO file ('owner', 'name', 'creationDate', 'sha', 'size', 'mime', 'autoMime', 'description', 'uploadState', 'file', 'uploadedSize') VALUES(?,?,?,?,?,?,?,?,?,?,?);", user.Id(), name, time.Now(), sha, size, "", true, "", UPLOADING, "", 0)
	
	
	
    if err != nil {
        fmt.Println(err)
        return nil
    }
    
    //TODO Check already exist

	return model.Files.GetBySha(sha)
}


func LoadFile(row *sql.Rows) *File {
    fmt.Println("LoadFile()")
	var file File

    var creationDate string
    var ownerId int

    err := row.Scan(&file.id, &file.name, &creationDate, &file.size, &file.uploadedSize , &file.sha, &file.uploadState, &file.file, &ownerId, &file.description, &file.mime, &file.autoMime)
    if err != nil {
		fmt.Println("LoadFile Scan Error: ",err)
		return nil
	}

    fmt.Println("loaded id=",file.id, " name=", file.name, " size=", file.size, " uploadedSize=", file.uploadedSize, " sha=", file.sha, " uploadState=", file.uploadState, " file=", file.file, " ownerId=", ownerId, " description=", file.description, " mime=", file.mime, " autoMime=", file.autoMime)
    return &file
}

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

func (this *File) Progress() float64 {
	if this.uploadState == UPLOADED {
		return 1
	}

	return float64(this.uploadedSize) / float64(this.size)

}

func (this *File) Parts() []*Part {
	rows, err := db.driver.Query("SELECT * FROM part WHERE file=?", this.id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		LoadPart(rows)
	}
	rows.Close()

	return nil
}

/*
def get_progress(self):
        if self.upload_state == "uploaded":
            return 1;

        uploaded = 0;

        for part in self.parts.all():
            uploaded += part.size

        return float(uploaded)/self.size*/

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

/////////////////
// Files list
/////////////////

type files struct {
}

func (this files) GetBySha(sha string) *File {
	fmt.Println("files: GetBySha sha=", sha)

	rows, err := db.driver.Query("SELECT * FROM file WHERE sha=?", sha)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println("has result, load")
		return LoadFile(rows)
	}
	fmt.Println("no result")
	return nil

}


