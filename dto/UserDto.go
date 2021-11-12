package dto

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
	Code     string
}

type LoginRequest struct {
	Email    string
	Password string
}

type UserModifyRequest struct {
	Name     string
	Gender   int
	Birthday string
	Sign     string
}

type PassModifyRequest struct {
	Password string
	Code     string
}

type AdminModifyUserRequest struct {
	ID    uint
	Name  string
	Email string
	Sign  string
}
