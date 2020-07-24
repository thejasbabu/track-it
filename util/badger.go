package util

import "github.com/dgraph-io/badger"

// Badger interface for abstracting badger operations
type Badger interface {
	Update(key, value []byte) error
	Read(accessFn func([]byte) error) error
	Delete(key []byte) error
}

// BadgerDB for BadgerOperations
type BadgerDB struct {
	db *badger.DB
}

// Open opens the badgerDB
func Open(path string) (BadgerDB, error) {
	opt := badger.DefaultOptions(path)
	db, err := badger.Open(opt)
	return BadgerDB{db: db}, err
}

// Update updates the badgerDB with key and value
func (b BadgerDB) Update(key, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
}

// Update updates the badgerDB with key and value
func (b BadgerDB) Delete(key []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
}

func (b BadgerDB) Read(accessFn func([]byte) error) error {
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		iter := txn.NewIterator(opts)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			err := item.Value(accessFn)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
