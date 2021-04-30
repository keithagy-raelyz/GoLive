package db

// TODO Work around circular dependency (db needs types declared in app); for now we are redeclaring the types
// Information attached to general users.
type User struct {
	Id    string
	Name  string
	Email string
}
