package service

import (
	"context"
	"election/auth"
	"election/model"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	Server      string
	Database    string
	Collection1 string
	Collection2 string
}

const uploadPath = "upload/"

var CollectionUserDetails *mongo.Collection
var CollectionElectionDetails *mongo.Collection
var ctx = context.TODO()

func (e *Connection) Connect() {
	clientOptions := options.Client().ApplyURI(e.Server)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	CollectionUserDetails = client.Database(e.Database).Collection(e.Collection1)
	CollectionElectionDetails = client.Database(e.Database).Collection(e.Collection2)
}

// ===================================userDetails============================================
func (e *Connection) SaveUserDetails(reqBody model.User) ([]*model.User, string, error) {
	var data []*model.User
	bool, err := ValidateByMailId(reqBody.MailId)
	if err != nil {
		return data, "", err
	}
	if !bool {
		return data, "", errors.New("User already present")
	}
	if err != nil {
		log.Println(err)
		return data, "", err
	}
	msg, err := UploadFile(reqBody.UploadedDocs.DocumentPath)
	if err != nil {
		log.Println(err)
		return data, "", errors.New("Unable to upload file")
	}
	fmt.Println("Upload file:", msg)
	reqBody.IsVerified = false
	finalData, err := CollectionUserDetails.InsertOne(ctx, reqBody)
	if err != nil {
		log.Println(err)
		return data, "", errors.New("Unable to store data")
	}
	result, err := CollectionUserDetails.Find(ctx, bson.D{primitive.E{Key: "_id", Value: finalData.InsertedID}})
	if err != nil {
		log.Println(err)
		return data, "", err
	}
	data, err = ConvertDbResultIntoUserStruct(result)
	if err != nil {
		log.Println(err)
		return data, "", err
	}
	return data, "Saved Successfully", nil
}

func (e *Connection) SearchUsersDetailsById(idStr string) ([]*model.User, string, error) {
	var finalData []*model.User

	id, err := ConvertStringIntoHex(idStr)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	searchData, err := CollectionUserDetails.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = ConvertDbResultIntoUserStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	return finalData, "Data Fetch Successfully", nil
}

func (e *Connection) UpdateUserDetailsById(reqData model.User, idStr string) (bson.M, string, error) {
	var updatedDocument bson.M
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return updatedDocument, "Error Occurred", err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	UpdateQuery := bson.D{}
	if reqData.Name != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "name", Value: reqData.Name})
	}
	if reqData.Password != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "password", Value: reqData.Password})
	}
	if reqData.PhoneNumber != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "phone_number", Value: reqData.PhoneNumber})
	}
	if reqData.PersonalInfo.FatherName != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.father_name", Value: reqData.PersonalInfo.FatherName})
	}
	if reqData.PersonalInfo.Age != 0 {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.age", Value: reqData.PersonalInfo.Age})
	}
	if reqData.PersonalInfo.DocumentType != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.document_type", Value: reqData.PersonalInfo.DocumentType})
	}
	if reqData.PersonalInfo.Address.City != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.city", Value: reqData.PersonalInfo.Address.City})
	}
	if reqData.PersonalInfo.Address.Street != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.street", Value: reqData.PersonalInfo.Address.Street})
	}
	if reqData.PersonalInfo.Address.ZipCode != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.zip_code", Value: reqData.PersonalInfo.Address.ZipCode})
	}
	if reqData.PersonalInfo.Address.State != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.state", Value: reqData.PersonalInfo.Address.State})
	}
	if reqData.PersonalInfo.Address.Country != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "personal_info.address.country", Value: reqData.PersonalInfo.Address.Country})
	}

	update := bson.D{primitive.E{Key: "$set", Value: UpdateQuery}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	r := CollectionUserDetails.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if r != nil {
		return updatedDocument, "Error Occurred", r
	}

	if updatedDocument == nil {
		return updatedDocument, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}

	return updatedDocument, "Document Updated Successfully", nil
}

