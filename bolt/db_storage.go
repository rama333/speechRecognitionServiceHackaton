package bolt1

import (
	"time"
	"github.com/boltdb/bolt"
)

const (
	processedFiles = "processed_files"
)
type Strorage struct {
	
}

func NewStrorage()  {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})


}
