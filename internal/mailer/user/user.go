package user

type User struct {
	ID      string
	Email   string
	Name    string
	Surname string
	Meta    map[string]any
}
