package main

import "cloud.google.com/go/datastore"

// Account ...
// Describes a users' account
type Account struct {
	K *datastore.Key `datastore:"__key__"`
	Email	string	`json:"email"`
	Password string	`json:"password"`
	Admin bool `json:"isAdmin"`
}

// Game ...
// Describes a Game
type Game struct {
	K *datastore.Key `datastore:"__key__"`
	Owner	float64 `json:"owner"`
	Name	string	`json:"name"`
	Version	string	`json:"version"`
}