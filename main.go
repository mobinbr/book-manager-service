package main

import (
	"BookManager/authenticate"
	"BookManager/config"
	"BookManager/db"
	"BookManager/handlers"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	// Load the config
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err.Error())
	}

	// Setup the logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// Create new instance of GormDB containing Configuration and DB
	gormDB, err := db.NewGormDB(cfg)
	if err != nil {
		logger.WithError(err).Fatal("error in connecting to the postgres database")
	}
	logger.Info("connected to the book_manager_db database...")

	// Create database schemas for models
	err = gormDB.CreateSchemas()
	if err != nil {
		logger.WithError(err).Fatal("error in database migration")
	}
	logger.Infoln("migrate tables and models successfully...")

	// Create authenticate
	auth, err := authenticate.NewAuth(gormDB, 10, logger)
	if err != nil {
		logger.WithError(err).Fatal("can not create the authenticate instance")
	}

	bookManagerServer := handlers.BookManagerServer{Db: gormDB, Logger: logger, Authenticate: auth}

	http.HandleFunc("/api/v1/auth/signup", bookManagerServer.HandleSignup)
	http.HandleFunc("/api/v1/auth/login", bookManagerServer.HandleLogin)
	http.HandleFunc("/profile", bookManagerServer.HandleProfile)
	http.HandleFunc("/api/v1/books", bookManagerServer.HandleCreateGetAllBooks)
	http.HandleFunc("/api/v1/books/", bookManagerServer.HandleUpdateRetrieveDeleteBook)

	logger.WithError(http.ListenAndServe(":8080", nil)).Fatalln("can not setup the server")
}
