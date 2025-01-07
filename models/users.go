package models

type User struct {
	Login    string `json:"login"`
	Password string `json:"-"`
	ID       int    `json:"id"`
	ReferUID int    `json:"refer_uid"`
	Points   int    `json:"points"`
}
