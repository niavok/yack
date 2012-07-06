package yack

import (
	_ "database/sql"
	_ "fmt"
	_ "math/rand"
	"time"
)

type UserGroup struct {
	id           int
	name         string
	creationDate time.Time
	owner        *User
}

func (this *UserGroup) Id() int {
	return this.id
}

func (this *UserGroup) Name() string {
	return this.name
}
