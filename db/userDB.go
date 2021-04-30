package db

// Information attached to general users.
type User struct {
	Id    string
	Name  string
	Email string
}

func (d *Database) getUser(userID string) ([]User, error) {
	result := d.b.QueryRow("SELECT UserID, Username FROM Users WHERE UserID = ?", userID)
	var user User
	err := result.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return []User{}, user.Scan()
	}
	defer ProductRows.Close()

	var Products = make([]Product, 0)
	for ProductRows.Next() {
		var newProduct Product
		err = ProductRows.Scan(&newProduct.Name, &newProduct.Id, &newProduct.Description, &newProduct.Thumbnail, &newProduct.Price, &newProduct.Quantity, newProduct.MerchID)
		if err != nil {
			return []Product{}, err
		}
		Products = append(Products, newProduct)
	}
	return Products, nil
}
