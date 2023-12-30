package converter

import "github.com/lspaccatrosi16/spac/lib/converter/time"

func list() []Convertable {
	return []Convertable{
		&time.ConvTime{},
	}
}
