package database

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func CreateSession() {
	Store = session.New()
}
