package servicediscovery

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	port      = 8085
	serviceId = "wishlist-service"
)

func RegisterService() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalf(err.Error())
	}
	addr := "localhost"
	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    "wishlist-server",
		Port:    port,
		Address: addr,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", addr, port, serviceId),
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	regiErr := consul.Agent().ServiceRegister(registration)
	if regiErr != nil {
		log.Fatal("failed to register service")
	} else {
		log.Printf("sucessfully registered service on %s:%d", addr, port)
	}

}
