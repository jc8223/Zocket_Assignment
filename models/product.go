package models

import "time"

type Product struct {
	Product_ID               int       `json:"product_id"`
	Product_Name             string    `json:"product_name"`
	Product_Description      string    `json:"product_description"`
	Product_Images           []string  `json:"product_images"`
	Product_Price            float64   `json:"product_price"`
	Compressed_Product_Images []string  `json:"compressed_product_images"`
	Created_At               time.Time `json:"created_at"`
	Updated_At               time.Time `json:"updated_at"`
}
