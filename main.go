package main

import (
	"context"
	"fmt"

	"github.com/jrh3k5/moo4plex/ui"
)

func main() {
	app := ui.NewApp()
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		fmt.Printf("ERROR: Failed to run application: %v\n", err)
	}
}
