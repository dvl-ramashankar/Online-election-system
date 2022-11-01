package main

import (
	"election/auth"
	"election/model"
	"election/service"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var ser = service.Connection{}
var finalResponse model.FinalResponse

func init() {
	ser.Server = "mongodb://localhost:27017"
	ser.Database = "ElectionSystem"
	ser.Collection1 = "UserDetails"
	ser.Collection2 = "ElectionDetails"

	ser.Connect()
}

func saveUserDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if dataBody.MailId == "" || dataBody.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter mailId or Password", "")
		return
	}
	if dataBody.Role == "Admin" || dataBody.Role == "Voter" {
		respondWithError(w, http.StatusBadRequest, "Please enter role field value Admin or Voter", "")
		return
	}
	if result, msg, err := ser.SaveUserDetails(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func searchUsersDetailsById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	segment := strings.Split(r.URL.Path, "/")
	id := segment[len(segment)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide Id for Search", "Error Occurred")
	}

	if result, msg, err := ser.SearchUsersDetailsById(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func updateUserDetailsById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	_, role, err := validateToken(token)
	if role != "Voter" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	path := r.URL.Path
	segments := strings.Split(path, "/")
	id := segments[len(segments)-1]

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "Error Occurred")
		return
	}

	if result, msg, err := ser.UpdateUserDetailsById(dataBody, id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func verifyUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	mail, role, err := validateToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	if role != "Admin" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}

	var dataBody model.VerifyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := ser.VerifyUser(dataBody, mail); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func searchUsersDetailsFilter(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody model.SearchFilterRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := ser.FilterOnUsersDetails(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func deactivateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	_, role, err := validateToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	if role != "Admin" {
		respondWithError(w, http.StatusBadRequest, "Token is invalid as it's role is different", "Invalid")
		return
	}

	segment := strings.Split(r.URL.Path, "/")
	id := segment[len(segment)-1]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide Id for Search", "Error Occurred")
	}

	if result, msg, err := ser.DeactivateUser(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var loginDetails model.LoginDetails
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "")
	}

	if result, msg, err := ser.GenerateToken(loginDetails); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func validateToken(token string) (string, string, error) {
	if token == "" {
		return "", "", errors.New("Please Enter Token")
	}
	mail, err := auth.ValidateToken(token)
	if err != nil {
		return "", "", errors.New("Either Token Is Invalid Or Expired")
	}
	role := ser.FetchRole(mail)
	return mail, role, err
}

func respondWithError(w http.ResponseWriter, code int, msg string, msg2 string) {
	respondWithJson(w, code, map[string]string{"error": msg}, "error", msg2)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}, err, msg string) {
	if err == "error" {
		finalResponse.Success = "false"
	} else {
		finalResponse.Success = "true"
	}
	finalResponse.SucessMsg = msg
	finalResponse.SucessCode = fmt.Sprintf("%v", code)
	finalResponse.Response = payload
	response, _ := json.Marshal(finalResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	http.HandleFunc("/generate-token", generateToken)
	http.HandleFunc("/user/create", saveUserDetails)
	http.HandleFunc("/user/search/", searchUsersDetailsById)
	http.HandleFunc("/user/search-filter/", searchUsersDetailsFilter)
	http.HandleFunc("/user/update/", updateUserDetailsById)
	http.HandleFunc("/user/verify/", verifyUser)
	http.HandleFunc("/user/deactivate/", deactivateUser)
	log.Println("Server started at 8080")
	http.ListenAndServe(":8080", nil)
}
