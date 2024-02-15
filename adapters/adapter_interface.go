package adapters

import "github.com/vishnusunil243/WishlistService/entities"

type AdapterInterface interface {
	CreateWishlist(req entities.Wishlist) error
	AddToWishlist(req entities.WishlistItems, userId int) error
	RemoveFromWishlist(productId, userId int) error
	GetWishlistItem(productId, userId int) (entities.WishlistItems, error)
	GetAllWishlistItems(userId int) ([]entities.WishlistItems, error)
}
