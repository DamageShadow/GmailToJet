package main

import "time"

func buildForwardMessage(message message) (message) {
	message.subject = "Fwd: " + message.subject
	message.dateTime = time.Now()
	//TODO work with gmail msg instead of messages
	return message
}
