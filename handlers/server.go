package handlers

import (
	"BookManager/authenticate"
	"BookManager/db"
	"github.com/sirupsen/logrus"
)

type BookManagerServer struct {
	Db           *db.GormDB
	Logger       *logrus.Logger
	Authenticate *authenticate.Auth
}
