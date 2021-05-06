package db

// MerchantUser has User's account details, with MerchDesc of storefront.
type MerchantUser struct {
	User
	MerchDesc string
}

// Merchant contains MerchantUser (storefront details) and inventory.
type Merchant struct {
	MerchantUser
	Products []Product
}

func (d *Database) GetAllMerchants() ([]Merchant, error) {
	merchantRows, err := d.b.Query("SELECT merchantID, Username, MerchDesc FROM merchants")
	if err != nil {
		return []Merchant{}, err
	}
	defer merchantRows.Close()

	var merchants = make([]Merchant, 0)
	for merchantRows.Next() {
		var newMerchant Merchant
		err = merchantRows.Scan(&newMerchant.Id, &newMerchant.Name, &newMerchant.MerchDesc)
		if err != nil {
			return []Merchant{}, err
		}
		merchants = append(merchants, newMerchant)
	}
	return merchants, nil
}

func (d *Database) GetInventory(merchID string) (Merchant, []Product, error) {
	merchProdsRows, err := d.b.Query("SELECT * from (SELECT username, merchants.merchantid, merchants.MerchDesc, products.ProductID, products.Product_Name, products.Quantity, products.Thumbnail, products.price, products.ProdDesc, products.Sales from merchants LEFT JOIN products on products.merchantid = merchants.merchantid) AS joinTable WHERE merchantid = ?;", merchID)
	if err != nil {
		// fmt.Println("Query error", err)
		return Merchant{}, []Product{}, err
	}
	defer merchProdsRows.Close()

	var merchProds []Product
	var merch = Merchant{}
	for merchProdsRows.Next() {
		var p Product
		// TODO Need to fix so merch only gets scanned ONCE
		err = merchProdsRows.Scan(&merch.Name, &merch.Id, &merch.MerchDesc, &p.Id, &p.Name, &p.Quantity, &p.Thumbnail, &p.Price, &p.ProdDesc, &p.Sales)
		if err != nil {
			// fmt.Println("Scan error", err)
			return merch, merchProds, err
		}
		merchProds = append(merchProds, p)
	}
	return merch, merchProds, nil
}

func (d *Database) CheckMerchant(merchant MerchantUser) error {
	var m MerchantUser
	err := d.b.QueryRow("SELECT username,email FROM users where Username=? OR email=?", merchant.Name, merchant.Email).Scan(m.Name, m.Email)
	if err != nil {
		//TODO return custom error msg
		return err
	}
	return nil
}

func (d *Database) CreateMerchant(merchant MerchantUser, password string) error {
	res, err := d.b.Exec("INSERT INTO merchants (username,password,email,MerchDesc) VALUES (?, ?,?, ?)", merchant.Name, password, merchant.Email, merchant.MerchDesc)
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

func (d *Database) UpdateMerchant(merchant MerchantUser) error {
	res, err := d.b.Exec("UPDATE merchants set username=?,email=?,MerchDesc=? where MerchantID=?", merchant.Name, merchant.Email, merchant.MerchDesc, merchant.Id)
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

func (d *Database) DeleteMerchant(merchID string) error {
	tx, err := d.b.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement1, err := tx.Prepare("DELETE FROM Products where MerchantID = ?")
	if err != nil {
		return err
	}
	defer statement1.Close()

	statement1.Exec(merchID)

	statement2, err := tx.Prepare("DELETE FROM Merchants where MerchantID = ?")
	if err != nil {
		return err
	}
	defer statement2.Close()

	statement2.Exec(merchID)

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
