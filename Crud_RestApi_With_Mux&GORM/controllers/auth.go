package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Crud_RestApi_With_Mux_GORM/entities"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	secretkey string = "secretkeyjwt"
)

// /////////////////////////////////////////////////////////////////////////////
// take password as input and generate new hash password from it
func GeneratehashPassword(password string) (string, error) {
	fmt.Println("calling GeneratehashPassword")

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// /////////////////////////////////////////////////////////////////////////////
// compare plain password with hash password
func CheckPasswordHash(password, hash string) bool {
	fmt.Println("calling CheckPasswordHash")

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// /////////////////////////////////////////////////////////////////////////////
// Generate JWT token
func GenerateJWT(email, role string) (string, error) {
	fmt.Println("calling GenerateJWT")

	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

///////////////////////////////////////////////////////////////////////////////
//---------------------MIDDLEWARE FUNCTION-----------------------

// check whether user is authorized or not
func IsAuthorized(roleArray []string, handler http.HandlerFunc) http.HandlerFunc {
	fmt.Println("calling IsAuthorized")

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err entities.Error
			err = SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte(secretkey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err entities.Error
			err = SetError(err, "Your Token has been expired.")
			json.NewEncoder(w).Encode(err)
			return
		}

		// آگر آرایه ورودی خالی باشد یعنی نیازی به بررسی رول کاربر نمی باشد و فقط اعبار توکن مهم است
		if len(roleArray) == 0 && token.Valid {
			fmt.Println("input role is empty")
			handler.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// if claims["role"] == "admin" {
			// 	r.Header.Set("Role", "admin")
			// 	handler.ServeHTTP(w, r)
			// 	return

			// } else if claims["role"] == "user" {
			// 	r.Header.Set("Role", "user")
			// 	handler.ServeHTTP(w, r)
			// 	return

			// }
			for _, role := range roleArray {
				if claims["role"] == role {
					fmt.Println("claims[role] == role")
					handler.ServeHTTP(w, r)
					return
				}
			}

		}


		fmt.Println("role no access")

		var reserr entities.Error
		reserr = SetError(reserr, "Not Authorized.")
		json.NewEncoder(w).Encode(reserr)



	}
}
