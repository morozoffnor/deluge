# deluge

A Go client library for the [Deluge](https://deluge-torrent.org/) BitTorrent client JSON-RPC API. Zero external dependencies.

## Installation

```bash
go get github.com/morozoffnor/deluge
```

## Usage

```go
package main

import (
	client "github.com/morozoffnor/deluge"
)

func main() {
	// Create a new client (with optional HTTP Basic Auth)
	deluge, err := client.New("https://your-deluge-url", "password", "basicUser", "basicPass")
	if err != nil {
		panic(err)
	}

	// Authenticate with Deluge
	_, err = deluge.Login()
	if err != nil {
		panic(err)
	}

	// Add a torrent from a magnet link
	torrentID, err := deluge.AddTorrentMagnet("magnet:?xt=urn:btih:...", client.TorrentOptions{
		DownloadLocation: "/mnt/media_root/shows",
	})

	// Add a torrent from a .torrent file
	torrentID, err = deluge.AddTorrentFile("/path/to/file.torrent", client.TorrentOptions{
		DownloadLocation: "/mnt/media_root/films",
	})
}
```

## Torrent Options

| Option | Type | Description |
|---|---|---|
| `DownloadLocation` | `string` | Where to save files |
| `MoveCompleted` | `bool` | Move files on completion |
| `MoveCompletedPath` | `string` | Where to move completed files |
| `MaxDownloadSpeed` | `float64` | KB/s, -1 for unlimited |
| `MaxUploadSpeed` | `float64` | KB/s, -1 for unlimited |
| `MaxConnections` | `int` | -1 for unlimited |
| `MaxUploadSlots` | `int` | -1 for unlimited |
| `Paused` | `bool` | Add in paused state |
| `PrioritizeFirstLast` | `bool` | Prioritize first/last pieces |
| `SeedMode` | `bool` | Skip hash check |
| `SuperSeeding` | `bool` | Enable super seeding |
