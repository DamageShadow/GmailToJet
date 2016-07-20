package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"github.com/goinggo/tracelog"
	"sort"
	"time"
	"os/exec"
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

type Handler struct {
	code   string
	rState string
}

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		// select level of access you want https://developer.github.com/v3/oauth/#scopes
		Scopes:       []string{gmail.MailGoogleComScope},
		//Endpoint:     githuboauth.Endpoint,
	}

	ch = make(chan string)

	oauthStateString = "thisshouldberandom"

	code string
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(ctx, config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}


// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch = make(chan string)
	code := ""

	randState := fmt.Sprintf("st%d", time.Now().UnixNano())

	config.RedirectURL = "http://localhost:4567/oauth2callback"

	authURL := config.AuthCodeURL(randState)



	/*go openURL(authURL)

	//http.ListenAndServe(":4567", handler)

	handler := &Handler{code:"", rState:randState}

	server := httptest.NewServer(handler)

	server.URL = "http://localhost:4567/oauth2callback"

	code = handler.code*/


	http.HandleFunc("/oauth2callback", handleGoogleCallback)


}

// /github_oauth_cb. Called by github after authorization is granted
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	/*token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}*/

	fmt.Println("Received code " + code)


	/*	oauthClient := oauthConf.Client(oauth2.NoContext, token)
		client := github.NewClient(oauthClient)
		user, _, err := client.Users.Get("")
		if err != nil {
			fmt.Printf("client.Users.Get() faled with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)*/
}

/*func (handler *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rState := handler.rState

	if req.URL.Path == "/favicon.ico" {
		log.Printf("404")
		http.Error(rw, "", 404)
		return
	}
	if state := req.FormValue("state"); state != rState {
		log.Printf("State doesn't match: sent = %s, receive = %s, req = %#v", rState, state, req)
		http.Error(rw, "", 500)
		log.Printf("500")
		return
	}
	if code := req.FormValue("code"); code != "" {
		fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
		log.Printf("Code received1: %s", code)
		rw.(http.Flusher).Flush()
		log.Printf("Code received2: %s", code)
		//ch <- code
		handler.code = code
		log.Printf("Code received3: %s", code)
		return
	}
	if error := req.FormValue("error"); error == "access_denied" {
		fmt.Fprintf(rw, "<h1>Rejected</h1>User rejected authorization.")
		rw.(http.Flusher).Flush()
		log.Printf("User rejected authorization")
		return
	}
	log.Printf("no code")
	http.Error(rw, "", 500)
}*/

//Open URL in Chrome

func openURL(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
	tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "openURL", "Error opening URL in browser.")
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("./keys/client_secret_web.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	svc, err := gmail.New(client)
	if err != nil {
		tracelog.Errorf(fmt.Errorf("Exception At..."), "gmail", "gmailMain", "Unable to create Gmail service: %v", err)
	}

	var total int64
	msgs := []message{}
	pageToken := ""

	//queryStr := fmt.Sprintf("from:%v newer_than:%v", "confirmation@mail.hotels.com", "50d")

	emails := getAllStoresEmails()

	queryStr := buildQueryString(emails)

	for {
		req := svc.Users.Messages.List("me").Q(queryStr)
		if pageToken != "" {
			req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			tracelog.Errorf(fmt.Errorf("Exception At..."), "gmail", "gmailMain",
				"Unable to retrieve messages: %v", err)
		}
		if (r.Messages != nil) {
			tracelog.Info("gmail", "gmailMain", "Processing %v messages...\n", len(r.Messages))
		} else {
			tracelog.Info("gmail", "gmailMain", "No messages to process")
		}
		for _, m := range r.Messages {
			msg, err := svc.Users.Messages.Get("me", m.Id).Do()
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
}

func sortByDate(messages timeSlice) {

	date_sorted_messages := make(timeSlice, 0, len(messages))

	for _, d := range messages {
		date_sorted_messages = append(date_sorted_messages, d)
	}

	sort.Sort(date_sorted_messages)
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

func strToDate(str string) time.Time {
	layout := "Mon, 2 Jan 2006 15:04:05 -0700"

	t, err := time.Parse(layout, str)

	if err != nil {
		tracelog.Errorf(err, "gmail", "strToDate", "Can't convert to date")
	}

	return t
}

func convertDateToDateTime(messages timeSlice) {
	tracelog.Trace("gmail", "convertDateToDateTime", "Started conversion od string dates to time.Time")

	for _, msg := range messages {
		msg.dateTime = strToDate(msg.date)
	}

	tracelog.Trace("gmail", "convertDateToDateTime", "Finished conversion od string dates to time.Time")
}

func forwardToJetAnywhere(message message) {

}