func (e *Connection) VerifyUser(req model.VerifyUserRequest, adminMail string) ([]*model.User, string, error) {
	var finalData []*model.User
	var adminData []*model.User

	data, err := CollectionUserDetails.Find(ctx, bson.D{primitive.E{Key: "mail_id", Value: adminMail}})
	adminData, err = ConvertDbResultIntoUserStruct(data)
	if len(adminData) == 0 {
		return finalData, "Error Occurred", errors.New("Data not present in db acc. to given tokenId")
	}
	filter := bson.D{}
	flag := true
	if req.Id != "" {
		id, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		filter = append(filter, primitive.E{Key: "_id", Value: id})
		flag = false
	}
	if flag {
		if req.MailId != "" {
			filter = append(filter, primitive.E{Key: "mail_id", Value: bson.M{"$regex": req.MailId}})
			flag = false
		}
	}
	UpdateQuery := bson.D{}
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "is_verified", Value: req.IsVerified})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "verified_by.id", Value: adminData[0].Id})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "verified_by.name", Value: adminData[0].Name})
	update := bson.D{primitive.E{Key: "$set", Value: UpdateQuery}}

	CollectionUserDetails.FindOneAndUpdate(ctx, filter, update)

	data, err = CollectionUserDetails.Find(ctx, filter)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	finalData, err = ConvertDbResultIntoUserStruct(data)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	//Send mail method
	return finalData, "Voter verified successfully!", nil
}

func (e *Connection) FilterOnUsersDetails(req model.SearchFilterRequest) ([]*model.User, string, error) {
	var finalData []*model.User
	query := bson.D{}

	if req.Id != "" {
		id, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		query = append(query, primitive.E{Key: "_id", Value: id})
	}
	if req.Name != "" {
		query = append(query, primitive.E{Key: "name", Value: req.Name})
	}
	if req.Role != "" {
		query = append(query, primitive.E{Key: "role", Value: req.Role})
	}
	if req.MailId != "" {
		query = append(query, primitive.E{Key: "mail_id", Value: req.MailId})
	}
	if req.IsVerified != false {
		query = append(query, primitive.E{Key: "is_verified", Value: req.IsVerified})
	}
	if req.PhoneNumber != "" {
		query = append(query, primitive.E{Key: "phone_number", Value: req.PhoneNumber})
	}
	if req.PersonalInfo.FatherName != "" {
		query = append(query, primitive.E{Key: "personal_info.father_name", Value: req.PersonalInfo.FatherName})
	}
	if req.PersonalInfo.Age != 0 {
		query = append(query, primitive.E{Key: "personal_info.age", Value: req.PersonalInfo.Age})
	}
	if req.PersonalInfo.DocumentType != "" {
		query = append(query, primitive.E{Key: "personal_info.document_type", Value: req.PersonalInfo.DocumentType})
	}
	if req.PersonalInfo.Address.City != "" {
		query = append(query, primitive.E{Key: "personal_info.address.city", Value: req.PersonalInfo.Address.City})
	}
	if req.PersonalInfo.Address.Street != "" {
		query = append(query, primitive.E{Key: "personal_info.address.street", Value: req.PersonalInfo.Address.Street})
	}
	if req.PersonalInfo.Address.ZipCode != "" {
		query = append(query, primitive.E{Key: "personal_info.address.zip_code", Value: req.PersonalInfo.Address.ZipCode})
	}
	if req.PersonalInfo.Address.State != "" {
		query = append(query, primitive.E{Key: "personal_info.address.state", Value: req.PersonalInfo.Address.State})
	}
	if req.PersonalInfo.Address.Country != "" {
		query = append(query, primitive.E{Key: "personal_info.address.country", Value: req.PersonalInfo.Address.Country})
	}
	if req.VerifiedBy.Name != "" {
		query = append(query, primitive.E{Key: "verified_by.name", Value: req.VerifiedBy.Name})
	}

	searchData, err := CollectionUserDetails.Find(ctx, query)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = ConvertDbResultIntoUserStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	return finalData, "Data Fetch Successfully", nil
}

