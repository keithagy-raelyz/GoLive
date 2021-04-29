package db

import (
	"database/sql"

	"github.com/keithagy-raelyz/GoLive/app"
)

type Database struct {
	b *sql.DB
}

func (d *Database) GetAllMerchants() ([]app.Merchant, error) {
	merchantRows, err := d.b.Query("SELECT merchantID, Username, description FROM merchants")
	if err != nil {
		return app.Merchant{}, err
	}
	defer merchantRows.Close()

	var merchants = make([]app.Merchant, 0)
	for merchantRows.Next() {
		var newMerchant app.Merchant
		err = merchantRows.Scan(&newMerchant.Id, &newMerchant.Name, &newMerchant.Description)
		if err != nil {
			return app.Merchant{}, err
		}
		merchants = append(merchants, newMerchant)
	}
	return merchants, nil
}

func (d *Database) GetInventory(merchID int) ([]app.Product, error) {
	merchProdsRows, err := d.b.Query("SELECT username, merchants.merchantid, merchants.description, products.ProductID, products.Product_Name, products.Quantity, products.Image, products.price,products.Description from merchants LEFT JOIN products on products.merchantid = merchants.merchantid where merchants.merchantid = ?;", merchID)
	if err != nil {
		return app.Product{}, err
	}
	defer merchProdsRows.Close()

	var merchProds = make([]app.Product, 0)
	var merch = &app.Merchant{}
	for merchProdsRows.Next() {
		var p app.Product
		// TODO Need to fix so merch only gets scanned ONCE
		err = merchProdsRows.Scan(&merch.Name, &merch.Id, &merch.Description, &p.Id, &p.Name, &p.Quantity, &p.Thumbnail, &p.Price, &p.Description)
		if err != nil {
			return app.Product{}, err
		}
		merchProds = append(merchProds, p)
	}
	return merchProds, nil
}
