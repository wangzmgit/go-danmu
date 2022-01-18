package dto

type AdminIdDto struct {
	ID uint
}

type AdminLoginDto struct {
	Email    string
	Password string
}

type AddAdminDto struct {
	Name      string
	Email     string
	Password  string
	Authority int
}
