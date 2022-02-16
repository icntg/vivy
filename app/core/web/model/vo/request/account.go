package request

type UserPass struct {
	user    string
	pass    string
	captcha string
}

type GoogleToken string
