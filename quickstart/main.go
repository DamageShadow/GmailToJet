package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/goinggo/tracelog"
)

// Flags
var (
	clientID = flag.String("clientid", "965925130911-lb3hu48uj4u5ab1l8acs7e3u40fqlve4.apps.googleusercontent.com",
		"OAuth 2.0 Client ID.  If non-empty, overrides --clientid_file")
	clientIDFile = flag.String("clientid-file", "",
		"Name of a file containing just the project's OAuth 2.0 Client ID from https://developers.google.com/console.")
	secret = flag.String("secret", "", "OAuth 2.0 Client Secret.  If non-empty, overrides --secret_file")
	secretFile = flag.String("secret-file",
		"./keys/client_secret.json",
		"Name of a file containing just the project's OAuth 2.0 Client Secret from https://developers.google.com/console.")
	cacheToken = flag.Bool("cachetoken", true, "cache the OAuth 2.0 token")
	debug = flag.Bool("debug", false, "show HTTP traffic")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: go-api-demo <api-demo-name> [api name args]\n\nPossible APIs:\n\n")
	for n := range demoFunc {
		fmt.Fprintf(os.Stderr, "  * %s\n", n)
	}
	os.Exit(2)
}

func main() {

	tracelog.StartFile(tracelog.LevelInfo, "./logs", 1)

	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	name := flag.Arg(0)
	demo, ok := demoFunc[name]
	if !ok {
		usage()
	}

	config := &oauth2.Config{
		ClientID:     valueOrFileContents(*clientID, *clientIDFile),
		ClientSecret: valueOrFileContents(*secret, *secretFile),
		Endpoint:     google.Endpoint,
		Scopes:       []string{demoScope[name]},
		RedirectURL:  "http://localhost:4567/oauth2callback",
	}

	ctx := context.Background()
	/*
		b, err := ioutil.ReadFile("client_secret.json")

		config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
		if err != nil {
			tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "main","Unable to parse client secret file to config: %v", err)
		}*/

	if *debug {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &logTransport{http.DefaultTransport},
		})
	}
	c := newOAuthClient(ctx, config)
	demo(c, flag.Args()[1:])
}

var (
	demoFunc = make(map[string]func(*http.Client, []string))
	demoScope = make(map[string]string)
)

func registerDemo(name, scope string, main func(c *http.Client, argv []string)) {
	if demoFunc[name] != nil {
		panic(name + " already registered")
	}
	demoFunc[name] = main
	demoScope[name] = scope
}

func osUserCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd":
		return filepath.Join(os.Getenv("HOME"), ".cache")
	}
	tracelog.Info("main", "osUserCacheDir", "TODO: osUserCacheDir on GOOS %q", runtime.GOOS)
	return "."
}

func tokenCacheFile(config *oauth2.Config) (string) {
	hash := fnv.New32a()
	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))
	fn := fmt.Sprintf("go-api-demo-tok%v", hash.Sum32())
	return filepath.Join(osUserCacheDir(), url.QueryEscape(fn))
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	if !*cacheToken {
		return nil, errors.New("--cachetoken is false")
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "saveToken", "Warning: failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}

func newOAuthClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile := tokenCacheFile(config)

	token, err := tokenFromFile(cacheFile)
	if err != nil {
		token = tokenFromWeb(ctx, config)
		saveToken(cacheFile, token)
	} else {
		tracelog.Info("main", "newOAuthClient", "Using cached token %#v from %q", token, cacheFile)
	}

	return config.Client(ctx, token)
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if code := req.FormValue("code"); code != "" {
			tracelog.Info("main", "tokenFromWeb", "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		if req.URL.Path == "/favicon.ico" {
			tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "tokenFromWeb", "404 - Not found")
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "tokenFromWeb", "State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}

		tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "tokenFromWeb", "no code")
		http.Error(rw, "", 500)
	}))

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	go openURL(authURL)
	tracelog.Info("main", "tokenFromWeb", "Authorize this app at: %s", authURL)
	code := <-ch
	tracelog.Info("main", "tokenFromWeb", "Got response code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		tracelog.Info("main", "tokenFromWeb", "Token exchange error: %v", err)
	}

	defer ts.Close()

	return token
}

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

func valueOrFileContents(value string, filename string) string {
	if value != "" {
		return value
	}
	slurp, err := ioutil.ReadFile(filename)
	if err != nil {
		tracelog.Errorf(fmt.Errorf("Exception At..."), "main", "valueOrFileContents", "Error reading %q: %v", filename, err)
	}
	return strings.TrimSpace(string(slurp))
}
