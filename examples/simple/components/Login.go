package components

type Login struct {
	Email    string
	Password string
}

func (l *Login) Update(p any) {
	l.Email = "email@from.remote"
}
