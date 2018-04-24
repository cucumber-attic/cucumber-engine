package formatter

import (
	"encoding/json"
)

// FormatAsEventProtocol formats the events with the event protocol formatter
func FormatAsEventProtocol(opts *CommonOptions) {
	go func() {
		for ev := range opts.EventChannel {
			bytes, err := json.Marshal(ev)
			if err != nil {
				panic(err)
			}
			opts.Stream.Write(bytes)
		}
		opts.Stream.Close()
	}()
}
