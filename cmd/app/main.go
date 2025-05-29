package main

import (
	"app/internal/config"
	"context"
)

func main() {

	ctx := context.Background()

	cfg := config.MustRead()
}
