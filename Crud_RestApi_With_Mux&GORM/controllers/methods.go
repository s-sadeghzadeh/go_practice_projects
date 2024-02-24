package controllers

import (
	"Crud_RestApi_With_Mux_GORM/entities"

)


// set error message in Error struct
func SetError(err entities.Error, message string) entities.Error {
	err.IsError = true
	err.Message = message
	return err
}
