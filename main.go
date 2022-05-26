package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/coocos/droplet-oauth/server"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func readPort() int {
	if port, ok := os.LookupEnv("SERVICE_PORT"); ok {
		portNumber, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal("Failed to parse port number:", err)
		}
		return portNumber
	}
	return 8000
}

func main() {

	sessionStore := sessions.NewFilesystemStore(os.TempDir(), []byte(os.Getenv("SESSION_KEY")))
	sessionStore.Options.MaxAge = 60
	sessionStore.Options.HttpOnly = true

	server := server.Server{
		Sessions: sessionStore,
		OAuth: oauth2.Config{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			RedirectURL:  os.Getenv("REDIRECT_URL"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://cloud.digitalocean.com/v1/oauth/authorize",
				TokenURL: "https://cloud.digitalocean.com/v1/oauth/token",
			},
			Scopes: []string{"read"},
		},
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/droplets", server.DropletHandler)
	http.HandleFunc("/login", server.LoginHandler)
	http.HandleFunc("/redirect", server.RedirectHandler)

	port := readPort()
	log.Println("Listening on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
