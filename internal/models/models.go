package models

type User struct {
	UserID   string
	UserName string
	Email    string
}

type Item struct {
	ItemID    string
	ItemName  string
	Category  string
	Price     float64
	SellerID  string
	Status    int
	CreatedAt string
}

type Order struct {
	OrderID    string
	ItemID     string
	ItemName   string
	BuyerID    string
	BuyerName  string
	CreatedAt  string
}
