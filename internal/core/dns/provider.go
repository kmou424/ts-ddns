package dns

type IProvider interface {
	Init(params map[string]string) error
	GetRecords() ([]Record, error)
	DeleteRecord(record Record) error
	CreateRecord(record Record) error
	UpdateRecord(record Record) error
}
