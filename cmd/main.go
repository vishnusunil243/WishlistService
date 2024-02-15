package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	jaegar "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/vishnusunil243/WishlistService/db"
	"github.com/vishnusunil243/WishlistService/initializer"
	"github.com/vishnusunil243/WishlistService/service"
	servicediscovery "github.com/vishnusunil243/WishlistService/service_discovery"
	"github.com/vishnusunil243/proto-files/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err.Error())
	}
	productConn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		productConn.Close()
	}()
	productRes := pb.NewProductServiceClient(productConn)
	service.ProductClient = productRes
	addr := os.Getenv("DB_KEY")
	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	services := initializer.Initializer(DB)
	server := grpc.NewServer()
	pb.RegisterWishlistServiceServer(server, services)
	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal("failed to listen on port 8085")
	}
	log.Printf("wishlist server listening on port 8085")
	go func() {
		time.Sleep(5 * time.Second)
		servicediscovery.RegisterService()
	}()
	healthService := &service.HealthChecker{}
	grpc_health_v1.RegisterHealthServer(server, healthService)
	tracer, closer := InitTracer()
	defer closer.Close()
	service.RetrieveTracer(tracer)
	if err = server.Serve(listener); err != nil {
		log.Fatal("failed to listen on port 8085")
	}
}
func InitTracer() (tracer opentracing.Tracer, closer io.Closer) {
	jaegarEndpoint := "http://localhost:14628/api/tracer"
	cfg := &config.Configuration{
		ServiceName: "wishlist-service",
		Sampler: &config.SamplerConfig{
			Type:  jaegar.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: jaegarEndpoint,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("updated")
	return
}
