package dto

type OssConfigDto struct {
	Storage         bool
	Bucket          string
	Endpoint        string
	AccesskeyId     string
	AccesskeySecret string
	Domain          string
}

type EmailConfigDto struct {
	Name     string
	Host     string
	Port     int
	Address  string
	Password string
}

type AdminConfigDto struct {
	Email    string
	Password string
}
