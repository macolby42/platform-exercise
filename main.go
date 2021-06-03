package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/macolby42/fender-platform-challenge/endpoints"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var authServer *server.Server
var clientStore *store.ClientStore
var tokenStore oauth2.TokenStore

func signup(w http.ResponseWriter, r *http.Request) {
	endpoints.Signup(w, r, clientStore)
}

func login(w http.ResponseWriter, r *http.Request) {
	endpoints.Login(w, r)
}

func logout(w http.ResponseWriter, r *http.Request) {
	endpoints.Logout(w, r, authServer, clientStore, tokenStore)
}

func update(w http.ResponseWriter, r *http.Request) {
	endpoints.Update(w, r, authServer, clientStore)
}

func delete(w http.ResponseWriter, r *http.Request) {
	endpoints.Delete(w, r, authServer, clientStore, tokenStore)
}

func main() {
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)

	manager := manage.NewDefaultManager()

	// client & token memory store
	clientStore = store.NewClientStore()
	var err error
	tokenStore, err = store.NewMemoryTokenStore()
	if err != nil {
		log.Fatal("could not create token store")
	}

	manager.MustTokenStorage(tokenStore, err)
	manager.MapClientStorage(clientStore)

	// set up auth server stuff
	authServer = server.NewDefaultServer(manager)
	authServer.SetAllowGetAccessRequest(true)
	authServer.SetClientInfoHandler(server.ClientFormHandler)

	authServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	authServer.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		authServer.HandleTokenRequest(w, r)
	})

	port := ":8080"
	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
