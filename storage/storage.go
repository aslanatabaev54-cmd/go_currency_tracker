package storage

type DataSaver interface {
	SaveCurrency(code string, rate float64) error
}
