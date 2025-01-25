package main

import (
	"github.com/Isudin/gator/internal/config"
	"github.com/Isudin/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
