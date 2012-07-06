package yack

import (
	"database/sql"
	"fmt"
)

type PartList struct {
	id     int
	offset int
	size   int
}

func LoadPartList(row *sql.Rows) *PartList {
	fmt.Println("PartList()")
	var partList PartList
	row.Scan(&partList.id, &partList.offset, &partList.size)
	return &partList
}

func (this *PartList) Size() int {
	return this.size
}
