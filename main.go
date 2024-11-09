package main

import (
	"fmt"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/config"
)

func main() {
  cfg := config.Read()

  cfg.SetUser("lane")

  cfg = config.Read()

  fmt.Println(cfg)
}
