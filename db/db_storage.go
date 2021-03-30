package bolt

import (
	"fmt"
	"time"
	"github.com/boltdb/bolt"


)

const (
	processedFiles = "processed_files"
)
type Strorage struct {
	db *bolt.DB
}

func NewStrorage() (*Strorage, error)  {

	db, err := bolt.Open("db.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return &Strorage{}, err
	}

	s:= &Strorage{db: db}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(processedFiles))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return s, nil
}

func (s * Strorage) Close() (error)  {
	return s.db.Close()
}

func (s *Strorage) AddProcessedFile(file string) (error)  {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(processedFiles))
		if b == nil {
			return fmt.Errorf("%s bucket doesntexists", processedFiles)
		}
		return b.Put([]byte(file), nil)
	})
}

func (s *Strorage) CheckFileIsProssed(file string) (prossed bool, err error)  {

	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(processedFiles))

		if b == nil {
			return fmt.Errorf("%s bucket doesntexists", processedFiles)
		}

		v := b.Get([]byte(file))

		prossed = v != nil

		return nil
	})

	return
}



