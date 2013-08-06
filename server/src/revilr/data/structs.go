package data

type (
	Revil struct {
		Type    string
		Url     string
		Comment string
		Date    string
	}

	User struct {
		Username string
		Password []byte
	}
)