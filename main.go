package main

import (
	"fmt"
	"net/http"

	"io/ioutil"
	"log"
	"time"

	"github.com/goinggo/tracelog"
	"golang.org/x/oauth2"
	googleOauth "golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"golang.org/x/net/context"
)

var (
	// You must register the app at https://github.com/settings/applications
	// Set callback to http://127.0.0.1:7000/github_oauth_cb
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

const htmlIndex = `<html><body><h1>Logged in with <a href="/login">Gmail</a><h1></body></html>`

// /
func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

// /login
func handleGmailLogin(w http.ResponseWriter, r *http.Request) {

	oauthConf.RedirectURL = "http://localhost:4567/GTJoauth2callback"

	oauthStateString = generateUUID().String()

	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)

	fmt.Println("URL will be = ", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// /Gmail oauth callback. Called by gmail after authorization is granted
func handleGmailCallback(w http.ResponseWriter, r *http.Request) {

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

	client := oauthConf.Client(context.Background(), token) //oauth2.NoContext

	fmt.Println("Client is ready = ", client)

	gmailService, err := gmail.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	if err != nil {
		tracelog.Errorf(fmt.Errorf("Exception At..."),
			"gmail", "gmailMain", "Unable to create Gmail service: %v", err)
	} else {
		fmt.Println("Service established = ", gmailService)

	}

	var total int64
	msgs := []message{}
	pageToken := ""

	//queryStr := fmt.Sprintf("from:%v newer_than:%v", "confirmation@mail.hotels.com", "50d")

	emails := getAllStoresEmails()

	queryStr := buildQueryString(emails)

	fmt.Println("Search string is = ", queryStr)

	for {
		req := gmailService.Users.Messages.List("me").Q(queryStr)
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			tracelog.Errorf(fmt.Errorf("Exception At..."), "gmail", "gmailMain",
				"Unable to retrieve messages: %v", err)
		}
		if r.Messages != nil {
			tracelog.Info("gmail", "gmailMain", "Processing %v messages...\n", len(r.Messages))
		} else {
			tracelog.Info("gmail", "gmailMain", "No messages to process")
		}
		for _, m := range r.Messages {
			msg, err := gmailService.Users.Messages.Get("me", m.Id).Do()
			if err != nil {
				tracelog.Errorf(fmt.Errorf("Exception At..."), "gmail",
					"gmailMain", "Unable to retrieve message %v: %v", m.Id, err)
			}
			total += msg.SizeEstimate
			date := ""
			subject := ""
			for _, h := range msg.Payload.Headers {
				if h.Name == "Date" {
					date = h.Value
				}
				if h.Name == "Subject" {
					subject = h.Value
				}
			}
			msgs = append(msgs, message{
				size:    msg.SizeEstimate,
				gmailID: msg.Id,
				date:    date,
				snippet: msg.Snippet,
				subject: subject,
			})
		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}
	tracelog.Info("gmail", "gmailMain", "total: %v\n", total)

	convertDateToDateTime(msgs)

	sortByDate(msgs)

	count := 0
	for _, m := range msgs {
		count++
		tracelog.Info("gmail", "gmailMain", "\nMessage URL: https://mail.google.com/mail/u/0/#all/%v\n" +
			"Subject: %q ,Size: %v, Date: %v, Snippet: %q\n",
			m.gmailID, m.subject, m.size, m.date, m.snippet)

		forwardToJetAnywhere(m)

		/*if err := svc.Users.Messages.Delete("me", m.gmailID).Do(); err != nil {*/
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

func forwardToJetAnywhere(message message) {

}

func main() {
	b, err := ioutil.ReadFile("./keys/client_secret_web.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json

	oauthConf, err = googleOauth.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	fmt.Println("oauthConf = ", oauthConf)

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGmailLogin)
	http.HandleFunc("/GTJoauth2callback", handleGmailCallback)

	fmt.Println("Started running on http://127.0.0.1:4567")

	fmt.Println(http.ListenAndServe(":4567", nil))

}
