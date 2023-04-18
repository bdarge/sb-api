package main

import (
	"fmt"
	"github.com/bdarge/api/out/customer"
	"github.com/bdarge/api/out/profile"
	"github.com/bdarge/api/out/transaction"
	"github.com/bdarge/api/pkg/config"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/services"
	"github.com/bdarge/api/pkg/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
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
		log.Fatalln("Failed to listing:", err)
	}

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("api service on", conf.ServerPort)

	transactionServer := services.Server{
		H: handler,
	}

	grpcServer := grpc.NewServer()

	transaction.RegisterTransactionServiceServer(grpcServer, &transactionServer)

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
}
