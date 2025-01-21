package dns

import (
	"github.com/kmou424/ts-ddns/pkgs/typed"
)

const (
	RecordTypeA    = "A"
	RecordTypeAAAA = "AAAA"
)

type Record struct {
	ID      string
	Type    string
	Domain  string
	IP      string
	TTL     int
	Comment string
	Extra   typed.Map[string]
}

func NewEmptyRecord() Record {
	return Record{
		Extra: typed.NewMap[string](),
	}
}
