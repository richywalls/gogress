package portals

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/plus/v1"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

//TODO: on error do something... :S
var hmacTestKey []byte

func init() {
	var err error
	absPath, _ := filepath.Abs("hmactest")
	hmacTestKey, err = ioutil.ReadFile(absPath)
	if err != nil {
		panic("No puedo leer el archivo hmactest")
	}
}

type SattelizerData struct {
	ClientId    string `json:"clientId"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirectUri"`
}

var config = &oauth2.Config{
	ClientID:     "164620448986-olal315lm7t73p7qgp47isa5jl31le8r.apps.googleusercontent.com", // from https://code.google.com/apis/console/
	ClientSecret: "LHfl8sX9PzLf4bxVLjLhSTP9",                                                 // from https://code.google.com/apis/console/
	Endpoint:     google.Endpoint,
	Scopes:       []string{plus.PlusLoginScope},
}

type UserData struct {
	Favourites []string `json:"favourites"`
	SysConfig  string   `json:"sysConfig"`
}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if userId := GetUserId(r); userId != "" {
		var userData UserData
		datastore.Get(c, datastore.NewKey(c, "UserData", userId, 0, datastore.NewKey(c, "User", userId, 0, nil)), &userData)
		b, _ := json.Marshal(userData)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, string(b))
		return
	}
	http.Error(w, "no user id present", http.StatusForbidden)
}
func SaveUserData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if userId := GetUserId(r); userId != "" {
		var userData UserData
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		log.Infof(c, "%s", body)
		err := json.Unmarshal(body, &userData)
		log.Infof(c, "%s", userData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		datastore.Put(c, datastore.NewKey(c, "UserData", userId, 0, datastore.NewKey(c, "User", userId, 0, nil)), &userData)
		return
	}
	http.Error(w, "no user id present", http.StatusForbidden)
}
func GetUserId(r *http.Request) string {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return hmacTestKey, nil
	})
	if err != nil {
		return ""
	}
	return token.Claims["id"].(string)
}
func GetUserAgentId(r *http.Request) int64 {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return hmacTestKey, nil
	})
	if err != nil {
		return 0
	}
	return int64(token.Claims["agentId"].(float64))
}
func Authenticate(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var sattData SattelizerData
	err := json.Unmarshal(body, &sattData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	config.RedirectURL = sattData.RedirectUri
	token, err := config.Exchange(c, sattData.Code)
	client := config.Client(c, token)
	plusService, err := plus.New(client)
	person, err := plusService.People.Get("me").Do()

	//Create user in DB if not exists
	var user User
	key := datastore.NewKey(c, "User", person.Id, 0, nil)
	if err := datastore.Get(c, key, &user); err != nil {
		//User doesn't exist, create
		user.Name = person.DisplayName
		user.Allowed = false
		if person.Emails != nil && person.Emails[0].Type == "account" {
			user.Email = person.Emails[0].Value
		}
		user.Id = person.Id
		key, err = datastore.Put(c, key, &user)
	}
	tokenString, err := createJWT(person.Id, person.DisplayName, person.Image.Url, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "{\"token\":\"%s\"}", tokenString)
}

func createJWT(id string, name string, imageUrl string, user User) (string, error) {
	jottoken := jwt.New(jwt.GetSigningMethod("HS256"))
	jottoken.Claims["id"] = id
	jottoken.Claims["name"] = name
	jottoken.Claims["imageUrl"] = imageUrl
	jottoken.Claims["agentId"] = user.AgentId
	jottoken.Claims["allowed"] = user.Allowed
	jottoken.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := jottoken.SignedString(hmacTestKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func AuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			return hmacTestKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if !token.Claims["allowed"].(bool) {
			http.Error(w, "user not allowed or smurf", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
