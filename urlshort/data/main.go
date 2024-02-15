package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func main() {
	_, err := SetupDB()
	if err != nil {
		panic(err)
	}
}

// Creates the db buckets and returns the db
func SetupDB() (*bolt.DB, error) {
	db, err := bolt.Open("paths.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("Opening db: %s", err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Paths"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte("/houston"), []byte("https://www.houstontx.gov/"))
		err = b.Put([]byte("/atlanta"), []byte("https://www.atlantaga.gov/"))
		err = b.Put([]byte("/nashville"), []byte("https://www.visitmusiccity.com/"))
		err = b.Put([]byte("/birmingham"), []byte("https://www.birminghamal.gov/"))
		err = b.Put([]byte("/austin"), []byte("https://www.austintexas.org/"))
		if err != nil {
			return fmt.Errorf("put item: %s", err)
		}
		return nil
	})
	return db, nil
}
