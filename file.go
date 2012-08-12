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
	id           int64
	name         string
	creationDate time.Time
	size         int
	uploadedSize int
	sha          string
	uploadState  string
	file         string
    ownerId      int64
	isPublic     bool
	description  string
	mime         string
	autoMime     bool

    //Cache	
	owner        *User
	
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

    err := row.Scan(&file.id, &file.name, &creationDate, &file.size, &file.uploadedSize , &file.sha, &file.uploadState, &file.file, &file.ownerId, &file.description, &file.mime, &file.autoMime)
    if err != nil {
		fmt.Println("LoadFile Scan Error: ",err)
		return nil
	}

    fmt.Println("loaded id=",file.id, " name=", file.name, " size=", file.size, " uploadedSize=", file.uploadedSize, " sha=", file.sha, " uploadState=", file.uploadState, " file=", file.file, " ownerId=", file.ownerId, " description=", file.description, " mime=", file.mime, " autoMime=", file.autoMime)
    return &file
}

func (this *File) Id() int64 {
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

func (this *File) Owner() *User {
    if this.owner == nil {
        this.owner = model.Users.GetById(this.ownerId)
    }
	return this.owner
}

func (this *File) CanWrite(user *User) bool {
	if model.Users.Equal(this.Owner(),user) {
		return true
	}
	return false
}

func (this *File) CanRead(user *User) bool {
    if model.Users.Equal(this.Owner(),user) {
		return true
	}

    //TODO share
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

func (this files) GetById(id int64) *File {
	fmt.Println("files: Get id=", id)

	rows, err := db.driver.Query("select * from file where id=?", id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		return LoadFile(rows)
	}
	return nil
}



