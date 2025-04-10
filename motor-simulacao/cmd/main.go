package main

import (
	"context"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/natanro/simlador-credito/motor-simulacao/datasource/db"
	"github.com/natanro/simlador-credito/motor-simulacao/datasource/rabbitmq"
	"github.com/natanro/simlador-credito/motor-simulacao/infra"
	"github.com/natanro/simlador-credito/motor-simulacao/interactor"
	"github.com/natanro/simlador-credito/motor-simulacao/transportlayer/amqp"
	"github.com/natanro/simlador-credito/motor-simulacao/transportlayer/rest"
)

var (
	// dbUser = os.Getenv("DB_USER")
	// dbPassword = os.Getenv("DB_PASSWORD")
	// dbHost = os.Getenv("DB_HOST")
	// dbPort = os.Getenv("DB_PORT")
	// dbName = os.Getenv("DB_NAME")

	// mongoURI = os.Getenv("MONGO_URI")

	dbUser     = "root"
	dbPassword = "rootpassword"
	dbHost     = "localhost"
	dbPort     = "3306"
	dbName     = "simlador_credito"

	mongoURI = "mongodb://root:rootpassword@mongodb:27017/simlador_credito?authSource=admin"
)

func main() {
	port, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	log.Println("Connecting to database...", dbUser, dbPassword, dbHost, dbPort, dbName)
	mysqlDB, err := infra.NewDatabaseConnection(&infra.DatabaseConfig{
		User:     dbUser,
		Password: dbPassword,
		Host:     dbHost,
		Port:     port,
		Name:     dbName,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	mongoDB, err := infra.NewMongoDBConnection(context.Background(), &infra.MongoDBConfig{
		URI: mongoURI,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	queue := infra.NewQueue("simulation queue", 10)

	simulationParamRepository := db.NewSimulationParamRepository(mongoDB)
	simulationRepository := db.NewSimulationRepository(mysqlDB)
	simulationQueueRepository := rabbitmq.NewSimulationQueue(queue)

	simulationRegister := interactor.NewSimulationRegister(simulationRepository, simulationQueueRepository)
	rateStrategy := interactor.NewRateStrategy(simulationParamRepository)
	simulationProcessor := interactor.NewSimulationProcessor(simulationRepository, rateStrategy)

	amqpHander := amqp.NewSimulationAmqpHandler(simulationProcessor)
	queue.RegisterObserver(amqpHander)

	handler := rest.NewMotorHandler(simulationRegister)
	router := echo.New()
	rest.RegisterHandlersWithBaseURL(router, handler, "/motor-simulacao")

	if err := router.Start("localhost:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
