package users

import "time"

// Only the Username type should be passed around.
type Username string

type Session struct {
	Key []byte
	TTL time.Time
}

// This should only be formed when querying the user hash
type User struct {
	Username     Username `json:"username"`
	PasswordHash string   `json:"password-hash"`
	Session      []Session
}

// func VerifyCredentials(login User) bool {

// }
