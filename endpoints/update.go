package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/macolby42/fender-platform-challenge/dataModels"

	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func Update(w http.ResponseWriter, r *http.Request, authServer *server.Server, clientStore *store.ClientStore) {
	tokenInfo, err := authServer.ValidationBearerToken(r)
	if err != nil {
		fmt.Printf("%d - unauthorized\n", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "unauthorized"}`))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	ctx := r.Context()
	clientInfo, err := clientStore.GetByID(ctx, tokenInfo.GetClientID())
	if err != nil {
		fmt.Printf("%d - bad request\n", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "bad update value"}`))
	}

	reader := r.Body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("%d - bad request\n", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "bad update value"}`))
	}

	var toUpdate dataModels.User
	err = json.Unmarshal(body, &toUpdate)
	if err != nil {
		fmt.Printf("%d - bad request\n", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "bad update value"}`))
		return
	}

	if clientInfo.GetID() != toUpdate.Email {
		fmt.Printf("%d - bad request\n", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "invalid email for update"}`))
		return
	}

	clientStore.Set(clientInfo.GetID(), &models.Client{
		ID:     clientInfo.GetID(),
		Secret: toUpdate.Password,
		Domain: clientInfo.GetDomain(),
		UserID: toUpdate.Name,
	})

	str, _ := clientStore.GetByID(ctx, clientInfo.GetID())
	w.Write([]byte(str.GetUserID()))
}
