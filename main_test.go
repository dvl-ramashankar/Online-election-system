package main

import (
	// "encoding/json"
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SearchReq struct {
	Id string
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
	req := httptest.NewRequest(http.MethodGet, "/636266374b19404e533683", nil)
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

// func TestSearchUsersDetailsFilter(t *testing.T) {
// 	request := `{"id": "635a13aa96d74d1397c5bd9d"}`
// 	var cod SearchReq
// 	err := json.Unmarshal([]byte(request), &cod)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	req := httptest.NewRequest(http.MethodPost, "", nil)
// 	w := httptest.NewRecorder()
// 	SearchUsersDetailsFilter(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()
// 	if res.StatusCode == http.StatusAccepted {
// 		t.Logf("Pass. Expected %d, Got %d\n", http.StatusAccepted, res.StatusCode)
// 	} else {
// 		t.Errorf("Failed. Expected %d ,Got %d\n", http.StatusAccepted, res.StatusCode)
// 	}
// }
