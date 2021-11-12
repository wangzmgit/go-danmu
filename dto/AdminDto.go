package dto

type AdminIDRequest struct {
	ID uint
}

type AdminLoginRequest struct {
	Email    string
	Password string
}

type AddAdminRequest struct {
	Name      string
	Email     string
	Password  string
	Authority int
}
