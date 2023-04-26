package model

type UserProfile struct {
	Cookie string `json:"cookie"`
	// sorted in descending time order
	Views []UserTagEvent `json:"views"`
	// sorted in descending time order
	Buys []UserTagEvent `json:"buys"`
}