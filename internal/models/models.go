package models

type Password struct {
	Login    string
	Password string
	Tag      string
	ID       string
}

type Data struct {
	Data []byte
	Tag  string
	ID   string
}

type Text struct {
	Data string
	Tag  string
	ID   string
}

type CreditCard struct {
	CardNum string
	Exp     string
	Name    string
	CVV     string
	ID      string
}
