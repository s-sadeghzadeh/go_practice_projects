package entities

import (
	"gorm.io/gorm"
)



//////////////////////////////////////////////////

type Contact struct {
	gorm.Model
	//ID           string  `json:"id"`
	Name         string  `json:"name"`
	Family       string  `json:"family"`
	MobileNumber string  `json:"mobileNumber"`
	PhoneNumber  string  `json:"phoneNumber"`
	Address      string  `json:"address"`
	CompanyInfo  Company `json:"companyInfo" gorm:"foreignKey:CompanyID"`
	CompanyID    int     `json:"companyID"`
}

type Company struct {
	gorm.Model
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

//////////////////////////////////////////////////

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}


