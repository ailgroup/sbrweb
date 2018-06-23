package datastore

type MemcacheDatastore struct {
}

func NewMemcacheDatastore(address string) (*MemcacheDatastore, error) {
	return &MemcacheDatastore{}, nil
}

func (m *MemcacheDatastore) Ping() DataStoreStatus {
	return DataStoreStatus{
		Code:   201,
		Status: "ok",
		Error:  nil,
	}
}
