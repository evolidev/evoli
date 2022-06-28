package components

type Login struct {
	Email    string
	Password string
}

func (l *Login) Update() {
	l.Email = "email@from.remote"
}
