package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/macolby42/fender-platform-challenge/dataModels"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

		var user dataModels.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			fmt.Printf("%d - bad request\n", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "no sign up data provided"}`))
		}

		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/token?grant_type=client_credentials&client_id=%s&client_secret=%s&scope=read", user.Email, user.Password))
		if err != nil {
			fmt.Printf("%d - unauthorized\n", http.StatusUnauthorized)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "unauthorized"}`))
		}

		reader = resp.Body
		response, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Printf("%d - bad request\n", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "no sign up data provided"}`))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
