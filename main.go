package main

import (
	"fmt"
	"net/http"

	"io/ioutil"
	"log"
	"time"

	"github.com/goinggo/tracelog"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	googleOauth "golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	// You must register the app at https://console.developers.google.com/apis/credentials
	// Set callback to http://127.0.0.1:4567/GTJoauth2callback
	// Set ClientId and ClientSecret to
	oauthConf = &oauth2.Config{
	/*		ClientID:     "965925130911-lb3hu48uj4u5ab1l8acs7e3u40fqlve4.apps.googleusercontent.com",
			ClientSecret: "6JctRfqfnGeOHPiDBluaRd-o",
			// select level of access you want https://developer.github.com/v3/oauth/#scopes
			Scopes:      []string{gmail.MailGoogleComScope},
			Endpoint:    googleOauth.Endpoint,
			RedirectURL: "http://localhost:4567/oauth2",*/
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "thisshouldberandom"
	htmlIndex        []byte //= `<html><body><h1>Logged in with <a href="/login">Gmail</a><h1></body></html>`
	errorf           = fmt.Errorf("Exception At...")
)

type message struct {
	size     int64
	gmailID  string
	date     string // retrieved from message header
	snippet  string
	subject  string
	dateTime time.Time
}

type timeSlice []message

// Respond on initial page request
func mainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlIndex))

		tracelog.Trace("gmail", "gmailMain", "%s %s took %s", r.Method, r.URL, time.Since(began))
	})
}

//After success page handler
func finishedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><h1>Successfully finished<h1></body></html>`))

		tracelog.Trace("gmail", "gmailMain", "%s %s took %s", r.Method, r.URL, time.Since(began))
	})
}

// login
func gmailLoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		oauthConf.RedirectURL = "http://localhost:4567/GTJoauth2callback"
		oauthStateString = generateUUID().String()
		url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		fmt.Println("URL will be = ", url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		tracelog.Trace("gmail", "gmailMain", "%s %s took %s", r.Method, r.URL, time.Since(began))
	})
}

// Gmail oauth callback. Called by gmail after authorization is granted
func gmailCallbackHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()

		state := r.FormValue("state")
		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		fmt.Println("Recevied URL callback = ", r.RequestURI)

		code := r.FormValue("code")
		fmt.Println("Recevied code = ", code)

		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		client := oauthConf.Client(context.Background(), token)
		gmailService, err := gmail.New(client)

		if err != nil {
			tracelog.Errorf(errorf,
				"gmail", "gmailMain", "Unable to create Gmail service: %v", err)
		} else {
			fmt.Println("Service established = ", gmailService)
		}

		processMessagesInInbox(gmailService)

		http.Redirect(w, r, "/finished", http.StatusTemporaryRedirect)

		tracelog.Trace("gmail", "gmailMain", "%s %s took %s", r.Method, r.URL, time.Since(began))
	})
}

func (messages timeSlice) Len() int {
	return len(messages)
}

func (messages timeSlice) Less(i, j int) bool {
	return messages[i].dateTime.Before(messages[j].dateTime)
}

func (messages timeSlice) Swap(i, j int) {
	messages[i], messages[j] = messages[j], messages[i]
}

func main() {

	tracelog.StartFile(tracelog.LevelTrace, "./logs", 1)

	clientSecret, err := ioutil.ReadFile("./keys/client_secret_web.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	htmlIndex, err = ioutil.ReadFile("./index.html")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json

	oauthConf, err = googleOauth.ConfigFromJSON(clientSecret, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	fmt.Println("oauthConf = ", oauthConf)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", mainHandler())
	http.Handle("/login", gmailLoginHandler())
	http.Handle("/GTJoauth2callback", gmailCallbackHandler())
	http.Handle("/finished", finishedHandler())

	fmt.Println("Started running on http://127.0.0.1:4567")

	fmt.Println(http.ListenAndServe(":4567", nil))

}
