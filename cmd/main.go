package main

import (
	"flag"
	"fmt"
	"github.com/bdarge/api/out/customer"
	"github.com/bdarge/api/out/profile"
	"github.com/bdarge/api/out/transaction"
	"github.com/bdarge/api/out/transactionItem"
	"github.com/bdarge/api/pkg/config"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"github.com/bdarge/api/pkg/services"
	"github.com/bdarge/api/pkg/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"time"
)

var (
	sleep  = flag.Duration("sleep", time.Second*5, "duration between changes in health")
	system = "" // empty string represents the health of the system
)

func main() {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "dev"
	}

	conf, err := config.LoadConfig(environment)
	if err != nil {
		panic(fmt.Errorf("Failed to load configuration: %s. ", err))
	}

	// init env object
	env := util.NewEnv()
	migrateData := env.GetBool("MIGRATOR")
	log.Printf("migrateData => %t", migrateData)

	handler := db.Init(conf.DSN)

	if migrateData {
		if err := util.Migrate(conf, handler); err != nil {
			log.Fatalln(err)
		}
		log.Println("Successfully migrated database.")
		os.Exit(0)
	}

	lis, err := net.Listen("tcp", conf.ServerPort)

	if err != nil {
		log.Fatalf("Listing on port %s has failed: %v", conf.ServerPort, err)
	}

	grpcServer := grpc.NewServer()

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)

	transactionServer := services.Server{
		H: handler,
	}
	transaction.RegisterTransactionServiceServer(grpcServer, &transactionServer)

	transactionItemServer := services.TransactionItemServer{
		H: handler,
	}
	transactionItem.RegisterTransactionItemServiceServer(grpcServer, &transactionItemServer)

	customerServer := services.CustomerServer{
		H: handler,
	}
	customer.RegisterCustomerServiceServer(grpcServer, &customerServer)

	profileServer := services.ProfileServer{
		H: handler,
	}
	profile.RegisterProfileServiceServer(grpcServer, &profileServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

	fmt.Printf("api service is listening on %s", conf.ServerPort)

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthcheck.SetServingStatus(system, next)
			err = isDbConnectionWorks(profileServer.H.DB)
			if err != nil {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}
			time.Sleep(*sleep)
		}
	}()
}

func isDbConnectionWorks(DB *gorm.DB) error {
	return DB.First(&models.Account{}).Error
}
