package main

import (
	"fmt"
	"github.com/goinggo/tracelog"
	"google.golang.org/api/gmail/v1"
)

func processMessagesInInbox(gmailService *gmail.Service) {

	msgs := []message{}
	pageToken := ""

	emails := getAllStoresEmails()

	queryStr := buildQueryString(emails)

	fmt.Println("Search string is = ", queryStr)

	for {
		rq, _ := gmailService.Users.GetProfile("me").Do()

		fmt.Println("Profile email = ", rq.EmailAddress)

		req := gmailService.Users.Messages.List("me").Q(queryStr)

		fmt.Println("Request = ", req)
		if pageToken != "" {
			fmt.Println("pageToken before = ", pageToken)
			req.PageToken(pageToken)
			fmt.Println("pageToken after = ", pageToken)
		}

		r, err := req.Do()
		if err != nil {
			tracelog.Errorf(errorf, "gmail", "gmailMain",
				"Unable to retrieve messages: %v", err)
		}
		if r.Messages != nil {
			tracelog.Info("gmail", "gmailMain", "Processing %v messages...\n", len(r.Messages))
		} else {
			tracelog.Info("gmail", "gmailMain", "No messages to process")
		}
		for _, m := range r.Messages {
			msg, err := gmailService.Users.Messages.Get("me", m.Id).Do()
			fmt.Println("Msgs", msg)
			if err != nil {
				tracelog.Errorf(errorf, "gmail",
					"gmailMain", "Unable to retrieve message %v: %v", m.Id, err)
			}

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
	sortByDate(msgs)

	count := 0
	for _, m := range msgs {
		count++
		tracelog.Info("gmail", "gmailMain", "\nMessage URL: https://mail.google.com/mail/u/0/#all/%v\n"+
			"Subject: %q ,Size: %v, Date: %v, Snippet: %q\n",
			m.gmailID, m.subject, m.size, m.date, m.snippet)

		forwardToJetAnywhere(m)
	}
}

func forwardToJetAnywhere(message message) {

	fmt.Println("Mock Forward to Jet")

}
