package main

import "time"

func buildForwardMessage(message message) (message){
	message.subject = "Fwd: " + message.subject
	message.dateTime = time.Now()

}
