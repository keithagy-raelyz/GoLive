package db

// TODO Work around circular dependency (db needs types declared in app); for now we are redeclaring the types
type Product struct {
	Name        string
	Id          int
	Description string
	Thumbnail   string
	Price       float64
	Quantity    int
	MerchID     int
}

func (d *Database) GetAllProducts() ([]Product, error) {
	ProductRows, err := d.b.Query("SELECT ProductID, Product_Name, description FROM Products")
	if err != nil {
		return []Product{}, err
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

func (d *Database) GetProduct(prodID string) (Product, error) {
	var p Product
	err := d.b.QueryRow("Select * from products where ProductID=?", prodID).Scan(&p.Id, &p.Name, &p.Quantity, &p.Thumbnail, &p.Price, &p.Description)
	if err != nil {
		//TODO return custom error msg
		return Product{}, err
	}
	return p, nil
}

func (d *Database) CreateProduct(product Product) error {
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("INSERT INTO Products (Product_Name,Quantity,Image,Price,Description,MerchantID) VALUES (?, ?,?, ?,?,?)", product.Name, product.Quantity, product.Thumbnail, product.Price, product.Description, product.MerchID)
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

func (d *Database) UpdateProduct(product Product) error {
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("Update Products set Product_Name=?,Quantity=?,Image=?,Price=?,Description=? where ProductID=? AND MerchantID =?", product.Name, product.Quantity, product.Thumbnail, product.Price, product.Description, product.Id, product.MerchID)
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