func (e *Connection) DeactivateUser(idStr string) (bson.M, string, error) {
	var updatedDocument bson.M
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return updatedDocument, "Error Occurred", err
	}

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"_id", id}},
				bson.D{{"is_verified", true}},
			},
		},
	}
	update := bson.D{{"$set", bson.D{primitive.E{Key: "is_verified", Value: false}}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	r := CollectionUserDetails.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if r != nil {
		return updatedDocument, "Error Occurred", r
	}

	if updatedDocument == nil {
		return updatedDocument, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}

	return updatedDocument, "User details deactivate successfully!", nil
}

func ConvertDbResultIntoUserStruct(fetchDataCursor *mongo.Cursor) ([]*model.User, error) {
	var finaldata []*model.User
	for fetchDataCursor.Next(ctx) {
		var data model.User
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = append(finaldata, &data)
	}
	return finaldata, nil
}

func ValidateByMailId(mailId string) (bool, error) {

	var result []*model.User
	data, err := CollectionUserDetails.Find(ctx, bson.D{primitive.E{Key: "mail_id", Value: mailId}})
	if err != nil {
		return false, err
	}
	result, err = ConvertDbResultIntoUserStruct(data)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return true, err
	}
	return false, err
}

func ConvertDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Println(err)
		return date, err
	}
	return date, err
}

func UploadFile(path string) (string, error) {
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	fileURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	segments := strings.Split(fileURL.Path, "/")
	fileName := segments[len(segments)-1]
	fileName = uploadPath + fileName
	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	resp, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Close()
	size, err := io.Copy(file, resp)
	defer file.Close()
	return "File Downloaded with size :" + fmt.Sprintf("%v", size), nil
}

func ConvertStringIntoHex(idStr string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return id, err
	}
	return id, err
}

func (e *Connection) FetchRole(mailId string) string {
	data, err := CollectionUserDetails.Find(ctx, bson.D{primitive.E{Key: "mail_id", Value: mailId}})
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	convData, err := ConvertDbResultIntoUserStruct(data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	if len(convData) == 0 {
		return ""
	}
	return convData[0].Role
}

// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
// ======================================Election=============================================
func (e *Connection) AddElection(data model.ElectionRequest) ([]*model.ElectionDetails, string, error) {

	var finalData []*model.ElectionDetails
	setData, err := SetValueInElectionModelStruct(data)
	if err != nil {
		return finalData, "Error occurred", err
	}
	insert, err := CollectionElectionDetails.InsertOne(ctx, setData)
	if err != nil {
		return finalData, "Error occurred", err
	}
	fetchData, err := CollectionElectionDetails.Find(ctx, bson.D{primitive.E{Key: "_id", Value: insert.InsertedID}})

	finalData, err = ConvertDbResultIntoElectionStruct(fetchData)
	if err != nil {
		return finalData, "Error occurred", err
	}
	return finalData, "Election details saved successfully!", nil
}

func (e *Connection) AddCandidate(data model.CandidatesRequest) ([]*model.ElectionDetails, string, error) {
	var finalData []*model.ElectionDetails
	electionId, err := primitive.ObjectIDFromHex(data.ElectionId)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userId, err := primitive.ObjectIDFromHex(data.UserId)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userData, err := CollectionUserDetails.Find(ctx, bson.D{primitive.E{Key: "_id", Value: userId}})
	if err != nil {
		return finalData, "Error Occurred", err
	}
	d, _ := ConvertDbResultIntoUserStruct(userData)

	if len(d) == 0 {
		return finalData, "Error Occurred", errors.New("Invalid UserId")
	}

	filter := bson.D{primitive.E{Key: "_id", Value: electionId}}
	UpdateQuery := bson.D{}
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "user_id", Value: userId})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "name", Value: data.Name})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "commitments", Value: data.Commitments})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "vote_sign", Value: data.VoteSign})
	UpdateQuery = append(UpdateQuery, primitive.E{Key: "nomation_status", Value: "Verification In Progress"})

	update := bson.D{primitive.E{Key: "candidates", Value: UpdateQuery}}
	update = bson.D{primitive.E{Key: "$push", Value: update}}

	CollectionElectionDetails.FindOneAndUpdate(ctx, filter, update)

	fetchData, err := CollectionElectionDetails.Find(ctx, filter)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	finalData, err = ConvertDbResultIntoElectionStruct(fetchData)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	return finalData, "Candidates details saved successfully!", nil
}

