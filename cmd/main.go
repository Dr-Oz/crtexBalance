package main

import (
	"context"
	c "crtexBalance/internal/config"
	"crtexBalance/internal/handler"
	"crtexBalance/internal/repository"
	"crtexBalance/internal/service"
	"crtexBalance/mq"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title crtexBalance API
// @version 1.0
// @description Microservice for working with user balance

// @host localhost:8081
// @BasePath /

func main() {
	var services *service.Service
	var db *sql.DB
	var err error
	var conf *c.Config
	var path string
	var defaultPath string = "./configs/config.yaml"

	path = defaultPath

	flag.StringVar(&path, "config", "./configs/config.yaml", "example -config ./configs/config.yaml")
	migrationup := flag.Bool("migrationup", false, "use migrationup to perform migrationup")
	migrationdown := flag.Bool("migrationdown", false, "use migrationdown to perform migrationdown")

	flag.Parse()

	conf, err = c.GetConfig(path)
	if err != nil {
		log.Printf("%s, use default config '%s'", err, defaultPath)
		conf, err = c.GetConfig(defaultPath)
		if err != nil {
			log.Println(err)
			return
		}
	}

	if err = repository.Migration(conf, *migrationup, *migrationdown); err != nil {
		log.Println(err)
	}

	if db, err = repository.Connect(conf); err != nil {
		log.Println(err)
		return
	}

	// Initialize RabbitMQ
	rabbitmq, err := mq.NewRabbitMQ("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	repos := repository.NewRepository(db)
	services = service.NewService(repos, conf, db, rabbitmq)
	handlers := handler.NewHandler(services)

	server := new(Server)
	server.conf = conf

	replenishmentConsumer, err := mq.NewReplenishmentConsumer("amqp://guest:guest@rabbitmq:5672/") // Update with your RabbitMQ URL
	if err != nil {
		log.Fatalf("Failed to create replenishment consumer: %v", err)
	}

	go func() {
		if err := replenishmentConsumer.Consume(); err != nil {
			log.Fatalf("Failed to consume replenishment messages: %v", err)
		}
	}()

	go func() {
		if err := server.Run(conf.Port, handlers.Init()); err != nil {
			log.Fatalf("ошибка при запуске http сервера: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("сервер останавливается")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.ContexTimeout)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("произошла ошибка при выключении сервера: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf("произошла ошибка при закрытии соединения с БД: %s", err.Error())
	}

}
