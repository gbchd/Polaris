package token

import (
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/guillaumebchd/polaris/pkg/authentication"
	"github.com/pkg/errors"
)

const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub

	LifetimeIdToken      = 24 * time.Hour
	LifetimeAccessToken  = 24 * time.Hour
	LifetimeRefreshToken = 7 * 24 * time.Hour
)

var (
	// We don't store the private key here, just in case
	publicKey []byte
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// read the key files before starting http handlers
func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		log.Fatal(err)
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func ServePubKeyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write(publicKey)
}

func CreateAccessToken(client_id string, user_mail string) (string, error) {
	Claims := jwt.MapClaims{}
	Claims["iss"] = "Polaris"
	Claims["aud"] = client_id
	Claims["email"] = user_mail
	Claims["type"] = "access"
	Claims["exp"] = time.Now().Unix() + int64(LifetimeAccessToken)
	Claims["iat"] = time.Now().Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims)

	token, err := t.SignedString(signKey)
	return token, err
}

func CreateRefreshToken(client_id string, user_mail string, access_token string) (string, error) {
	h := sha256.New()
	h.Write([]byte(access_token))

	Claims := jwt.MapClaims{}
	Claims["iss"] = "Polaris"
	Claims["aud"] = client_id
	Claims["email"] = user_mail
	Claims["type"] = "refresh"
	Claims["at_hash"] = h.Sum(nil)
	Claims["exp"] = time.Now().Unix() + int64(LifetimeAccessToken)
	Claims["iat"] = time.Now().Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims)

	token, err := t.SignedString(signKey)
	return token, err
}

func CreateIdToken(client_id string, user authentication.User) (string, error) {
	Claims := jwt.MapClaims{}
	Claims["iss"] = "Polaris"
	Claims["aud"] = client_id
	Claims["name"] = user.Name
	Claims["email"] = user.Email
	Claims["type"] = "id"
	Claims["exp"] = time.Now().Unix() + int64(LifetimeAccessToken)
	Claims["iat"] = time.Now().Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims)

	token, err := t.SignedString(signKey)
	return token, err
}

func verifyToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func CheckToken(t string) (*jwt.Token, bool) {
	token, err := verifyToken(t)
	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	return token, ok && token.Valid && claims.VerifyExpiresAt(time.Now().Unix(), true)
}
