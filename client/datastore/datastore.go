package datastore

import "errors"

const (
	MEMCACHED = iota
)

type Datastore interface {
	Ping() DataStoreStatus
}

type DataStoreStatus struct {
	Code   int
	Status string
	Error  error
}

func NewDataStore(dstoreType int, conn string) (Datastore, error) {
	switch dstoreType {
	case MEMCACHED:
		return NewMemcacheDatastore(conn)
	default:
		return nil, errors.New("Datastore not registered")
	}
}