func (e *Connection) VerifyCandidate(req model.VerifyCandidates, adminMail string) (bson.M, string, error) {
	var finalData bson.M
	var adminData []*model.User
	var elections []*model.ElectionDetails

	data, err := CollectionUserDetails.Find(ctx, bson.D{primitive.E{Key: "mail_id", Value: adminMail}})
	adminData, err = ConvertDbResultIntoUserStruct(data)
	if len(adminData) == 0 {
		return finalData, "Error Occurred", errors.New("Data not present in db acc. to given tokenId")
	}
	filter := bson.D{}
	electionId, err := primitive.ObjectIDFromHex(req.ElectionId)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return finalData, "Error Occurred", err
	}

	cur, err := CollectionElectionDetails.Find(ctx, bson.D{primitive.E{Key: "_id", Value: electionId}})
	if err != nil {
		return finalData, "Error Occurred", errors.New("unable to query db")
	}
	elections, err = ConvertDbResultIntoElectionStruct(cur)
	Candidates := elections[0].Candidates
	for i := range Candidates {
		if Candidates[i].UserId == userId {
			Candidates[i].NominationStatus = req.NominationStatus
			Candidates[i].NominationVerifiedBy = adminData[0].Id
		}
	}

	filter = bson.D{primitive.E{Key: "_id", Value: electionId}}
	update := bson.D{primitive.E{Key: "$set", Value: elections[0]}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = CollectionElectionDetails.FindOneAndUpdate(ctx, filter, update, opts).Decode(&finalData)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	fmt.Println(finalData)
	if finalData == nil {
		return finalData, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}
	return finalData, "Candidates verified successfully!", nil
}

func (e *Connection) FindElectionById(idStr string) ([]*model.ElectionDetails, string, error) {
	var finalData []*model.ElectionDetails

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"_id", id}},
				bson.D{{"election_status", bson.M{"$ne": "Deactivated"}}},
			},
		},
	}

	searchData, err := CollectionElectionDetails.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = ConvertDbResultIntoElectionStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	if finalData == nil {
		return finalData, "Data Fetched Successfully", errors.New("Either data is not present or it is deactivated")
	}
	return finalData, "Data Fetched Successfully", nil
}

func (e *Connection) SearchFilterOnElectionDetails(req model.SearchFilterElectionReq) ([]*model.ElectionDetails, string, error) {
	var finalData []*model.ElectionDetails
	query := bson.D{}

	if req.Id != "" {
		id, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		query = append(query, primitive.E{Key: "_id", Value: id})
	}
	if req.Location != "" {
		query = append(query, primitive.E{Key: "location", Value: req.Location})
	}
	if req.Result != "" {
		query = append(query, primitive.E{Key: "result", Value: req.Result})
	}
	if req.ElectionStatus != "" {
		query = append(query, primitive.E{Key: "election_status", Value: req.ElectionStatus})
	}
	if req.ElectionDate != "" {
		electionDate, err := ConvertDate(req.ElectionDate)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		query = append(query, primitive.E{Key: "election_date", Value: electionDate})
	}
	if req.ResultDate != "" {
		resultDate, err := ConvertDate(req.ResultDate)
		if err != nil {
			return finalData, "Error Occurred", err
		}
		query = append(query, primitive.E{Key: "result_date", Value: resultDate})
	}

	filter := bson.D{
		{"$and",
			bson.A{
				query,
				bson.D{{"election_status", bson.M{"$ne": "Deactivated"}}},
			},
		},
	}
	searchData, err := CollectionElectionDetails.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	finalData, err = ConvertDbResultIntoElectionStruct(searchData)
	if err != nil {
		log.Println(err)
		return finalData, "Error Occurred", err
	}
	if finalData == nil {
		return finalData, "Data Fetched Successfully", errors.New("Either data is not present or it is deactivated")
	}
	return finalData, "Data Fetched Successfully", nil
}

