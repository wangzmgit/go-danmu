package dto

type RegisterDto struct {
	Name     string
	Email    string
	Password string
	Code     string
}

type LoginDto struct {
	Email    string
	Password string
}

type ModifyUserDto struct {
	Name     string
	Gender   int
	Birthday string
	Sign     string
}

type ModifyPasswordDto struct {
	Password string
	Code     string
}

type AdminModifyUserDto struct {
	ID    uint
	Name  string
	Email string
	Sign  string
}
