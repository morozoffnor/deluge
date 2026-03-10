package main

import (
	client "github.com/morozoffnor/deluge"
)

func main() {
	deluge, err := client.New("https://your-deluge-url", "pass", "admin", "admin")
	if err != nil {
		panic(err)
	}
	_, err = deluge.Login()
	if err != nil {
		panic(err)
	}

	_, err = deluge.AddTorrentFile("/path/to/file.torrent", client.TorrentOptions{
		DownloadLocation: "/mnt/media_root/films",
	})
	if err != nil {
		panic(err)
	}

	_, err = deluge.AddTorrentMagnet("magnet:?xt=urn:btih:08ada5a7a6183aae1e09d831df6748d566095a10", client.TorrentOptions{
		DownloadLocation: "/mnt/media_root/shows",
	})
	if err != nil {
		panic(err)
	}
}
