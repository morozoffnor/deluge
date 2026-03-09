package main

import (
	"context"

	"github.com/morozoffnor/deluge/pkg"
)

func main() {
	deluge, err := deluge.New("https://torrent.kitburg.ru", "134562", "admin", "134562Morozoff")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	deluge.Start(ctx)
}
