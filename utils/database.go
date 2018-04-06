package utils

import (
  "github.com/boltdb/bolt"
)

type DatabaseObject struct {
  DbFile string
  Bucket string
  Db *bolt.DB
}

func CreateNewDatabaseObject(dbFile string, bucket string) DatabaseObject{
  db, err := bolt.Open(dbFile, 0600, nil)
  RequireNil(err)
  return DatabaseObject{dbFile, bucket, db}
}

func (databaseObject *DatabaseObject) Set(key []byte, blob []byte) {
  err := databaseObject.Db.Update(func(tx *bolt.Tx) error {
    b, err := tx.CreateBucket([]byte(databaseObject.Bucket))
    RequireNil(err)
    err = b.Put(key, blob)
    RequireNil(err)
    return nil;
  })
  RequireNil(err)
}

func (databaseObject *DatabaseObject) Get(key []byte) []byte {
  var ret []byte
  err := databaseObject.Db.View(func(tx *bolt.Tx) error{
    ret = tx.Bucket([]byte(databaseObject.Bucket)).Get(key)
    return nil;
  })
  RequireNil(err)
  return ret
}

func (databaseObject *DatabaseObject) SetLastHash(hash []byte) {
  databaseObject.Set([]byte("l"),hash)
}

func (databaseObject *DatabaseObject) GetLastHash() []byte {
  return databaseObject.Get([]byte("l"))
}
func (databaseObject *DatabaseObject) GetLastBlockData() []byte {
  lastHash := databaseObject.Get([]byte("l"))
  lastBlockData := databaseObject.Get(lastHash)
  return lastBlockData
}
