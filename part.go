package yack

import (
	"database/sql"
	"fmt"
)

type Part struct {
	id     int
	offset int
	size   int

}

func (this *Part) Size() int {
	return this.size
}

func LoadPart(row *sql.Rows) *Part {
   fmt.Println("LoadPart()")
   var part Part
   row.Scan(&part.id, &part.offset, &part.size)
   return &part
}