func (e *Connection) UpdateElectionDetailsById(reqData model.ElectionDetailsRequest, idStr string) (bson.M, string, error) {
	var finalData bson.M
	data, _, err := e.FindElectionById(idStr)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	electionStatus := data[0].ElectionStatus
	if electionStatus == "Completed" {
		return finalData, "Election Details Updated Successfully!", errors.New("Election has been completed already")
	}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return finalData, "Error Occurred", err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	UpdateQuery := bson.D{}
	if reqData.Location != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "location", Value: reqData.Location})
	}
	if reqData.ElectionDate != "" {
		electionDate, err := ConvertDate(reqData.ElectionDate)
		if err != nil {
			return finalData, "Error occurred", err
		}
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "election_date", Value: electionDate})
	}
	if reqData.ResultDate != "" {
		resultDate, err := ConvertDate(reqData.ResultDate)
		if err != nil {
			return finalData, "Error occurred", err
		}
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "result_date", Value: resultDate})
	}
	if reqData.ElectionStatus != "" {
		var resultDateStr string
		if reqData.ResultDate == "" {
			resultDateStr = time.Now().Format("2006-01-02")
		} else {
			resultDateStr = reqData.ResultDate
		}
		resultDate, err := ConvertDate(resultDateStr)
		if err != nil {
			return finalData, "Error occurred", err
		}
		if resultDate != data[0].ResultDate {
			return finalData, "Error Occurred", errors.New("Election Status can't be updated")
		}
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "election_status", Value: reqData.ElectionStatus})
	}
	if reqData.Result != "" {
		UpdateQuery = append(UpdateQuery, primitive.E{Key: "result", Value: reqData.Result})
	}
	query := bson.D{{"$set", UpdateQuery}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	er := CollectionElectionDetails.FindOneAndUpdate(ctx, filter, query, opts).Decode(&finalData)
	if er != nil {
		return finalData, "Error occurred", errors.New("Error while updating document")
	}

	if finalData == nil {
		return finalData, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}
	return finalData, "Election Details Updated Successfully!", nil
}

func (e *Connection) DeactivateElection(idStr string) (bson.M, string, error) {
	var updatedDocument bson.M
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return updatedDocument, "Error Occurred", err
	}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"_id", id}},
				bson.D{{"election_status", bson.M{"$ne": "Deactivated"}}},
			},
		},
	}
	update := bson.D{{"$set", bson.D{primitive.E{Key: "election_status", Value: "Deactivated"}}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	r := CollectionElectionDetails.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if r != nil {
		return updatedDocument, "Error Occurred", r
	}

	if updatedDocument == nil {
		return updatedDocument, "Error Occurred", errors.New("Data not present in db given by Id or it is deactivated")
	}

	return updatedDocument, "Election details deactivated successfully!", nil
}

func SetValueInElectionModelStruct(data model.ElectionRequest) (model.ElectionDetails, error) {
	var electionData model.ElectionDetails
	electionDate, err := ConvertDate(data.ElectionDate)
	if err != nil {
		return electionData, err
	}
	resultDate, err := ConvertDate(data.ResultDate)
	if err != nil {
		return electionData, err
	}
	electionData.ElectionDate = electionDate
	electionData.ResultDate = resultDate
	electionData.Location = data.Location
	electionData.ElectionStatus = data.ElectionStatus
	return electionData, err
}

func ConvertDbResultIntoElectionStruct(fetchDataCursor *mongo.Cursor) ([]*model.ElectionDetails, error) {
	var finaldata []*model.ElectionDetails
	for fetchDataCursor.Next(ctx) {
		var data model.ElectionDetails
		err := fetchDataCursor.Decode(&data)
		if err != nil {
			return finaldata, err
		}
		finaldata = append(finaldata, &data)
	}
	return finaldata, nil
}

// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
// ======================================Token=============================================
func (e *Connection) GenerateToken(request model.LoginDetails) (string, string, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"mail_id", request.MailId}},
				bson.D{{"password", request.Password}},
			},
		},
	}

	// check if email exists and password is correct
	record, err := CollectionUserDetails.Find(ctx, filter)
	if err != nil {
		return "", "Error", err
	}
	convData, err := ConvertDbResultIntoUserStruct(record)
	if err != nil {
		return "", "Error", err
	}
	if len(convData) != 0 {
		tokenString, err := auth.GenerateJWT(request.MailId)
		if err != nil {
			return "", "Error", err
		}
		return tokenString, "Token Generated Successfully", err
	} else {
		return "", "Unable to Generate Token", errors.New("Invalid Credentials")
	}
}

//XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXS
//===========================================================================================

func (e *Connection) ElectionResultById(idStr string) ([]model.Candidates, string, error) {
	var updatedDocument []*model.ElectionDetails
	var candidate []model.Candidates
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return candidate, "Error Occurred", err
	}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"_id", id}},
				bson.D{{"election_status", bson.M{"$eq": "Completed"}}},
			},
		},
	}
	data, err := CollectionElectionDetails.Find(ctx, filter)
	if err != nil {
		return candidate, "Error Occurred", err
	}
	updatedDocument, err = ConvertDbResultIntoElectionStruct(data)
	if updatedDocument == nil {
		return candidate, "Result fetch successfully", errors.New("Given id is not valid or the result date is not today")
	}
	candidates := updatedDocument[0].Candidates
	for i := range candidates {
		if candidates[i].NominationStatus == "Verified" {
			candidate = append(candidate, candidates[i])
		}
	}
	return candidate, "Result fetch successfully", nil
}

