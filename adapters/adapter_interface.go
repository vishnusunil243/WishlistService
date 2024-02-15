package adapters

import "github.com/vishnusunil243/WishlistService/entities"

type AdapterInterface interface {
	CreateWishlist(req entities.Wishlist) error
	AddToWishlist(req entities.WishlistItems, userId int) error
}
