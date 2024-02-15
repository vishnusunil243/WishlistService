package initializer

import (
	"github.com/vishnusunil243/WishlistService/adapters"
	"github.com/vishnusunil243/WishlistService/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.WishlistService {
	adapter := adapters.NewWishlistAdapter(db)
	service := service.NewWishlistService(adapter)
	return service
}
