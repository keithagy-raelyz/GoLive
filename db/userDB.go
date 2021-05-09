package db

import "errors"

// Information attached to general users.
type User struct {
	Id       string
	Name     string
	Email    string
	Password string
}

func (d *Database) GetUsers() ([]User, error) {
	result, err := d.b.Query("SELECT UserID, Username, Email FROM Users")
	if err != nil {
		return []User{}, err
	}
	var users []User
	for result.Next() {
		var u User
		err := result.Scan(&u.Id, &u.Name, &u.Email)
		if err != nil {
			return []User{}, err
		}
		users = append(users, u)
	}
	return users, nil

}

func (d *Database) GetUser(username string) (User, error) {
	result := d.b.QueryRow("SELECT * FROM Users WHERE Username = ?", username)
	var user User

	return user, result.Scan(&user.Id, &user.Name, &user.Password, &user.Email)

}

func (d *Database) CreateUser(user User, password string) error {
	res, err := d.b.Exec("INSERT INTO Users (Username,Password,Email) VALUES (?,?,?)", user.Name, password, user.Email)
	if err != nil {
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		return err
	}
	return nil
}

func (d *Database) CheckUser(user User) error {

	res, err := d.b.Query("SELECT username,email FROM users where Username=? OR email=?", user.Name, user.Email)
	if err != nil {
		return err
	}
	if res.Next() {
		return errors.New("User exists")
	}
	return nil
}

func (d *Database) UpdateUser(user User) error {
	res, err := d.b.Exec("Update Users set Username=?,Email=? where UserID=?", user.Name, user.Email, user.Id)
	if err != nil {
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		return err
	}
	return nil
}

func (d *Database) DeleteUser(userID string) error {
	res, err := d.b.Exec("DELETE FROM users where UserID =? ", userID)
	if err != nil {
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		return err
	}
	return nil
}
