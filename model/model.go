package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Role         string             `bson:"role,omitempty" json:"role,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	MailId       string             `bson:"mail_id,omitempty" json:"mail_id,omitempty"`
	Password     string             `bson:"password,omitempty" json:"password,omitempty"`
	PhoneNumber  string             `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	IsVerified   bool               `bson:"is_verified" json:"is_verified"`
	PersonalInfo PersonalInfo       `bson:"personal_info,omitempty" json:"personal_info,omitempty"`
	VerifiedBy   VerifiedBy         `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
	UploadedDocs UploadedDocs       `bson:"uploaded_docs,omitempty" json:"uploaded_docs,omitempty"`
	Voted        Voted              `bson:"voted,omitempty" json:"voted,omitempty"`
	CreatedOn    time.Time          `bson:"created_on,omitempty" json:"created_on,omitempty"`
}

type PersonalInfo struct {
	FatherName   string      `bson:"father_name,omitempty" json:"father_name,omitempty"`
	Dob          time.Time   `bson:"dob,omitempty" json:"dob,omitempty"`
	Age          int64       `bson:"age,omitempty" json:"age,omitempty"`
	VoterId      string      `bson:"voter_id,omitempty" json:"voter_id,omitempty"`
	DocumentType string      `bson:"document_type,omitempty" json:"document_type,omitempty"`
	Address      AddressInfo `bson:"address,omitempty" json:"address,omitempty"`
}

type AddressInfo struct {
	City    string `bson:"city,omitempty" json:"city,omitempty"`
	Street  string `bson:"street,omitempty" json:"street,omitempty"`
	State   string `bson:"state,omitempty" json:"state,omitempty"`
	ZipCode string `bson:"zip_code,omitempty" json:"zip_code,omitempty"`
	Country string `bson:"country,omitempty" json:"country,omitempty"`
}

type VerifiedBy struct {
	Id   primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
}

type UploadedDocs struct {
	DocumentType             string `bson:"document_type,omitempty" json:"document_type,omitempty"`
	DocumentIdentificationNo string `bson:"document_identification_no,omitempty" json:"document_identification_no,omitempty"`
	DocumentPath             string `bson:"document_path,omitempty" json:"document_path,omitempty"`
}

type Voted struct {
	ElectionId []primitive.ObjectID `bson:"election_id,omitempty" json:"election_id,omitempty"`
}

type FinalResponse struct {
	Success    string      `json:"success,omitempty"`
	SucessCode string      `json:"successCode,omitempty"`
	SucessMsg  string      `json:"successMsg,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

type LoginDetails struct {
	MailId   string ` json:"mail_id,omitempty"`
	Password string ` json:"password,omitempty"`
}

type VerifyUserRequest struct {
	Id         string `bson:"_id,omitempty" json:"id,omitempty"`
	Role       string `bson:"role,omitempty" json:"role,omitempty"`
	Name       string `bson:"name,omitempty" json:"name,omitempty"`
	MailId     string `bson:"mail_id,omitempty" json:"mail_id,omitempty"`
	IsVerified bool   `bson:"is_verified,omitempty" json:"is_verified,omitempty"`
}

type SearchFilterRequest struct {
	Id           string       `bson:"_id,omitempty" json:"id,omitempty"`
	Role         string       `bson:"role,omitempty" json:"role,omitempty"`
	Name         string       `bson:"name,omitempty" json:"name,omitempty"`
	MailId       string       `bson:"mail_id,omitempty" json:"mail_id,omitempty"`
	Password     string       `bson:"password,omitempty" json:"password,omitempty"`
	PhoneNumber  string       `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	IsVerified   bool         `bson:"is_verified,omitempty" json:"is_verified,omitempty"`
	PersonalInfo PersonalInfo `bson:"personal_info,omitempty" json:"personal_info,omitempty"`
	VerifiedBy   VerifiedBy   `bson:"verified_by,omitempty" json:"verified_by,omitempty"`
}

type ElectionDetails struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Location       string             `bson:"location,omitempty" json:"location,omitempty"`
	ElectionDate   time.Time          `bson:"election_date,omitempty" json:"election_date,omitempty"`
	ResultDate     time.Time          `bson:"result_date,omitempty" json:"result_date,omitempty"`
	Result         string             `bson:"result,omitempty" json:"result,omitempty"`
	ElectionStatus string             `bson:"election_status,omitempty" json:"election_status,omitempty"`
	Candidates     []Candidates       `bson:"candidates,omitempty" json:"candidates,omitempty"`
}

type Candidates struct {
	UserId                primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name                  string             `bson:"name" json:"name"`
	Commitments           []string           `bson:"commitments" json:"commitments"`
	Votecount             string             `bson:"vote_count" json:"vote_count"`
	VoteSign              string             `bson:"vote_sign" json:"vote_sign"`
	NominationStatus      string             `bson:"nomination_status" json:"nomination_status"`
	NomainationVerifiedBy primitive.ObjectID `bson:"nomaination_verified_by,omitempty" json:"nomaination_verified_by,omitempty"`
}

type ElectionRequest struct {
	Location       string `json:"location,omitempty"`
	ElectionDate   string `json:"election_date,omitempty"`
	ResultDate     string `json:"result_date,omitempty"`
	ElectionStatus string `json:"election_status,omitempty"`
}

type CandidatesRequest struct {
	ElectionId  string   `json:"election_id,omitempty"`
	UserId      string   `json:"user_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Commitments []string `json:"commitments,omitempty"`
	VoteSign    string   `json:"vote_sign,omitempty"`
}

type VerifyCandidates struct {
	ElectionId       string `json:"election_id"`
	UserId           string `json:"user_id"`
	NominationStatus string `json:"nomination_status"`
}

type SearchFilterElectionReq struct {
	Id             string `json:"_id,omitempty"`
	Location       string `json:"location,omitempty"`
	ElectionDate   string `json:"election_date,omitempty"`
	ResultDate     string `json:"result_date,omitempty"`
	Result         string `json:"result,omitempty"`
	ElectionStatus string `json:"election_status,omitempty"`
	CandidateName  string `json:"election_status,omitempty"`
}
