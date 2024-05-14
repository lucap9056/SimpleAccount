package AccountStruct

type Token struct {
	Id     int
	User   *User
	Token  string
	Device string
}
