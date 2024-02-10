package controllers

import (
	"encoding/json"
	"net/http"
	"prj_crud/database"
	"prj_crud/entities"
	"strconv"
	"strings"
	"fmt"

	//"gorm.io/gorm"
	"github.com/gorilla/mux"
)


/////////////////////////////////////////////////////////////////////////////
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"welcome to home page")
}



// ////////////////////////////////////////////////////////////////////////////////
func AddContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contact entities.Contact
	json.NewDecoder(r.Body).Decode(&contact)
	database.Instance.Create(&contact)
	json.NewEncoder(w).Encode(contact)
}

// ////////////////////////////////////////////////////////////////////////////////
func GetContactByIDs(w http.ResponseWriter, r *http.Request) {
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
	var contacts []entities.Contact
	// database.Instance.Find(&contacts)
	database.Instance.Preload("CompanyInfo").Find(&contacts)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contacts)
}

// ////////////////////////////////////////////////////////////////////////////////
func UpdateContact(w http.ResponseWriter, r *http.Request) {

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
