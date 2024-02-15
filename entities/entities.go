package entities

type Wishlist struct {
	Id     uint
	UserId uint
}
type WishlistItems struct {
	Id         uint
	WishlistId uint
	Wishlist   `gorm:"foreignKey:WishlistId"`
	ProductId  uint
}
