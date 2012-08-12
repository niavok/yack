package yack

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Pack struct {
	id           int64
	name         string
	creationDate time.Time
    ownerId      int64
	parentPackId   int64
	isPublic     bool
	
	//Cache
	owner        *User
	parentPack   *Pack
}


func NewPack(name string) *Pack {
    fmt.Println("NewPack()")
	result, err := db.driver.Exec("insert into pack ('name', 'creationDate', 'isPublic', 'owner', 'parentPack') values(?, ?, ?, ?, ?);", name, time.Now(), false, -1, -1)
	if err != nil {
		fmt.Println(err)
		return nil
	}
    id, _ := result.LastInsertId()
    fmt.Println("LastInsertId pack ", id)
	var pack = model.Packs.GetById(id)
	return pack;
}


func LoadPack(row *sql.Rows) *Pack {
	fmt.Println("LoadPack()")
	var pack Pack
	var creationDate string

	err := row.Scan(&pack.id, &pack.name, &creationDate, &pack.ownerId, &pack.parentPackId, &pack.isPublic)
	if err != nil {
		fmt.Println("LoadPack Scan Error: ", err)
		return nil
	}
	pack.creationDate, _ = time.Parse("2006-01-02 15:04:05", creationDate)
	
	return &pack
}

func (this *Pack) Id() int64 {
	return this.id
}

func (this *Pack) Owner() *User {
    if this.owner == nil {
        this.owner = model.Users.GetById(this.ownerId)
    }
	return this.owner
}

func (this *Pack) CanRead(user *User) bool {
	if model.Users.Equal(this.Owner(),user) {
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

	rows, err := db.driver.Query("SELECT id, name, creationDate, size, uploadedSize, sha, uploadState, file.file, owner, description, mime, autoMime FROM pack_file, file WHERE pack=?", this.id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	
	files := []*File{}
	
	for rows.Next() {
        files = append(files, LoadFile(rows))
	}

	return files
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

func (this Pack) SetOwner(user *User) {
    _, err := db.driver.Exec("UPDATE pack SET owner=? WHERE id=?", user.Id(), this.Id())
	if err != nil {
		fmt.Println("Error in SetOwner: ",err)
	}
	
	this.ownerId = user.Id()
    this.owner = user
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

	var id, _ = strconv.ParseInt(pathSegments[0], 10, 64)

	var user *User = model.Users.GetById(id)

	if user == nil {
		fmt.Println("packs: GetByPath wrong user id: " + pathSegments[0])
		return nil
	}

	pack = user.RootPack()
	fmt.Println("packs: GetByPath RootPack: ", pack)

	for i := 1; i < len(pathSegments); i++ {
		var packId, _ = strconv.ParseInt(pathSegments[i], 10, 64)
		pack = model.Packs.GetByParent(pack.Id(), packId)
		if pack == nil {
			break
		}
	}

	return pack
}

func (this packs) GetByParent(parentId int64, id int64) *Pack {
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

func (this packs) GetById(id int64) *Pack {
	fmt.Println("packs: Get id=", id)

	rows, err := db.driver.Query("select * from pack where id=?", id)
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
