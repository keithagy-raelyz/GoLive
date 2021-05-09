package db

import (
	"errors"
)

type Product struct {
	Id        string
	Name      string
	Quantity  int
	Thumbnail string
	Price     float64
	ProdDesc  string
	MerchID   string
	Sales     int
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
		err = ProductRows.Scan(&newProduct.Id, &newProduct.Name, &newProduct.Quantity, &newProduct.Thumbnail, &newProduct.Price, &newProduct.ProdDesc, &newProduct.MerchID, &newProduct.Sales)
		if err != nil {
			return []Product{}, err
		}
		Products = append(Products, newProduct)
	}
	return Products, nil
}

func (d *Database) GetProduct(prodID string) (Product, error) {
	var newProduct Product
	err := d.b.QueryRow("Select * from products where ProductID=?", prodID).Scan(&newProduct.Id, &newProduct.Name, &newProduct.Quantity, &newProduct.Thumbnail, &newProduct.Price, &newProduct.ProdDesc, &newProduct.MerchID, &newProduct.Sales)
	if err != nil {
		return Product{}, err
	}
	return newProduct, nil
}

func (d *Database) CreateProduct(product Product) error {
	res, err := d.b.Exec("INSERT INTO Products (Product_Name,Quantity,Thumbnail,Price,ProdDesc,MerchantID,Sales) VALUES (?, ?,?, ?,?,?)", product.Name, product.Quantity, product.Thumbnail, product.Price, product.ProdDesc, product.MerchID, 0)
	if err != nil {
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		return err
	}
	return nil
}

func (d *Database) UpdateProduct(product Product) error {
	res, err := d.b.Exec("Update Products set Product_Name=?,Quantity=?,Thumbnail=?,Price=?,ProdDesc=?,Sales=? where ProductID=? AND MerchantID =?", product.Name, product.Quantity, product.Thumbnail, product.Price, product.ProdDesc, product.Sales+1, product.Id, product.MerchID)
	if err != nil {
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCount > 1 {
		return errors.New("more than 1 row is affected")
	}
	return nil
}

func (d *Database) DeleteProduct(prodID string, merchID string) error {
	res, err := d.b.Exec("DELETE FROM products where ProductID =? AND MerchantID =?", prodID, merchID)
	if err != nil {
		return err
	}
	rowCount, err := res.RowsAffected()
	if err != nil || rowCount != 1 {
		return err
	}
	return nil
}
