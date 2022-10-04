package dictionary

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"time"
)

type Dictionary struct {
	db *badger.DB
}
type Entry struct {
	Word       string
	Definition string
	CreatedAt  time.Time
}

func (e Entry) String() string {
	created := e.CreatedAt.Format(time.Stamp)
	return fmt.Sprintf("%-10v\t - %-50v %-6v", e.Word, e.Definition, created)
}

func New(dir string) (*Dictionary, error) {
	opts := badger.DefaultOptions(dir)
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, errors.New("Error db")
	}
	dictionary := &Dictionary{db}
	return dictionary, nil
}

func (d *Dictionary) Close() {
	d.db.Close()
}
