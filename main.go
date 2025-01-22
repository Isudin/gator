package main

import (
	"fmt"

	"github.com/Isudin/gator/internal/config"
)

var cfg config.Config

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.DbUrl)

}
