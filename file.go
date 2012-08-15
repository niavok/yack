package yack

import (
	"database/sql"
	"fmt"
	"time"
	"io"
	"os"
	"crypto/sha1"
)

const (
	UPLOADED = "uploaded"
	UPLOADING = "uploading"
)

type File struct {
	id           int64
	name         string
	creationDate time.Time
	size         int64
	uploadedSize int64
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

func (this *File) Size() int64 {
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

func (this *File) SetFile(filePath string) {
    _, err := db.driver.Exec("UPDATE file SET file=? WHERE id=?", filePath, this.Id())
	if err != nil {
		fmt.Println("Error in SetFile: ",err)
	}
	
	this.file = filePath
}

func (this *File) SetUploadState(state string) {
    _, err := db.driver.Exec("UPDATE file SET uploadState=? WHERE id=?", state, this.Id())
	if err != nil {
		fmt.Println("Error in SetUploadState: ",err)
	}
	
	this.uploadState = state
}

func (this *File) SetUploadedSize(size int64) {
    _, err := db.driver.Exec("UPDATE file SET uploadedSize=? WHERE id=?", size, this.Id())
	if err != nil {
		fmt.Println("Error in SetUploadSize: ",err)
	}
	
	this.uploadedSize = size
}



func (this *File) NewPart(offset int64) *Part{
    return NewPart(this, offset, 0)
}


func (this *File) Parts() []*Part {
	rows, err := db.driver.Query("SELECT * FROM part WHERE file=? ORDER BY offset", this.id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	
	parts := []*Part{}
	
	for rows.Next() {
        parts = append(parts, LoadPart(rows))
	}

	return parts
}

func (this *File) AddData(offset int64 , size int64, sha string, data io.ReadCloser) {
    // Create physical file if needed
    var file *os.File
    var err error
    
    if this.file == "" {
        //TODO Use write file path instead of /home/fred/.local/share/yack/files/
        var fileName string = "/home/fred/.local/share/yack/files/"+this.Sha()
        file, err = os.Create(fileName)
        err = file.Truncate(this.Size())
        this.SetFile(fileName)
    } else {
        file, err = os.OpenFile(this.file, os.O_RDWR, 0666)
    }
    
    if err != nil {
        fmt.Println(err)
		return
	}
    
    defer file.Close()
    
    // Check existing parts
    var parts = this.Parts()
    var matchPart *Part = nil
    var writeOffset int64 = 0
    var writeSize int64 = 0
    
    // Find part matching to begining
    for _, part := range(parts) {
        if offset >= part.Offset() && offset <= part.Offset() + part.Size() {
            matchPart = part
            writeOffset = part.Offset() + part.Size()
            var overwrite int64 = writeOffset - offset
            writeSize = size - overwrite
            break
        }
    }
    
    // If no match, create an empty part at offset
    if matchPart == nil {
        matchPart = this.NewPart(offset)
        writeOffset = offset
        writeSize = size
        parts = this.Parts()
    }
    
    
    // Find overwrite part at end and reduce write size
    var writeEndOffset = writeOffset + writeSize
    for _, part := range(parts) {
        if writeEndOffset >= part.Offset() && writeEndOffset < part.Offset() + part.Size() {
            var overwrite = writeEndOffset - part.Offset()
            writeSize = writeSize - overwrite
        }
    }
    
    // Seek overwrite bytes
    var inputOffset = writeOffset - offset
    
    var buffer[1000000]byte
    var bufferSize int64 = int64(len(buffer))
    
    var skipBytes int64 = 0
    for skipBytes < inputOffset {
        var byteToRead = inputOffset-skipBytes
        if byteToRead > bufferSize {
            byteToRead = bufferSize
        }
        n, err := data.Read(buffer[0:byteToRead])
        if err != nil {
            fmt.Println("Error reading the data to skip ",err)
            return
        }
        skipBytes += int64(n)
    }
    
    //Write data
    file.Seek(writeOffset, 0) // 0 mean absolute seek
    
    var writeBytes int64 = 0
    for writeBytes < writeSize {
        var byteToRead = writeSize-writeBytes
        if byteToRead > bufferSize {
            byteToRead = bufferSize
        }
        nRead, errRead := data.Read(buffer[0:byteToRead])
        if errRead != nil {
            fmt.Println("Error reading the data to write ",errRead)
            return
        }
        
        nWrite, errWrite := file.Write(buffer[0:nRead])
        if errWrite != nil {
            fmt.Println("Error writing the file ", errWrite)
            return
        }
        writeBytes += int64(nWrite)
    }

    // Check sha
    file.Seek(offset, 0)
    h := sha1.New()
    var checkBytes int64 = 0
    for checkBytes < size {
        var byteToRead = size-checkBytes
        if byteToRead > bufferSize {
            byteToRead = bufferSize
        }
        nRead, errRead := file.Read(buffer[0:byteToRead])
        if errRead != nil {
            fmt.Println("Error reading the data to check ",errRead)
            return
        }
        
        nWrite, errWrite := h.Write(buffer[0:nRead])
        if errWrite != nil {
            fmt.Println("Error write the data to check ",errWrite)
            return
        }
        checkBytes += int64(nWrite)
    }
    
    var checksha = fmt.Sprintf("%x", h.Sum(nil))
    if checksha != sha {
       fmt.Printf("Wrong sha1 : %s expected but %s computed.", sha, checksha)
       return
    }

    // Check ok : Update parts
    matchPart.SetSize(matchPart.Size() + writeSize)
    this.SetUploadedSize(this.uploadedSize + writeSize)
    
    // Merge parts
    var partCount = len(parts)
    var mergePartIndex = 1
    var firstPart = parts[0]
    
    for mergePartIndex < len(parts) {
        var mergePart *Part = parts[mergePartIndex]
        if firstPart.Offset() + firstPart.Size() == mergePart.Offset() {
            firstPart.SetSize(firstPart.Size() + mergePart.Size())
            mergePart.Delete()
            partCount--
            mergePartIndex ++
        } else {
            break
        }
    }
    
    // Update state
    if firstPart.Size() == this.Size() {
        this.SetUploadState(UPLOADED)
    }
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



