package main

import (
	"net/http"
	"sort"

	gmail "google.golang.org/api/gmail/v1"
	"github.com/goinggo/tracelog"
	"fmt"
	"time"
)

func init() {
	registerDemo("gmail", gmail.MailGoogleComScope, gmailMain)
}

type message struct {
	size     int64
	gmailID  string
	date     string // retrieved from message header
	snippet  string
	subject  string
	dateTime time.Time
}

type timeSlice []message

// gmailMain is an example that demonstrates calling the Gmail API.
// It iterates over all messages of a user that are larger
// than 5MB, sorts them by size, and then interactively asks the user to
// choose either to Delete, Skip, or Quit for each message.
//
// Example usage:
//   go build -o go-api-demo *.go
//   go-api-demo -clientid="my-clientid" -secret="my-secret" gmail
func gmailMain(client *http.Client, argv []string) {

	if len(argv) != 0 {
		tracelog.Errorf(fmt.Errorf("Exception At..."), "gmail", "gmailMain", "Proper Usage of arguments: gmail")
		return
	}

	svc, err := gmail.New(client)
	if err != nil {
		tracelog.Errorf(fmt.Errorf("Exception At..."), "gmail", "gmailMain", "Unable to create Gmail service: %v", err)
	}

	var total int64
	msgs := []message{}
	pageToken := ""

	//queryStr := fmt.Sprintf("from:%v newer_than:%v", "confirmation@mail.hotels.com", "50d")

	//emails := getAllStoresEmails()

	queryStr := ""//buildQueryString(emails)

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

		tracelog.Info("gmail", "gmailMain", "Processing %v messages...\n", len(r.Messages))
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