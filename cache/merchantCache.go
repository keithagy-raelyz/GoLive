package cache

import (
	"GoLive/db"
	"time"
)

//merchantCache stores recently accessed merchant pages.
type merchantCache struct {
	cache
}

//getAllmerchants pulls a fresh row from teh DB if the total cached merchants is less than 5.
func (m *merchantCache) getAllMerchants(database *db.Database) ([]*CachedMerchant, error) {
	var count int
	var allMerchants []*CachedMerchant
	for _, v := range (*m).cache {
		allMerchants = append(allMerchants, v.(*CachedMerchant))
		count++
	}
	if count < 5 {
		var allMerchantsFromDB []*CachedMerchant
		allMerch, err := database.GetAllMerchants()
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(allMerch); j++ {
			merchantToCache := &CachedMerchant{
				&allMerch[j], expiry{time.Now().Add(SessionLife * time.Minute)},
			}
			(*m).cache[allMerch[j].Id] = merchantToCache
			allMerchantsFromDB = append(allMerchantsFromDB, merchantToCache)
		}
		return allMerchantsFromDB, nil
	} else {
		return allMerchants, nil
	}
}

//updateMerchant updates the merchant stored in the cache and the DB.
func (m *merchantCache) updateMerchant(merchant db.MerchantUser, database *db.Database) error {
	if cachedMerch, ok := m.cache[merchant.Id]; ok {
		cachedMerch.(*CachedMerchant).Update(merchant)
		return database.UpdateMerchant(merchant)
	}
	return database.UpdateMerchant(merchant)
}

//CachedMerchant stores the data of a recently accessed merchant.
type CachedMerchant struct {
	merchantPage *db.Merchant
	expiry
}

//getKey returns the merchant's ID
func (c *CachedMerchant) getKey() string {
	return c.merchantPage.Id
}

//deleteItem deletes an item within the cache.
func (c *CachedMerchant) deleteItem(prodID string) {
	size := len(c.merchantPage.Products)
	if size == 1 {
		c.merchantPage.Products = make([]db.Product, 0)
		return
	}
	for i := 0; i < size; i++ {
		if c.merchantPage.Products[i].Id == prodID {
			if i == 0 {
				c.merchantPage.Products = c.merchantPage.Products[i+1:]
				return
			} else if i == len(c.merchantPage.Products)-1 {
				c.merchantPage.Products = c.merchantPage.Products[0:i]
				return
			} else {
				c.merchantPage.Products = append(c.merchantPage.Products[0:i], c.merchantPage.Products[i+1:]...)
				return
			}
		}
	}
}

//GetCachedMerchant returns the stored merchant in the cache.
func (c *CachedMerchant) GetCachedMerchant() *db.Merchant {
	return c.merchantPage
}

//deleteProduct deletes the product stored in the cachedmerchant and in the database.
func (m *merchantCache) deleteProduct(prodID string, merchID string, database *db.Database) error {
	if merch, ok := m.cache[merchID]; ok {
		merch.(*CachedMerchant).deleteItem(prodID)
	}
	return database.DeleteProduct(prodID, merchID)
}

//deleteMerchant deletes the merchant in the cache and in the database.
func (m *merchantCache) deleteMerchant(merchID string, database *db.Database) error {
	delete(m.cache, merchID)
	return database.DeleteMerchant(merchID)
}

//Update updates the merchant stored in the cache.
func (c *CachedMerchant) Update(merchant db.MerchantUser) {
	c.merchantPage.Id = merchant.Id
	c.merchantPage.Name = merchant.Name
	c.merchantPage.Email = merchant.Email
	c.merchantPage.Password = merchant.Password
}
