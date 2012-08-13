package yack

import (
	"database/sql"
	"fmt"
)

type Part struct {
	id     int64
	offset int64
	size   int64

}

func (this *Part) Size() int64 {
	return this.size
}

func (this *Part) Offset() int64 {
	return this.offset
}

func LoadPart(row *sql.Rows) *Part {
   fmt.Println("LoadPart()")
   var part Part
   err := row.Scan(&part.id, &part.offset, &part.size)
   
   if err != nil {
		fmt.Println("LoadPart Scan Error: ",err)
		return nil
    }
   return &part
}

