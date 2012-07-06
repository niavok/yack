package yack

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Pack struct {
	id           int
	name         string
	creationDate time.Time
	owner        *User
	parentPack   *Pack
	isPublic     bool
}

func LoadPack(row *sql.Rows) *Pack {
	fmt.Println("LoadPack()")
	var pack Pack
	row.Scan(&pack.id, &pack.name, &pack.creationDate, &pack.isPublic)
	return &pack
}

func (this *Pack) Id() int {
	return this.id
}

func (this *Pack) CanRead(user *User) bool {
	if this.owner == user {
		return true
	}

	if this.IsSharedToUser(user) {
		return true
	}

	return this.isPublic
}

func (this Pack) IsSharedToUser(user *User) bool {
	fmt.Println("packs: IsSharedToUser user=", user.DisplayName())

	// Is shared to this user
	rows_user, err := db.driver.Query("select * from pack_share_user where user=? AND pack=?", user.Id(), this.Id())
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows_user.Close()

	for rows_user.Next() {
		return true
	}

	// Is shared to a group of this user
	rows_groups, err := db.driver.Query("select * from user_usergroup where user=? ", user.Id())
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows_groups.Close()

	for rows_groups.Next() {
		var group UserGroup
		rows_groups.Scan(&group.id)
		if this.IsSharedToUserGroup(&group) {
			return true
		}
	}
	return false

}

func (this *Pack) GetFiles() []*File {

	rows, err := db.driver.Query("SELECT * FROM file_pack WHERE pack=?", this.id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		println(id, name)
	}
	rows.Close()

	return nil

}

func (this Pack) IsSharedToUserGroup(group *UserGroup) bool {
	fmt.Println("packs: IsSharedToUserGroup group=", group.Name())

	// Is shared to this user
	rowsGroup, err := db.driver.Query("select * from pack_share_usergroup where group=? AND pack=?", group.Id(), this.Id())
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rowsGroup.Close()

	for rowsGroup.Next() {
		return true
	}

	// Is shared to a group of this user
	rows_groups, err := db.driver.Query("select * from usergroup_usergroup where parentGroup=? ", group.Id())
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows_groups.Close()

	for rows_groups.Next() {
		var childGroup UserGroup
		rows_groups.Scan(&childGroup.id)
		if this.IsSharedToUserGroup(&childGroup) {
			return true
		}
	}
	return false

}

///////////////////////
// Pack list
///////////////////////

type packs struct {
}

func (this packs) GetByPath(path string) *Pack {
	fmt.Println("packs: GetByPath path=", path)

	var pathSegments = strings.Split(path, "/")

	if len(pathSegments) == 0 {
		fmt.Println("packs: GetByPath no segments")
		return nil
	}

	var pack *Pack = nil

	var id, _ = strconv.Atoi(pathSegments[0])

	var user *User = model.Users.GetById(id)

	if user == nil {
		fmt.Println("packs: GetByPath wrong user id: " + pathSegments[0])
		return nil
	}

	pack = user.RootPack()

	for i := 1; i < len(pathSegments); i++ {
		var packId, _ = strconv.Atoi(pathSegments[i])
		pack = model.Packs.GetByParent(pack.Id(), packId)
		if pack == nil {
			break
		}
	}

	return pack
}

func (this packs) GetByParent(parentId int, id int) *Pack {
	fmt.Println("packs: GetByParent parentId=", parentId, " id=", id)

	rows, err := db.driver.Query("select * from pack where id=? AND parent=?", id, parentId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		return LoadPack(rows)
	}
	return nil

}

/* if user in self.allowedUsers.all():
       return True

   for group in self.allowedGroups.all():
       if group.contain_user(user):
           return True*/
