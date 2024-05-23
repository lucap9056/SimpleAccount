package Extension

import (
	"context"
	"fmt"
	"net/http"
	"simple_account/app/Auths"
	"simple_account/app/Database"
	"simple_account/app/Email"
	ExtensionGet "simple_account/app/ExtensionChannel/Get"
	ExtensionPost "simple_account/app/ExtensionChannel/Post"
	"simple_account/app/Http/Message"
	"simple_account/app/Http/Url"
	"simple_account/app/Logger"
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

	mux.HandleFunc("/", api.MainHandler)

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
	url := Url.New(req.URL)

	ctx := &Message.Context{
		Database: api.DB,
		Auth:     api.Auth,
		Writer:   writer,
		Request:  req,
		Logs:     api.Logs,
		Url:      &url,
	}

	switch ctx.Request.Method {
	case http.MethodGet:
		ExtensionGet.Handler(ctx)
	case http.MethodPost:
		ExtensionPost.Handler(ctx)
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
