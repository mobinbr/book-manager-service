package handlers

import (
	"BookManager/authenticate"
	"BookManager/db"
	"encoding/json"
	"io"
	"net/http"
)

type signupRequestBody struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
}

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (bm *BookManagerServer) HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body for the new user
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var srb signupRequestBody
	err = json.Unmarshal(reqData, &srb)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not unmarshal the new user request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add the user to the database
	user := &db.User{
		Username:    srb.Username,
		Email:       srb.Email,
		Password:    srb.Password,
		FirstName:   srb.FirstName,
		LastName:    srb.LastName,
		PhoneNumber: srb.PhoneNumber,
		Gender:      srb.Gender,
	}
	err = bm.Db.CreateNewUser(user)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not create a new user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create the response body
	response := map[string]interface{}{
		"message": "user has been created successfully",
	}
	resBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusAccepted)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func (bm *BookManagerServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body for the new user
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var lrb loginRequestBody
	err = json.Unmarshal(reqData, &lrb)
	if err != nil {
		bm.Logger.WithError(err).Warn("can not unmarshal the login request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check the authentication
	token, err := bm.Authenticate.Login(authenticate.Credentials{
		Username: lrb.Username,
		Password: lrb.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Create the response body
	res := map[string]interface{}{
		"access_token": token.TokenString,
	}
	resBody, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
