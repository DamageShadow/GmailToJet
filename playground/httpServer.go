package main

import (
	"github.com/goinggo/tracelog"
	"fmt"
	"net/http/httptest"
	"net/http"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"time"
	"os/exec"
)

func main() {
	var (
		ctx context.Context
		config *oauth2.Config
	)

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
	defer ts.Close()

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
	fmt.Println(token)

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

