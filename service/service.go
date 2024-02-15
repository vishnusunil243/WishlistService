package service

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/vishnusunil243/WishlistService/adapters"
	"github.com/vishnusunil243/WishlistService/entities"
	"github.com/vishnusunil243/proto-files/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	Tracer        opentracing.Tracer
	ProductClient pb.ProductServiceClient
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type WishlistService struct {
	Adapter adapters.AdapterInterface
	pb.UnimplementedWishlistServiceServer
}

func NewWishlistService(adapter adapters.AdapterInterface) *WishlistService {
	return &WishlistService{
		Adapter: adapter,
	}
}
func (wishlist *WishlistService) CreateWishlist(ctx context.Context, req *pb.CreateWishlistRequest) (*emptypb.Empty, error) {
	span := Tracer.StartSpan("create Wishlist grpc")
	defer span.Finish()
	err := wishlist.Adapter.CreateWishlist(entities.Wishlist{
		UserId: uint(req.UserId),
	})
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
func (wishlist *WishlistService) AddToWishlist(ctx context.Context, req *pb.AddToWishlistRequest) (*pb.CreateWishlistRequest, error) {
	product, err := ProductClient.GetProduct(context.TODO(), &pb.GetProductById{
		Id: int32(req.ProductId),
	})
	if err != nil {
		return &pb.CreateWishlistRequest{}, fmt.Errorf("error fetching product from product service")
	}
	if product.Name == "" {
		return &pb.CreateWishlistRequest{}, fmt.Errorf("product with the given id is not found")
	}
	wishlistItem, err := wishlist.Adapter.GetWishlistItem(int(req.ProductId), int(req.UserId))
	if err != nil {
		return &pb.CreateWishlistRequest{}, err
	}
	if wishlistItem.ProductId != 0 {
		return &pb.CreateWishlistRequest{}, fmt.Errorf("item already present in wishlist")
	}
	err = wishlist.Adapter.AddToWishlist(entities.WishlistItems{
		ProductId: uint(req.ProductId),
	}, int(req.UserId))
	if err != nil {
		return &pb.CreateWishlistRequest{}, err
	}
	res := &pb.CreateWishlistRequest{
		UserId: req.UserId,
	}
	return res, nil
}
func (wishlist *WishlistService) RemoveFromWishlist(ctx context.Context, req *pb.AddToWishlistRequest) (*pb.CreateWishlistRequest, error) {
	item, err := wishlist.Adapter.GetWishlistItem(int(req.ProductId), int(req.UserId))
	if err != nil {
		return &pb.CreateWishlistRequest{}, err
	}
	if item.ProductId == 0 {
		return &pb.CreateWishlistRequest{}, fmt.Errorf("product not found in wishlist")
	}
	err = wishlist.Adapter.RemoveFromWishlist(int(req.ProductId), int(req.UserId))
	if err != nil {
		return &pb.CreateWishlistRequest{}, err
	}
	return &pb.CreateWishlistRequest{UserId: req.UserId}, nil
}
func (wishlist *WishlistService) GetAllWishlistItems(req *pb.CreateWishlistRequest, srv pb.WishlistService_GetAllWishlistItemsServer) error {
	wishlistItems, err := wishlist.Adapter.GetAllWishlistItems(int(req.UserId))
	if err != nil {
		return err
	}
	for _, item := range wishlistItems {
		res := &pb.GetAllWishlistResponse{
			ProductId: uint32(item.ProductId),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Println("check called")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
