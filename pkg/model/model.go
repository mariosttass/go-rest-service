package model

type Object struct {
	ID        int    `json:"id"`
	Number    int    `json:"number"`
	Online    bool   `json:"online"`
	CreatedAt string `json:"created_at"`
}

type ObjectList struct {
	Objects []Object `json:"objects"`
}
