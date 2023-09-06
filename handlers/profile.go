package handlers

import (
	"encoding/json"
	"net/http"
)

type userInfo struct {
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

func (bm *BookManagerServer) HandleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Grab Authorization Header
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve the related account to the token
	username, err := bm.Authenticate.GetUsernameByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve user from database
	user, err := bm.Db.GetUserByUsername(*username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the response body
	res, err := json.Marshal(&userInfo{
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
