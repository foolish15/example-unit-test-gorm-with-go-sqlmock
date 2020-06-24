package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func getConnection(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sdb, mock, err := sqlmock.New()
	assert.Nil(t, err)
	db, err := gorm.Open("mysql", sdb)
	assert.Nil(t, err)
	return db, mock
}

func TestMigratation(t *testing.T) {
	db, mock := getConnection(t)
	mock.ExpectQuery("SHOW TABLES FROM `` WHERE `Tables_in_` = ?").
		WithArgs("products").
		WillReturnRows(mock.NewRows([]string{"Tables_in_"}).AddRow("products"))
	mock.ExpectQuery("SHOW COLUMNS FROM `products` FROM `` WHERE Field = ?").
		WithArgs("id").
		WillReturnRows(
			mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
				AddRow("id", "int(10) unsigned", "NO", "PRI", nil, "auto_increment"),
		)
	mock.ExpectQuery("SHOW COLUMNS FROM `products` FROM `` WHERE Field = ?").
		WithArgs("product_type_id").
		WillReturnRows(
			mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}),
		)
	//Add column when mission
	mock.ExpectExec("ALTER TABLE `products` ADD `product_type_id` int unsigned;").
		WillReturnResult(
			sqlmock.NewResult(0, 0),
		)
	mock.ExpectQuery("SHOW COLUMNS FROM `products` FROM `` WHERE Field = ?").
		WithArgs("name").
		WillReturnRows(
			mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
				AddRow("name", "varchar(255)", "YES", "", nil, ""),
		)
	mock.ExpectQuery("SHOW COLUMNS FROM `products` FROM `` WHERE Field = ?").
		WithArgs("product_type").
		WillReturnRows(
			mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}),
		)
	mock.ExpectQuery("SHOW INDEXES FROM `products` FROM `` WHERE Key_name = ?").
		WithArgs("idx_products_product_type_id").
		WillReturnRows(
			mock.NewRows([]string{"Table", "Non_unique", "Key_name", "Seq_in_index", "Column_name", "Collation", "Cardinallity", "Sub_part", "Packed", "Null", "Index_type", "Comment", "Index_comment"}),
		)
	//Create index when mission
	mock.ExpectExec("CREATE INDEX idx_products_product_type_id ON `products`(.+)").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SHOW TABLES FROM `` WHERE `Tables_in_` = ?").
		WithArgs("product_types").
		WillReturnRows(mock.NewRows([]string{"Tables_in_"}).AddRow("product_types"))
	mock.ExpectQuery("SHOW COLUMNS FROM `product_types` FROM `` WHERE Field = ?").
		WithArgs("id").
		WillReturnRows(
			mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
				AddRow("id", "int(10) unsigned", "NO", "PRI", nil, "auto_increment"),
		)
	mock.ExpectQuery("SHOW COLUMNS FROM `product_types` FROM `` WHERE Field = ?").
		WithArgs("name").
		WillReturnRows(
			mock.NewRows([]string{"Field", "Type", "Null", "Key", "Default", "Extra"}).
				AddRow("name", "varchar(255)", "YES", "", nil, ""),
		)

	err := db.LogMode(true).AutoMigrate(&Product{}, &ProductType{}).Error
	assert.Nil(t, err)
}
