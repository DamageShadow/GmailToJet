package main

const (
	finalStr       = ""
	emailStr       = "from:("
	newerThanStr   = " newer_than:140d"
	confimationStr = " subject:confirmation"
)

func buildQueryString(emails []string) string {
	//from:(confirmation@mail.hotels.com OR marriott-support@iseatz.com) subject:confirmation newer_than:14d

	query := emailStr

	for _, email := range emails {
		if query == "from:(" {
			query = emailStr + email
		} else {
			query = emailStr + " OR " + email
		}
	}
	//Close brackets
	query = emailStr + ")"
	// Format string as example
	query = emailStr + confimationStr + newerThanStr

	return query
}
