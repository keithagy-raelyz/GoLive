package db

import (
	"net/http"
)

//Auth internal management code
type Auth struct {
	session map[string]string
}

func (d *Database) InitializeAndGetAuth() *Auth {
	d.a = &Auth{}
	return d.a
}

func (a *Auth) Middleware(endPoint http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//authenticate here if authenticated
		if 1 == 1 {
			//send session
			endPoint.ServeHTTP(w, r)
		} else {
			//don't send session
			endPoint.ServeHTTP(w, r)
		}

	})
}

func (a *Auth) ManageSession() {

}

func (a *Auth) UpdateSession() {

}

//DB Interaction code
func (d *Database) GetSessions() {

}

func (d *Database) GetSession(userSession string) {

}

func (d *Database) CreateSession() {

}

func (d *Database) UpdateSession() {

}
func (d *Database) DeleteSession(sessionID string) {
	//TODO consider delete to User called by random Curl requests ie no Authentication
	//res, err := d.b.Exec("DELETE FROM sessions where sessionID =? ", sessionID)
	//if err != nil {
	//	//TODO return custom error msg
	//	return err
	//}
	//rowCount, err := res.RowsAffected()
	//if err != nil || rowCount != 1 {
	//	//TODO return custom error msg
	//	return err
	//}
	//return nil
}
