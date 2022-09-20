package validation

type Error struct {
}

func (e Error) Error() string {
	return "bla"
}
