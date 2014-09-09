package app

import "crypto/sha256"
import "encoding/hex"
import "time"
import "appengine"
import "appengine/datastore"

type entity struct {
  Time time.Time `datastore:"time,noindex"`
}

func stringID(key string) string {
  sha := sha256.New()
  sha.Write([]byte(key))
  sum := sha.Sum(nil)
  return hex.EncodeToString(sum)
}

func makeDatastoreKey(c appengine.Context, key string) *datastore.Key {
  return datastore.NewKey(c, "timestamp", stringID(key), 0, nil)
}

func getTimestamp(c appengine.Context, key string) (*time.Time, error) {
  datastoreKey := makeDatastoreKey(c, key)
  entity := entity{}
  err := datastore.Get(c, datastoreKey, &entity)
  if err == datastore.ErrNoSuchEntity {
    return nil, nil
  } else if err != nil {
    return nil, err
  } else {
    return &entity.Time, nil
  }
}

func putTimestamp(c appengine.Context, key string) (*time.Time, error) {
  datastoreKey := makeDatastoreKey(c, key)
  entity := entity{time.Now()}
  _, e := datastore.Put(c, datastoreKey, &entity)
  if e == nil {
    return &entity.Time, nil
  } else {
    return nil, e
  }
}
