package endpoints

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func Delete(w http.ResponseWriter, r *http.Request, authServer *server.Server, clientStore *store.ClientStore, tokenStore oauth2.TokenStore) {
	if r.Method == "DELETE" {
		tokenInfo, err := authServer.ValidationBearerToken(r)
		if err != nil {
			fmt.Printf("%d - unauthorized\n", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "unauthorized"}`))
			return
		}

		ctx := r.Context()
		clientInfo, err := clientStore.GetByID(ctx, tokenInfo.GetClientID())
		if err != nil {
			fmt.Printf("%d - bad request\n", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "bad update value"}`))
		}

		clientStore.Set(clientInfo.GetID(), nil)
		at := strings.Split(r.Header.Get("Authorization"), " ")[1]
		tokenStore.RemoveByAccess(ctx, at)
		str, _ := clientStore.GetByID(ctx, clientInfo.GetID())
		if str == nil {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("deleted"))
		}
	} else {
		fmt.Printf("%d - method not allowed\n", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"message": "method not allowed"}`))
	}
}
