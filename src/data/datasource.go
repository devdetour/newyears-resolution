package data

// Interface for all data sources
type IDataSource interface {
	Connect() bool
	Measure() struct{}
}

type StravaDataSource struct {
}

func (d *StravaDataSource) Connect() bool {
	return true
}

func (d *StravaDataSource) Measure() struct{} {
	s := struct{}{}
	return s
}
