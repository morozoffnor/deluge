package main

import (
	"context"
	"deluge/config"
	"deluge/server"
	"deluge/tgbot"
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.New()
	b, err := tgbot.New(cfg)
	if err != nil {
		panic(err)
	}

	tgbot.InitReplyKeyboard(b)

	s := server.New(cfg, b)
	g.Go(func() error {
		fmt.Println("Server started on port", cfg.APIPort, "\nThis server is meant to run on localhost only "+
			"besides deluge instance. Use bundled .sh scripts in deluge UI")
		err := s.ListenAndServe()
		println("server stopped", err.Error())
		return err
	})

	b.Start(ctx)

}
