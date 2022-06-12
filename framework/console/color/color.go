package color

import "fmt"

func Text(code int, value interface{}) string {
	return fmt.Sprintf("\u001b[38;5;%dm%s\u001b[0m", code, value)
}

func Bg(code int, value interface{}) string {
	return fmt.Sprintf("\u001b[48;5;%dm%s\u001b[0m", code, value)
}
