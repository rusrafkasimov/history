package dto

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
} // @Name Error
