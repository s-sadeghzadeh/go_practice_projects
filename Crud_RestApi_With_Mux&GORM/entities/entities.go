package entities

import (
	"gorm.io/gorm"
)




type Contact struct {
	gorm.Model
	//ID           string  `json:"id"`
	Name         string  `json:"name"`
	Family       string  `json:"family"`
	MobileNumber string  `json:"mobileNumber"`
	PhoneNumber  string  `json:"phoneNumber"`
	Address      string  `json:"address"`
	CompanyInfo  Company  `json:"companyInfo" gorm:"foreignKey:CompanyID"`
	CompanyID	 int	  `json:"companyID"`

}

type Company struct {
	gorm.Model
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	
}







