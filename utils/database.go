package utils

import (
  "log"
  "github.com/boltdb/bolt"
)

type DatabaseObject struct {
  DbFile string
  Bucket string
  db *bolt.DB
}

func CreateNewDatabaseObject(dbFile string, bucket string, db *bolt.DB) DatabaseObject{
  db, err := bolt.Open(dbFile, 0600, nil)
  RequireNil(err)
  return DatabaseObject{dbFile, bucket, db}
}

func (databaseObject *DatabaseObject) Update(key []byte, blob []byte) {
  err := databaseObject.db.Update(func(tx *bolt.Tx) error {
    b, err := tx.CreateBucket([]byte(databaseObject.Bucket))
    RequireNil(err)
    err = b.Put(key, blob)
    RequireNil(err)
    return nil;
  })
  RequireNil(err)
}
