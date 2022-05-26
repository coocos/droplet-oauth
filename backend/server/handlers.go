package server

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/digitalocean/godo"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const TemplatePath = "templates"

var dropletTemplate = template.Must(template.ParseFiles(filepath.Join(TemplatePath, "layout.html"), filepath.Join("templates", "droplets.html")))
var loginTemplate = template.Must(template.ParseFiles(filepath.Join(TemplatePath, "layout.html"), filepath.Join("templates", "login.html")))

type Server struct {
	OAuth    oauth2.Config
	Sessions *sessions.FilesystemStore
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.Sessions.Get(r, "dashboard")
	if err != nil {
		log.Println("Failed to read session:", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if session.IsNew {
		session.Options = &sessions.Options{
			HttpOnly: true,
			MaxAge:   0,
		}
		session.Values["state"] = string(uuid.New().String())
		err := session.Save(r, w)
		if err != nil {
			log.Println("Failed to save session:", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
	if err = loginTemplate.ExecuteTemplate(w, "layout", struct {
		LoginUrl string
	}{
		s.OAuth.AuthCodeURL(session.Values["state"].(string), oauth2.AccessTypeOnline),
	}); err != nil {
		log.Println("Failed to render template:", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func (s *Server) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state string
	state := r.URL.Query().Get("state")
	session, _ := s.Sessions.Get(r, "dashboard")
	if state == "" || state != session.Values["state"] {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Grab access token and store it in the session
	token, err := s.OAuth.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		log.Println("Failed exchange code for token:", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session.Values["token"] = token.AccessToken
	err = session.Save(r, w)
	if err != nil {
		log.Println("Failed to store token:", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/droplets", http.StatusTemporaryRedirect)
}

func (s *Server) DropletHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := s.Sessions.Get(r, "dashboard")
	if session.Values["token"] == nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	client := godo.NewFromToken(session.Values["token"].(string))
	droplets, _, err := client.Droplets.List(context.Background(), nil)
	if err != nil {
		http.Error(w, "Failed to list droplets", http.StatusInternalServerError)
		return
	}

	if err = dropletTemplate.ExecuteTemplate(w, "layout", struct {
		Droplets []godo.Droplet
	}{droplets}); err != nil {
		log.Println("Failed to render template:", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
