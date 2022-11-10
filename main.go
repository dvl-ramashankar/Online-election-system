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

func SaveUserDetails(w http.ResponseWriter, r *http.Request) {
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
	if dataBody.Role == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter role field value Admin or Voter", "")
		return
	} else {
		if dataBody.Role != "Admin" && dataBody.Role != "Voter" {
			respondWithError(w, http.StatusBadRequest, "Please enter role field value Admin or Voter", "")
			return
		}
	}
	if result, msg, err := ser.SaveUserDetails(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func SearchUsersDetailsById(w http.ResponseWriter, r *http.Request) {
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

func UpdateUserDetailsById(w http.ResponseWriter, r *http.Request) {
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

func VerifyUser(w http.ResponseWriter, r *http.Request) {
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

func SearchUsersDetailsFilter(w http.ResponseWriter, r *http.Request) {
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

func DeactivateUser(w http.ResponseWriter, r *http.Request) {
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

func AddElection(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
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

	var dataBody model.ElectionRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := ser.AddElection(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func AddCandidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody model.CandidatesRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := ser.AddCandidate(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func verifyCandidate(w http.ResponseWriter, r *http.Request) {
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

	var dataBody model.VerifyCandidates
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := ser.VerifyCandidate(dataBody, mail); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func FindElectionById(w http.ResponseWriter, r *http.Request) {
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

	if result, msg, err := ser.FindElectionById(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func SearchFilterOnElectionDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}

	var dataBody model.SearchFilterElectionReq
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if result, msg, err := ser.SearchFilterOnElectionDetails(dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func UpdateElectionDetailsById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	_, role, err := validateToken(token)
	if role != "Admin" {
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

	var dataBody model.ElectionDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "Error Occurred")
		return
	}

	if result, msg, err := ser.UpdateElectionDetailsById(dataBody, id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func DeactivateElection(w http.ResponseWriter, r *http.Request) {
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

	if result, msg, err := ser.DeactivateElection(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func ElectionResultById(w http.ResponseWriter, r *http.Request) {
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

	if result, msg, err := ser.ElectionResultById(id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func CaseVoteByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid Method", "Invalid")
		return
	}
	token := r.Header.Get("tokenid")
	mailId, _, err := validateToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), "Error Occurred")
		return
	}
	var dataBody model.CastVoteReq
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request", "")
		return
	}

	if dataBody.ElectionId == "" || dataBody.CandidateId == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter electionId or CandidateId", "Error occurred")
		return
	}
	if result, msg, err := ser.CaseVote(dataBody, mailId); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err), msg)
	} else {
		respondWithJson(w, http.StatusAccepted, result, "", msg)
	}
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("MailId:", mail)
	fmt.Println("Role:", role)
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
	http.HandleFunc("/generate-token/", GenerateToken)
	http.HandleFunc("/user/create/", SaveUserDetails)
	http.HandleFunc("/user/search/", SearchUsersDetailsById)
	http.HandleFunc("/user/search-filter/", SearchUsersDetailsFilter)
	http.HandleFunc("/user/update/", UpdateUserDetailsById)
	http.HandleFunc("/user/verify/", VerifyUser)
	http.HandleFunc("/user/deactivate/", DeactivateUser)
	http.HandleFunc("/election/add_election/", AddElection)
	http.HandleFunc("/election/add_candidate/", AddCandidate)
	http.HandleFunc("/election/verify_candidate/", verifyCandidate)
	http.HandleFunc("/election/find_election/", FindElectionById)
	http.HandleFunc("/election/election_search/", SearchFilterOnElectionDetails)
	http.HandleFunc("/election/update/", UpdateElectionDetailsById)
	http.HandleFunc("/election/deactivate/", DeactivateElection)
	http.HandleFunc("/result/", ElectionResultById)
	http.HandleFunc("/vote/", CaseVoteByUser)
	log.Println("Server started at 8080")
	http.ListenAndServe(":8080", nil)
}
