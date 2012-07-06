package yack

type Part struct {
	id     int
	offset int
	size   int
	file   string
	sha    string
}

func (this *Part) Size() int {
	return this.size
}
