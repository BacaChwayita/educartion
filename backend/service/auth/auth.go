package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/P-SEN371-Group-3/educartion/config"
	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	// 12 chosen as per AI suggestion:
	// Cost 12 is a good default (adjust based on performance)
	// Note that password is salted here by the bcrypt library, so no need to handle salting separately
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func checkPassword(password, storedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

// Register
// Adds a user to the Account table in the db.
// Will hash the password and store the hashed value
// @param cfg config of the application, used for db connection
// @param rr RegisterRequest struct - request of the endpoint
// @return error any errors
func Register(cfg *config.Config, rr model.RegisterRequest) error {

	hashedPassword, err := hashPassword(rr.Password_text)
	if err != nil {
		return errors.New("Failed to hash password: " + err.Error())
	}

	newAccount := model.Account{
		Full_name:     rr.Full_name,
		Email:         rr.Email,
		Password_hash: hashedPassword,
		Password_salt: "", // TODO: Remove Password_salt, since bcrypt does salting itself
		Role:          "customer",
		Created_at:    time.Now(),
	}
	fmt.Printf("hashed %s to %s:", rr.Password_text, hashedPassword) // TODO: Remove

	_, err = db.InsertAccount(cfg, newAccount)
	if err != nil {
		return errors.New("Failed to Insert Account: " + err.Error())
	}

	return nil
}
