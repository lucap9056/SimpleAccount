package simple_account_http

import (
	"context"
	"fmt"
	"net/http"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	DELETE "simple_account/app/Http/Delete"
	GET "simple_account/app/Http/Get"
	POST "simple_account/app/Http/Post"
	PUT "simple_account/app/Http/Put"
	"time"
)

type API struct {
	Server      http.Server
	DB          *Database.API
	EmailSender *Email.Sender
	Auth        *Auths.Auth
}

func New(port int, databaseConfig Database.Config, emailConfig Email.Config, auth *Auths.Auth) (*API, error) {
	api := API{}
	mux := http.NewServeMux()

	mux.HandleFunc("/", api.MainHandler)

	api.Server = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	DB, err := Database.New(&databaseConfig)
	if err != nil {
		return nil, err
	}
	api.DB = DB

	Em, err := Email.New(emailConfig)
	if err != nil {
		return nil, err
	}
	api.EmailSender = Em

	api.Auth = auth

	return &api, nil
}

func (api *API) MainHandler(writer http.ResponseWriter, req *http.Request) {

	author := GetAuthor(api.Auth, writer, req)
	switch req.Method {
	case http.MethodGet:
		GET.Handler(author, api.Auth, api.DB, writer, req)
	case http.MethodPost:
		POST.Handler(api.Auth, api.DB, writer, req, api.EmailSender)
	case http.MethodPut:
		PUT.Handler(author, api.DB, writer, req)
	case http.MethodDelete:
		DELETE.Handler(author, api.DB, writer, req)
	}
}

func (api *API) Start() {
	err := api.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (api *API) Exit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return api.Server.Shutdown(ctx)
}
