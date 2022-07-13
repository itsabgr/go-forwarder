package handy

import "time"

func Stop(t *time.Timer) {
	for !t.Stop() {
		_ = <-t.C
	}
}
