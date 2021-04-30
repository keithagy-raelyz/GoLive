package db

import "github.com/keithagy-raelyz/GoLive/app"

func (d *Database) GetAllProducts() ([]app.Product, error) {
	ProductRows, err := d.b.Query("SELECT ProductID, Product_Name, description FROM Products")
	if err != nil {
		return []app.Product{}, err
	}
	defer ProductRows.Close()

	var Products = make([]app.Product, 0)
	for ProductRows.Next() {
		var newProduct app.Product
		err = ProductRows.Scan(&newProduct.Name, &newProduct.Id, &newProduct.Description, &newProduct.Thumbnail, &newProduct.Price, &newProduct.Quantity)
		if err != nil {
			return []app.Product{}, err
		}
		Products = append(Products, newProduct)
	}
	return Products, nil
}

func (d *Database) GetProduct(prodID string) (app.Product, error) {
	var p app.Product
	err := d.b.QueryRow("Select * from products where ProductID=?", prodID).Scan(&p.Id, &p.Name, &p.Quantity, &p.Thumbnail, &p.Price, &p.Description)
	if err != nil {
		//TODO return custom error msg
		return app.Product{}, err
	}
	return p, nil
}

func (d *Database) CreateProduct(product app.Product, merchID string) error {
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("INSERT INTO Products (Product_Name,Quantity,Image,Price,Description,MerchantID) VALUES (?, ?,?, ?,?,?)", product.Name, product.Quantity, product.Thumbnail, product.Price, product.Description, merchID)
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

func (d *Database) UpdateProduct(product app.Product, merchID string) error {
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("Update Products set Product_Name=?,Quantity=?,Image=?,Price=?,Description=? where ProductID=? AND MerchantID =?", product.Name, product.Quantity, product.Thumbnail, product.Price, product.Description, product.Id, merchID)
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

func (d *Database) DeleteProduct(prodID string, merchID string) error {
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("DELETE FROM products where ProductID =? AND merchID =?", prodID, merchID)
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
