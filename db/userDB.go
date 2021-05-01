package db

// Information attached to general users.
type User struct {
	Id       string
	Name     string
	Email    string
	Password string
}

func (d *Database) getUser(userID string) (User, error) {
	result := d.b.QueryRow("SELECT UserID, Username FROM Users WHERE UserID = ?", userID)
	var user User

	return user, result.Scan(&user.Id, &user.Name, &user.Email)

}

func (d *Database) CreateUser(user User) error {
	res, err := d.b.Exec("INSERT INTO Users (Username,Password,Email) VALUES (?,?,?)", user.Name, user.Password, user.Email)
	if err != nil {
		//TODO return custom error msg
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		//TODO return custom error msg
		return err
	}
	return nil
}

func (d *Database) CheckUser(user User) error {
	var u User
	err := d.b.QueryRow("SELECT username,email FROM users where Username=? OR email=?", user.Name, user.Email).Scan(u.Name, u.Email)
	if err != nil {
		//TODO return custom error msg
		return err
	}
	return nil
}

func (d *Database) UpdateUser(user User) error {
	//TODO consider updates to User called by random Curl requests ie no Authentication
	res, err := d.b.Exec("Update Users set Username=?,Password=?,Email=? where UserID=?", user.Name, user.Password, user.Email, user.Id)
	if err != nil {
		//TODO return custom error msg
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		//TODO return custom error msg
		return err
	}
	return nil
}

func (d *Database) DeleteUser(userID string) error {
	//TODO consider delete to User called by random Curl requests ie no Authentication
	res, err := d.b.Exec("DELETE FROM users where UserID =? ", userID)
	if err != nil {
		//TODO return custom error msg
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		//TODO return custom error msg
		return err
	}
	return nil
}
