package main

import (
	"encoding/json"
	"fmt"
	"github.com/goinggo/tracelog"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func difference(slice1 []string, slice2 []string) ([]string, []string) {
	diffStr1 := []string{}
	diffStr2 := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 2
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr1 = append(diffStr1, mKey)
		}
		if mVal == 2 {
			diffStr2 = append(diffStr2, mKey)
		}
	}

	return diffStr1, diffStr2
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

func sortByDate(messages timeSlice) {

	date_sorted_messages := make(timeSlice, 0, len(messages))

	for _, d := range messages {
		date_sorted_messages = append(date_sorted_messages, d)
	}

	sort.Sort(date_sorted_messages)
}

func generateUUID() *uuid.UUID {

	u, err := uuid.NewV4()
	if err != nil {
		tracelog.Errorf(err, "gmail", "generateUUID", "Can't generate UUID")
	}

	return u
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

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() string {
	/*usr, err := user.Current()
	 */
	tokenCacheDir := filepath.Join("./", ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)

	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-quickstart.json"))
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

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getCachedClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile := tokenCacheFile()

	tok, _ := tokenFromFile(cacheFile)

	return config.Client(ctx, tok)
}

//
//
