package adapters

import (
	"github.com/vishnusunil243/WishlistService/entities"
	"gorm.io/gorm"
)

type WishlistAdapter struct {
	DB *gorm.DB
}

func NewWishlistAdapter(db *gorm.DB) *WishlistAdapter {
	return &WishlistAdapter{
		DB: db,
	}
}
func (wishlist *WishlistAdapter) CreateWishlist(req entities.Wishlist) error {
	query := `INSERT INTO wishlists (user_id) VALUES ($1)`
	if err := wishlist.DB.Exec(query, req.UserId).Error; err != nil {
		return err
	}
	return nil
}
func (wishlist *WishlistAdapter) AddToWishlist(req entities.WishlistItems, userId int) error {
	var wishlistId int
	getWishlistId := `SELECT id FROM wishlists WHERE user_id=?`
	if err := wishlist.DB.Raw(getWishlistId, req.UserId).Scan(&wishlistId).Error; err != nil {
		return err
	}
	insertQuery := `INSERT INTO wishlist_items (product_id) VALUES ($1)`
	if err := wishlist.DB.Exec(insertQuery, req.ProductId).Error; err != nil {
		return err
	}
	return nil

}
