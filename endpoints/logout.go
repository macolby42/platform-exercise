package endpoints

import (
	"net/http"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func Logout(w http.ResponseWriter, r *http.Request, authServer *server.Server, clientStore *store.ClientStore, tokenStore oauth2.TokenStore) {
	ctx := r.Context()
	at := strings.Split(r.Header.Get("Authorization"), " ")[1]
	err := tokenStore.RemoveByAccess(ctx, at)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("logged out!"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}
}
