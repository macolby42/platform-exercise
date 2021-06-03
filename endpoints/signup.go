package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/macolby42/fender-platform-challenge/dataModels"
)

func Signup(w http.ResponseWriter, r *http.Request, clientStore *store.ClientStore) {
	if r.Method == "POST" {
		fmt.Printf("%s%s %s ", r.Host, r.RequestURI, r.Method)
		w.Header().Set("Content-Type", "application/json")

		if r.Body == nil {
			// TODO: reduce repeated code
			fmt.Printf("%d - bad request\n", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "no sign up data provided"}`))
		}

		reader := r.Body
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Printf("%d - bad request\n", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "no sign up data provided"}`))
		}

		var newUser dataModels.User
		err = json.Unmarshal(body, &newUser)
		if err != nil {
			fmt.Printf("%d - bad request\n", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "no sign up data provided"}`))
		}

		addClient(clientStore, newUser)

		w.WriteHeader(http.StatusAccepted)
	}
}

func addClient(clientStore *store.ClientStore, newUser dataModels.User) {
	clientStore.Set(newUser.Email, &models.Client{
		ID:     newUser.Email,
		Secret: newUser.Password,
		Domain: "http://localhost",
		UserID: newUser.Name,
	})
}
