package db

type db struct {
	driver string
}

type Db interface {
	StoreWsnSensorTelementyData(data []byte) error
}

func New() (Db, error) {
	db := db{
		driver: "hi",
	}
	return &db, nil
}

func (db *db) StoreWsnSensorTelementyData(data []byte) error {
	return nil
}
