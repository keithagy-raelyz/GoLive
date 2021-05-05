package db

import (
	"errors"
	"fmt"
)

type Product struct {
	Id        string
	Name      string
	Quantity  int
	Thumbnail string
	Price     float64
	ProdDesc  string
	MerchID   string
	Sales     int // TODO flow change through to app/DB operators
}

func (d *Database) GetAllProducts() ([]Product, error) {
	ProductRows, err := d.b.Query("SELECT * FROM Products WHERE Quantity != 0")
	if err != nil {
		return []Product{}, err
	}
	defer ProductRows.Close()

	var Products = make([]Product, 0)
	for ProductRows.Next() {
		var newProduct Product
		err = ProductRows.Scan(&newProduct.Id, &newProduct.Name, &newProduct.Quantity, &newProduct.Thumbnail, &newProduct.Price, &newProduct.ProdDesc, &newProduct.MerchID)
		if err != nil {
			return []Product{}, err
		}
		Products = append(Products, newProduct)
	}
	return Products, nil
}

func (d *Database) GetProduct(prodID string) (Product, error) {
	var newProduct Product
	err := d.b.QueryRow("Select * from products where ProductID=?", prodID).Scan(&newProduct.Id, &newProduct.Name, &newProduct.Quantity, &newProduct.Thumbnail, &newProduct.Price, &newProduct.ProdDesc, &newProduct.MerchID)
	if err != nil {
		//TODO return custom error msg
		return Product{}, err
	}
	return newProduct, nil
}

func (d *Database) CreateProduct(product Product) error {
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("INSERT INTO Products (Product_Name,Quantity,Thumbnail,Price,ProdDesc,MerchantID) VALUES (?, ?,?, ?,?,?)", product.Name, product.Quantity, product.Thumbnail, product.Price, product.ProdDesc, product.MerchID)
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
	fmt.Println(product)
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("Update Products set Product_Name=?,Quantity=?,Thumbnail=?,Price=?,ProdDesc=? where ProductID=? AND MerchantID =?", product.Name, product.Quantity, product.Thumbnail, product.Price, product.ProdDesc, product.Id, product.MerchID)
	if err != nil {
		//TODO return custom error msg
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCount > 1 {
		//TODO return custom error msg
		return errors.New("More than 1 row is affected")
	}
	return nil
}

func (d *Database) DeleteProduct(prodID string, merchID string) error {
	//TODO think about inserting products that do not belong to the specific merchant, how do ensure integrity of the data created
	//TODO session stores ID, read from Session take ID from session not from browser
	res, err := d.b.Exec("DELETE FROM products where ProductID =? AND MerchantID =?", prodID, merchID)
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
