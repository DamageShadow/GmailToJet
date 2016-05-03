package main

func buildQueryString(emails []string) (string) {
	//from:(confirmation@mail.hotels.com OR marriott-support@iseatz.com) subject:confirmation newer_than:14d
	finalStr := ""
	emailStr := "from:("
	newerThanStr := " newer_than:14d"
	confimationStr := " subject:confirmation"

	for email := range emails {
		if (emailStr == "") {
			emailStr = emailStr + email
		} else {
			emailStr = emailStr + " OR " + email
		}
	}
	//Close brackets
	emailStr = emailStr + ")"
	// Format string as example
	finalStr = emailStr + confimationStr + newerThanStr

	return finalStr
}
