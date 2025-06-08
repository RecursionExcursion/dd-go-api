package betbot

type CompressedFsData struct {
	Id      string `json:"id"`
	Created int64  `json:"created"`
	Data    string `json:"data"`
}

type User struct {
	MId      string `json:"_id"`
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