func (e *Connection) CaseVote(request model.CastVoteReq, mailId string) (string, string, error) {
	candidateId, err := primitive.ObjectIDFromHex(request.CandidateId)
	if err != nil {
		return "", "Error Occurred", err
	}
	electionId, err := primitive.ObjectIDFromHex(request.ElectionId)
	if err != nil {
		return "", "Error Occurred", err
	}
	c, err := CheckUserIsValidToVote(mailId)
	if err != nil {
		return "", "Error Occurred", err
	}
	if !c {
		return "", "Error Occurred", errors.New("Invalid User")
	}
	b, err := CheckIfUserAlreadyVoted(mailId, electionId)
	if err != nil {
		return "", "Error Occurred", err
	}
	if !b {
		return "", "Error Occurred", errors.New("Already Voted")
	}
	str, err := AddVoteInElectionDb(electionId, candidateId)
	fmt.Println(str)
	if err != nil {
		return "", "Error Occurred", err
	}
	str2, err := AddElectionIdInUserDb(electionId, mailId)
	fmt.Println(str2)
	if err != nil {
		return "", "Error Occurred", err
	}
	return "Voted successfully", "Voted successfully", nil
}

func CheckUserIsValidToVote(mailId string) (bool, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"mail_id", mailId}},
				bson.D{{"is_verified", true}},
			},
		},
	}
	data, err := CollectionUserDetails.Find(ctx, filter)
	if err != nil {
		return false, err
	}
	convUser, err := ConvertDbResultIntoUserStruct(data)
	if err != nil {
		return false, err
	}
	if len(convUser) == 0 {
		return true, err
	}
	return false, err
}

func CheckIfUserAlreadyVoted(mailId string, electionId primitive.ObjectID) (bool, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"mail_id", mailId}},
				bson.D{{"voted.election_id", electionId}},
			},
		},
	}
	data, err := CollectionUserDetails.Find(ctx, filter)
	if err != nil {
		return false, err
	}
	convUser, err := ConvertDbResultIntoUserStruct(data)
	if err != nil {
		return false, err
	}
	if len(convUser) == 0 {
		return true, err
	}
	return false, err
}

func AddVoteInElectionDb(electionId, candidateId primitive.ObjectID) (string, error) {
	var doc bson.M

	filter := bson.D{primitive.E{Key: "_id", Value: electionId}}
	cur, err := CollectionElectionDetails.Find(ctx, filter)
	if err != nil {
		return "Error Occurred", err
	}
	convElection, err := ConvertDbResultIntoElectionStruct(cur)
	if err != nil {
		return "Error Occurred", err
	}
	Candidates := convElection[0].Candidates
	for i := range Candidates {
		if Candidates[i].UserId == candidateId {
			if Candidates[i].NominationStatus == "Verified" {
				Candidates[i].Votecount = Candidates[i].Votecount + 1
			} else {
				return "Error Occurred", errors.New("Invalid Candidate")
			}
		}
	}
	update := bson.D{primitive.E{Key: "$set", Value: convElection[0]}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	CollectionElectionDetails.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if doc == nil {
		return "Error Occurred While Voting", errors.New("Unable to vote")
	}
	return "Voted Successfully", nil
}

func AddElectionIdInUserDb(electionId primitive.ObjectID, mailId string) (string, error) {
	var doc bson.M
	filter := bson.D{primitive.E{Key: "mail_id", Value: mailId}}
	updateQuery := bson.D{primitive.E{Key: "election_id", Value: electionId}}
	update := bson.D{primitive.E{Key: "voted", Value: updateQuery}}
	query := bson.D{primitive.E{Key: "$push", Value: update}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	CollectionUserDetails.FindOneAndUpdate(ctx, filter, query, opts).Decode(&doc)
	if doc == nil {
		return "Error Occurred While adding ElectionId", errors.New("Unable to add electionId")
	}
	return "Added ElectionId Successfully", nil
}
