package service

import (
	"election/model"
	"fmt"
	"testing"
)

var req []*model.User

type Tests struct {
	name     string
	request  []*model.User
	response bool
}

func TestConvertDate(t *testing.T) {
	result, err := ConvertDate("2022-10-15")
	want := "2022-10-15 00:00:00 +0000 UTC"
	if err != nil {
		t.Errorf("Conversion Failed. Error Occurred %s", err)
	} else {
		if result.String() != want {
			t.Errorf("Conversion Failed. Expected %s ,Got %s\n", want, result.String())
		} else {
			t.Logf("Conversion Pass. Expected %s, Got %s\n", want, result.String())
		}
	}
}

func TestConvertStringIntoHex(t *testing.T) {
	result, err := ConvertStringIntoHex("636266374b19404e533683d1")
	id := "636266374b19404e533683d1"
	want := fmt.Sprintf("ObjectID(%q)", id)
	if err != nil {
		t.Errorf("Conversion Failed. Error Occurred %s", err)
	} else {
		if result.String() != want {
			t.Errorf("Conversion Failed. Expected %s ,Got %s\n", want, result.String())
		} else {
			t.Logf("Conversion Pass. Expected %s, Got %s\n", want, result.String())
		}
	}
}
