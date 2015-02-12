package portals

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Agent struct {
	CodeName 	string `json:"codeName"`
	RealName 	string `json:"realName"`
	Email 		string `json:"email"`
}

func SaveAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringkey := vars["key"]
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var agent Agent
	err := json.Unmarshal(body, &agent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Agent", stringkey, 0, nil)
	key, err = datastore.Put(c, key, &agent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func GetAgents(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Agent").Limit(10)
	agents := make([]Agent, 0, 10)
	if _, err := q.GetAll(c, &agents); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(&agents)
	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, string(b))
}

func GetAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringkey := vars["key"]
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Agent", stringkey, 0, nil)
	var agent Agent
	if err := datastore.Get(c, key, &agent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(&agent)
	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, "[%s]", string(b))
}
