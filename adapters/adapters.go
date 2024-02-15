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
	if err := wishlist.DB.Raw(getWishlistId, userId).Scan(&wishlistId).Error; err != nil {
		return err
	}
	insertQuery := `INSERT INTO wishlist_items (product_id,wishlist_id) VALUES ($1,$2)`
	if err := wishlist.DB.Exec(insertQuery, req.ProductId, wishlistId).Error; err != nil {
		return err
	}
	return nil

}
func (wishlist *WishlistAdapter) RemoveFromWishlist(productId, userId int) error {
	var wishlistId int
	tx := wishlist.DB.Begin()
	getId := `SELECT id FROM wishlists WHERE user_id=?`
	if err := tx.Raw(getId, userId).Scan(&wishlistId).Error; err != nil {
		return err
	}
	removeItems := `DELETE FROM wishlist_items WHERE wishlist_id=$1 AND product_id=$2`
	if err := tx.Exec(removeItems, wishlistId, productId).Error; err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
func (wishlist *WishlistAdapter) GetWishlistItem(productId, userId int) (entities.WishlistItems, error) {
	var wishlistId int
	getId := `SELECT id FROM wishlists WHERE user_id=?`
	if err := wishlist.DB.Raw(getId, userId).Scan(&wishlistId).Error; err != nil {
		return entities.WishlistItems{}, err
	}
	var res entities.WishlistItems
	getItem := `SELECT * FROM wishlist_items WHERE wishlist_id=$1 AND product_id=$2`
	if err := wishlist.DB.Raw(getItem, wishlistId, productId).Scan(&res).Error; err != nil {
		return entities.WishlistItems{}, err
	}
	return res, nil
}
func (wishlist *WishlistAdapter) GetAllWishlistItems(userId int) ([]entities.WishlistItems, error) {
	var wishlistId int
	getId := `SELECT id FROM wishlists WHERE user_id=?`
	if err := wishlist.DB.Raw(getId, userId).Scan(&wishlistId).Error; err != nil {
		return []entities.WishlistItems{}, err
	}
	var res []entities.WishlistItems
	getAll := `SELECT * FROM wishlist_items WHERE wishlist_id=?`
	if err := wishlist.DB.Raw(getAll, wishlistId).Scan(&res).Error; err != nil {
		return []entities.WishlistItems{}, err
	}
	return res, nil
}
