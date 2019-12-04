package db

import (
	"fmt"
	"time"
)

type User struct {
	ID        int
	Email     string `pg:",unique,notnull"`
	Password  string `pg:",notnull"`
	Role      string
	Addresses []*Address `pg:"fk:user_id"`
}

type Address struct {
	DeletedAt  time.Time `pg:",soft_delete"`
	ID         int
	Name       string `pg:",notnull"`
	Line1      string `pg:",notnull"`
	Line2      string
	Line3      string
	City       string `pg:",notnull"`
	Region     string `pg:",notnull"`
	PostalCode string `pg:",notnull"`
	Country    string `pg:",notnull"`
	UserID     int
	User       *User
}

func (address Address) String() string {
	return fmt.Sprintf("%s-|-%s-|-%s-|-%s-|-%s-|-%s-|-%s",
		address.Line1, address.Line2, address.Line3, address.City, address.Region, address.PostalCode, address.Country)
}

func (address Address) Raw() interface{} {
	return address
}

type Image struct {
	DeletedAt time.Time `pg:",soft_delete"`
	ID        int
	Name      string
	Raw       string
	Thumbnail string
	Height600 string
}

type Product struct {
	DeletedAt       time.Time `pg:",soft_delete"`
	ID              int
	Slug            string `pg:",unique,notnull"`
	Name            string `pg:",notnull"`
	Description     string `pg:",notnull"`
	Details         string
	Published       bool
	ProductImages   []*ProductImage   `pg:"fk:product_id"`
	ProductOptions  []*ProductOption  `pg:"fk:product_id"`
	ProductVariants []*ProductVariant `pg:"fk:product_id"`
}

type ProductImage struct {
	ProductID int
	Product   *Product
	ImageID   int
	Image     *Image
}

type ProductOption struct {
	DeletedAt time.Time `pg:",soft_delete"`
	ID        int
	Label     string                `pg:",notnull"`
	Values    []*ProductOptionValue `pg:"fk:product_option_id"`
	ProductID int                   `pg:",notnull"`
	Product   *Product
}

type ProductOptionValue struct {
	DeletedAt       time.Time `pg:",soft_delete"`
	ID              int
	Value           string `pg:",notnull"`
	ProductOptionID int    `pg:",notnull"`
	ProductOption   *ProductOption
}

type ProductVariant struct {
	DeletedAt       time.Time `pg:",soft_delete"`
	ID              int
	Name            string
	Price           int                     `pg:",notnull"`
	Length          float64                 `pg:",notnull"`
	Width           float64                 `pg:",notnull"`
	Height          float64                 `pg:",notnull"`
	Weight          float64                 `pg:",notnull"`
	SelectedOptions []*ProductVariantOption `pg:"fk:product_variant_id"`
	ProductID       int                     `pg:",notnull"`
	Product         *Product
	ShipsFromID     int
	ShipsFrom       *Address
	Images          []*ProductVariantImage
}

type ProductVariantImage struct {
	ProductVariantID int
	ProductVariant   *ProductVariant
	ImageID          int
	Image            *Image
}

type ProductVariantOption struct {
	DeletedAt            time.Time `pg:",soft_delete"`
	ProductOptionValueID int       `pg:",notnull"`
	ProductOptionValue   *ProductOptionValue
	ProductVariantID     int `pg:",notnull"`
	ProductVariant       *ProductVariant
	ProductID            int `pg:",notnull"`
	Product              *Product
}

type Transaction struct {
	ID                  int
	Subtotal            int `pg:",notnull"`
	Taxes               int `pg:",notnull"`
	Shipping            int `pg:",notnull"`
	Total               int `pg:",notnull"`
	BraintreeID         string
	ShippoRateID        string
	ShippoTransactionID string
	UserID              int
	User                *User
	Addresses           *TransactionAddressInfo `pg:"fk:transaction_id"`
	LineItems           []*TransactionLineItem  `pg:"fk:transaction_id"`
	Status              []*TransactionStatus    `pg:"fk:transasction_id"`
}

type TransactionAddressInfo struct {
	ID                int
	TransactionID     int `pg:",notnull"`
	Transaction       *Transaction
	BillingAddressID  int `pg:",notnull"`
	BillingAddress    *Address
	ShippingAddressID int `pg:",notnull"`
	ShippingAddress   *Address
}

type TransactionLineItem struct {
	ID               int
	TransactionID    int `pg:",notnull"`
	Transaction      *Transaction
	ProductVariantID int `pg:",notnull"`
	ProductVariant   *ProductVariant
	Price            int `pg:",notnull"`
	Quantity         int `pg:",notnull"`
}

type TransactionStatus struct {
	ID            int
	CreatedAt     time.Time
	Status        string
	Carrier       string
	TrackingID    string
	TransactionID int `pg:",notnull"`
	Transaction   *Transaction
}
