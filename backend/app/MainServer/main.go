package Server

import (
	"context"
	"fmt"
	"net/http"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	"simple_account/app/Http/Author"
	"simple_account/app/Http/Message"
	"simple_account/app/Http/Url"
	"simple_account/app/Logger"
	DELETE "simple_account/app/MainServer/Delete"
	GET "simple_account/app/MainServer/Get"
	POST "simple_account/app/MainServer/Post"
	PUT "simple_account/app/MainServer/Put"
)

type API struct {
	Server http.Server
	DB     *Database.API
	Email  *Email.Manager
	Auth   *Auths.Auth
	Logs   *Logger.Manager
}

func New(port int, db *Database.API, email *Email.Manager, auth *Auths.Auth, logs *Logger.Manager) (*API, error) {
	api := API{}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/", api.MainHandler)

	api.Server = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	api.DB = db

	api.Email = email

	api.Auth = auth

	api.Logs = logs

	return &api, nil
}

func (api *API) MainHandler(writer http.ResponseWriter, req *http.Request) {

	author := Author.Get(api.Auth, writer, req)
	url := Url.New(req.URL)
	url.Shift()

	ctx := &Message.Context{
		Author:   author,
		Database: api.DB,
		Auth:     api.Auth,
		Writer:   writer,
		Request:  req,
		Email:    api.Email,
		Logs:     api.Logs,
		Url:      &url,
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
