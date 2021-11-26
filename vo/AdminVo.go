package vo

type AdminListVo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Authority int    `json:"authority"`
}

type AdminVo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Authority int    `json:"authority"`
}
