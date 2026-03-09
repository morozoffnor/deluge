package deluge

import "context"

const (
	update_type_recieved_magnet = iota
	update_type_recieved_file
)

func (d *DelugeClient) processUpdate(ctx context.Context, upd *DelugeUpdate) {

}
