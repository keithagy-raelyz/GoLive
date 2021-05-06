package db

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
