package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	svix "github.com/svix/svix-webhooks/go"
)

const (
	defaultBind = ":8080"

	healthCheckPath = "/_healthcheck"

	svixServer = "https://api.eu.svix.com"
	svixToken  = "testsk_CSkhagouqu-JXgZznr35dG2TYTmsCPnb"
)

var Version = "v0.0"

var svixAppID = "app_2BEv2hBcE2ICiB6hq1QOVTVBWgF"

func main() {
	fmt.Println("version:", Version)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	router := mux.NewRouter()
	router.Handle(healthCheckPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("health check OK")
			w.WriteHeader(http.StatusOK)
		}),
	)

	if err := godotenv.Load(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	serverUrl, _ := url.Parse(svixServer)
	svixClient := svix.New(svixToken, &svix.SvixOptions{
		ServerUrl: serverUrl,
	})
	spew.Dump(svixClient)

	app, err := svixClient.Application.GetOrCreate(&svix.ApplicationIn{
		Name: "Formance Webhooks Application",
		Uid:  &svixAppID,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	spew.Dump(app)

	logger.Infof("starting http server on address: %s", defaultBind)
	if err := http.ListenAndServe(defaultBind, router); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
