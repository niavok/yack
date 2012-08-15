package yack

import (
	"database/sql"
	"fmt"
)

type Part struct {
	id     int64
	offset int64
	size   int64
	fileId   int64
}

func NewPart(file *File, offset int64, size int64) *Part {
	result, err := db.driver.Exec("INSERT INTO part ('file', 'offset', 'size') VALUES(?,?,?);", file.Id(), offset, size)
	
    if err != nil {
        fmt.Println(err)
        return nil
    }
    
    id, _ := result.LastInsertId()
    fmt.Println("LastInsertId part ", id)
	var part = model.Parts.GetById(id)
	return part;
}

func (this *Part) Size() int64 {
	return this.size
}

func (this *Part) Offset() int64 {
	return this.offset
}

func (this *Part) SetSize(size int64) {
    _, err := db.driver.Exec("UPDATE part SET size=? WHERE id=?", size, this.id)
	if err != nil {
		fmt.Println("Error in SetSize: ",err)
	}
	
	this.size = size
}

func (this *Part) Delete() {
    _, err := db.driver.Exec("DELETE FROM part WHERE id=?", this.id)
	if err != nil {
		fmt.Println("Error in Delete: ",err)
	}
}

func LoadPart(row *sql.Rows) *Part {
   fmt.Println("LoadPart()")
   var part Part
   err := row.Scan(&part.id, &part.fileId, &part.offset, &part.size)
   
   if err != nil {
		fmt.Println("LoadPart Scan Error: ",err)
		return nil
    }
   return &part
}

/////////////////
// Parts list
/////////////////

type parts struct {
}

func (this parts) GetById(id int64) *Part {
	fmt.Println("parts: Get id=", id)

	rows, err := db.driver.Query("select * from part where id=?", id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		return LoadPart(rows)
	}
	return nil
}
