package auth

import (
	"errors"
	"time"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func getJWTKey() []byte {
	key := "my_secret_key_123" // TODO: Move to secret place and get proper key
	return []byte(key)
}

func createJWT() (string, error) {
	key := getJWTKey()

	// Registered Claims: Pre-defined, standard claims (IANA) for ensuring interoperability.
	// - iss (Issuer): Identifies who issued the token.
	// - sub (Subject): Identifies the entity the token represents.
	// - aud (Audience): The recipient for whom the token is intended.
	// - exp (Expiration Time): Time after which the token is invalid.
	// - nbf (Not Before): Time before which the token must not be accepted.
	// - iat (Issued At): Time the token was created.
	// - jti (JWT ID): Unique identifier for the token.

	nbfTime := time.Now()
	expTime := nbfTime.AddDate(0, 1, 0) // Add one month

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New(),
		"nbf": nbfTime.Unix(),
		"exp": expTime.Unix(),
	})

	tokenString, err := token.SignedString(key)

	return tokenString, err
}

func getAccountIdFromJWT(tokenString string) (model.Account, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return getJWTKey(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return model.Account{}, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	acc, err := db.GetAccountByJTI(claims["jti"].(string))
	if err != nil {
		return model.Account{}, err
	}

	return acc, err

}

// Register
// Adds a user to the Account table in the db.
// Will hash the password and store the hashed value
//
// @param rr RegisterRequest struct - request of the endpoint
// @return error any errors
func Register(rr model.RegisterRequest) error {
	hashedPassword, err := hashPassword(rr.Password_text)
	if err != nil {
		return ErrHashPassword
	}

	newAccount := model.Account{
		Full_name:      rr.Full_name,
		Email:          rr.Email,
		Password_hash:  hashedPassword,
		Role:           "customer",
		Login_attempts: 0,
		Is_active:      true,
		Created_at:     time.Now(),
	}

	_, err = db.InsertAccount(newAccount)
	if err != nil {
		return ErrDBInsert
	}

	return nil
}

// LoginWithEmail
// Logs a user in using email and password.
// If successful, a JWT token will be returned, as well as the user profile.
// If unsuccessful, the login_attempts is incremented.
// Once login_attemts reach it's max value, all login_attempts will fail.
func LoginWithEmail(lr model.LoginWithEmailRequest) (string, model.Account, error) {
	const max_login_attempts int = 3
	var acc model.Account

	acc, err := db.GetAccountByEmail(lr.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", model.Account{}, ErrAccountNotFound
		}

		return "", model.Account{}, ErrUnkown
	}

	// If account is inactive then need to reset account (forget password)
	if !acc.Is_active {
		return "", model.Account{}, ErrAccountInactive
	}

	ok := checkPassword(lr.Password_text, acc.Password_hash)
	if !ok {
		db.IncLoginAttempts(acc.Account_id)
		return "", model.Account{}, ErrPasswordIncorrect
	}

	tokenString, err := createJWT()
	if err == nil {
		accLogin := model.AccountLogin{
			Token_id:     -1,
			Token_string: tokenString,
			Account_id:   acc.Account_id,
		}
		accLogin, err := db.InsertAccountLogin(accLogin)
		if err != nil {
			return "", acc, ErrFailedToCreateJWT // Failed to make a JWT, but still managed to login. User can continue.
		}
	}

	return tokenString, acc, nil
}

func Logout(lr model.LogoutRequest) error {

	rowcount, err := db.DeleteAccountLoginByToken(lr.Token)
	if err != nil {
		return ErrUnkown
	}

	if rowcount == 0 {
		return ErrAccountLoginNotFound
	}

	if rowcount != 1 {
		return ErrUnkown
	}

	return nil

}
