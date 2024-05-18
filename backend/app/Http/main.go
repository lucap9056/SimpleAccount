package Server

import (
	"context"
	"fmt"
	"net/http"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	"simple_account/app/Http/Author"
	DELETE "simple_account/app/Http/Delete"
	GET "simple_account/app/Http/Get"
	"simple_account/app/Http/Message"
	POST "simple_account/app/Http/Post"
	PUT "simple_account/app/Http/Put"
	"simple_account/app/Logger"
)

type API struct {
	Server http.Server
	DB     *Database.API
	Email  *Email.Manager
	Auth   *Auths.Auth
	Logs   *Logger.Manager
}

func New(port int, databaseConfig Database.Config, emailConfig Email.Config, auth *Auths.Auth, logs *Logger.Manager) (*API, error) {
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

	email, err := Email.New(emailConfig)
	if err != nil {
		return nil, err
	}
	api.Email = email

	api.Auth = auth

	api.Logs = logs

	return &api, nil
}

func (api *API) MainHandler(writer http.ResponseWriter, req *http.Request) {

	author := Author.Get(api.Auth, writer, req)
	ctx := &Message.Context{
		Author:   author,
		Database: api.DB,
		Auth:     api.Auth,
		Writer:   writer,
		Request:  req,
		Email:    api.Email,
		Logs:     api.Logs,
	}

	switch req.Method {
	case http.MethodGet:
		GET.Handler(ctx)
	case http.MethodPost:
		POST.Handler(ctx)
	case http.MethodPut:
		PUT.Handler(ctx)
	case http.MethodDelete:
		DELETE.Handler(ctx)
	}
}

func (api *API) Start() {
	err := api.Server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (api *API) Exit(ctx context.Context) error {
	return api.Server.Shutdown(ctx)
}