package main

type ProductType struct {
	ID   uint   `gorm:"column:id;primary_key;"`
	Name string `gorm:"column:name;"`
}

type Product struct {
	ID            uint   `gorm:"column:id;primary_key;"`
	ProductTypeID uint   `gorm:"column:product_type_id;index;"`
	Name          string `gorm:"column:name;"`

	ProductType ProductType
}
