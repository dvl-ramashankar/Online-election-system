package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var TokenIdAdmin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJrQGdtYWlsLmNvbSIsImV4cCI6MTY2ODAwNDQ3OH0.hKVC6nxL3NyZanVplLCrdmNrrIUCn4d0VpanMHlX7s4"
var TokenIdUser = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJrQGdtYWlsLmNvbSIsImV4cCI6MTY2ODAwNDQ3OH0.hKVC6nxL3NyZanVplLCrdmNrrIUCn4d0VpanMHlX7s4"

func init() {
	ser.Server = "mongodb://localhost:27017"
	ser.Database = "ElectionSystem"
	ser.Collection1 = "UserDetails_test"
	ser.Collection2 = "ElectionDetails_test"

	ser.Connect()
}

func TestGenerateToken(t *testing.T) {
	request := `{
		"mail_id": "test14324@gmail.com",
		"password": "test1"
	}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(request))
	w := httptest.NewRecorder()
	GenerateToken(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestSaveUserDetails(t *testing.T) {
	request := `{
		"role" : "Admin",
		"name": "test12345",
		"mail_id": "rk@gmail.com",
		"password": "rk124323",
		"phone_number": "1234554326",
		"personal_info":{
			"father_name":"demo",
			"dob":"2001-10-01T00:00:00Z",
			"address":{
				"street" :"MgfdsH"
			}
		},
		"uploaded_docs": {
			"document_type": "Lorem",
			"document_identification_no": "Lorem",
			"document_path": "D:/Back_Up/529_3_334363_1666252880_Databricks - Generic.pdf"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(request))
	w := httptest.NewRecorder()
	SaveUserDetails(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestSearchUsersDetailsById(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/635a13aa96d74d1397c5bd9d", nil)
	w := httptest.NewRecorder()
	SearchUsersDetailsById(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestSearchUsersDetailsFilter(t *testing.T) {
	request := `{"id": "635a13aa96d74d1397c5bd9d"}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(request))
	w := httptest.NewRecorder()
	SearchUsersDetailsFilter(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestUpdateUserDetailsById(t *testing.T) {
	request := `{
		"name": "test1Update"
	}`
	req := httptest.NewRequest(http.MethodPut, "/636a486f76afdab19f61a19b", bytes.NewBufferString(request))
	req.Header.Set("tokenid", TokenIdUser)
	w := httptest.NewRecorder()
	UpdateUserDetailsById(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestFindElectionById(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/636266374b19404e533683d1", nil)
	w := httptest.NewRecorder()
	FindElectionById(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestDeactivateElection(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/636266374b19404e533683", nil)
	req.Header.Set("tokenid", TokenIdAdmin)
	w := httptest.NewRecorder()
	DeactivateElection(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestElectionResultById(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/636a08a6e95e88f80db6ecbc", nil)
	w := httptest.NewRecorder()
	ElectionResultById(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestVerifyUser(t *testing.T) {
	request := `{
	"id": "636a486f76afdab19f61a19b",
	"is_verified":true
	}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(request))
	req.Header.Set("tokenid", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJrQGdtYWlsLmNvbSIsImV4cCI6MTY2NzkxNDcxOX0.FGdopauErehATAf8G8WNLj99TVWdGYBNGzSH7YCgRMA")
	w := httptest.NewRecorder()
	VerifyUser(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestDeactivateUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/636a486f76afdab19f61a19b", nil)
	req.Header.Set("tokenid", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InJrQGdtYWlsLmNvbSIsImV4cCI6MTY2NzkxNDcxOX0.FGdopauErehATAf8G8WNLj99TVWdGYBNGzSH7YCgRMA")
	w := httptest.NewRecorder()
	DeactivateUser(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestAddElection(t *testing.T) {
	request := `{
		"location": "Gurugram",
		"election_status": "Nomination Started",
		"election_date": "2022-10-10",
		"result_date": "2023-10-15"
	}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(request))
	req.Header.Set("tokenid", TokenIdAdmin)
	w := httptest.NewRecorder()
	AddElection(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestAddCandidate(t *testing.T) {
	request := `{
			"name":"rk",
			"user_id" :"636b7a8974c7cd15a324b433",
			"election_id" :"636bacd841a4e2396f912ceb",
			"commitments":["dumm4","dummy","htrds"],
			"vote_sign":"Chakra"
		}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(request))
	req.Header.Set("tokenid", TokenIdAdmin)
	w := httptest.NewRecorder()
	AddCandidate(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestSearchFilterOnElectionDetails(t *testing.T) {
	request := `{
		"location":"Mumbai"
		}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(request))
	w := httptest.NewRecorder()
	SearchFilterOnElectionDetails(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
func TestUpdateElectionDetailsById(t *testing.T) {
	request := `{
		"election_status": "Voting In-Progress"
	}`
	req := httptest.NewRequest(http.MethodPut, "/636a08a6e95e88f80db6ecbc", bytes.NewBufferString(request))

	w := httptest.NewRecorder()
	UpdateElectionDetailsById(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode == http.StatusAccepted {
		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
	} else {
		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
	}
}
