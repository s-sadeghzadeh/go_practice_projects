package controllers

import (
	"Crud_RestApi_With_Mux_GORM/database"
	"Crud_RestApi_With_Mux_GORM/entities"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	//"gorm.io/gorm"
	"github.com/gorilla/mux"
)

// ///////////////////////////////////////////////////////////////////////////
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome to home page")
}

// ///////////////////////////////////////////////////////////////////////////////
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}

// //////////////////////////////////////////////////
func SignUp(w http.ResponseWriter, r *http.Request) {

	//connection := database.Connect(AppConfig.ConnectionString)
	//defer CloseDatabase(connection)

	//defer database.CloseDatabase()

	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err entities.Error
		err = SetError(err, "Error in reading payload.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var dbuser entities.User
	database.Instance.Where("email = ?", user.Email).First(&dbuser)

	//check email is alredy registered or not
	if dbuser.Email != "" {
		var err entities.Error
		err = SetError(err, "Email already in use")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password, err = GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("Error in password hashing.")
	}

	//insert user details in database
	//connection.Create(&user)
	database.Instance.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ////////////////////////////////////////////////////////
func SignIn(w http.ResponseWriter, r *http.Request) {
	//connection := client.Instance()
	//defer database.CloseDatabase()

	var authDetails entities.Authentication

	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		var err entities.Error
		err = SetError(err, "Error in reading payload.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authUser entities.User
	database.Instance.Where("email = 	?", authDetails.Email).First(&authUser)

	if authUser.Email == "" {
		var err entities.Error
		err = SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	check := CheckPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		var err entities.Error
		err = SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		var err entities.Error
		err = SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var token entities.Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// //////////////////////////////////////////////////////
func AdminIndex(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}

// ////////////////////////////////////////////////////////
func UserIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling UserIndex")

	if r.Header.Get("Role") != "user" {
		w.Write([]byte("Not Authorized."))
		return
	}
	w.Write([]byte("Welcome, User."))
}

// /////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
func AddContact(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var contact entities.Contact
	json.NewDecoder(r.Body).Decode(&contact)
	database.Instance.Create(&contact)
	json.NewEncoder(w).Encode(contact)
}

// ////////////////////////////////////////////////////////////////////////////////
func GetContactByIDs(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}

	contactId := mux.Vars(r)["id"]

	contactIDsArray := strings.Split(contactId, ",")

	var intSlice []int
	for _, str := range contactIDsArray {
		intNum, _ := strconv.Atoi(str)
		intSlice = append(intSlice, intNum)
	}

	if checkIfEntitiesExists(intSlice) == false {
		json.NewEncoder(w).Encode("Contact Not Found!")
		return
	}
	var contacts []entities.Contact
	database.Instance.Preload("CompanyInfo").Find(&contacts, intSlice)
	// database.Instance.Where(&contacts)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)

}

// ////////////////////////////////////////////////////////////////////////////////
func GetContacts(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}

	var contacts []entities.Contact
	// database.Instance.Find(&contacts)
	database.Instance.Preload("CompanyInfo").Find(&contacts)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contacts)
}

// ////////////////////////////////////////////////////////////////////////////////
func UpdateContact(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}

	contactId := mux.Vars(r)["id"]
	if checkIfContactExists(contactId) == false {
		json.NewEncoder(w).Encode("Contact Not Found!")
		return
	}
	var contact entities.Contact
	database.Instance.First(&contact, contactId)
	json.NewDecoder(r.Body).Decode(&contact)
	database.Instance.Save(&contact)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
}

// ////////////////////////////////////////////////////////////////////////////////
func DeleteContact(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	contactId := mux.Vars(r)["id"]
	contactIDsArray := strings.Split(contactId, ",")
	// if checkIfEntitiesExists(contactId) == false {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	json.NewEncoder(w).Encode("Contact Not Found!")
	// 	return
	// }

	var intSlice []int
	for _, str := range contactIDsArray {
		intNum, _ := strconv.Atoi(str)
		intSlice = append(intSlice, intNum)
	}

	if checkIfEntitiesExists(intSlice) == false {
		json.NewEncoder(w).Encode("Contact Not Found!")
		return
	}

	var contact entities.Contact
	database.Instance.Unscoped().Delete(&contact, intSlice)
	json.NewEncoder(w).Encode("Contact Deleted Successfully!")
}

// ////////////////////////////////////////////////////////////////////////////////
// ///// باید اصلاح شود
func checkIfContactExists(contactId string) bool {
	var contact entities.Contact
	database.Instance.First(&contact, contactId)
	if contact.ID == 0 {
		return false
	}
	return true
}

func checkIfEntitiesExists(contactIDsArray []int) bool {
	var contact []entities.Contact
	result := database.Instance.Find(&contact, contactIDsArray)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

// ///////////////////////////////////////////////////////////////
func AddCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var company entities.Company
	json.NewDecoder(r.Body).Decode(&company)
	database.Instance.Create(&company)
	json.NewEncoder(w).Encode(company)
}

//////////////////////////////////////////////
