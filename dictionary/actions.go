package dictionary

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"strings"
	"time"
)

func (d *Dictionary) AddEntry(word, definition string) (err error) {
	entry := Entry{strings.ToTitle(word), definition, time.Now()}
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(entry)
	if err != nil {
		return err
	}
	err = d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(word), buffer.Bytes())
	})
	if err != nil {
		return err
	}
	return
}

func (d *Dictionary) Get(key string) (Entry, error) {
	var entry Entry
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		data, err := getEntry(item)
		if err != nil {
			return err
		}
		entry = data
		return nil

	})
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}

func getEntry(item *badger.Item) (Entry, error) {
	var buffer bytes.Buffer
	var entry Entry
	err := item.Value(func(val []byte) error {
		_, err := buffer.Write(val)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Entry{}, err
	}
	dec := gob.NewDecoder(&buffer)
	dec.Decode(&entry)
	return entry, nil
}

func (d *Dictionary) Delete(key string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})

}

func (d *Dictionary) List() ([]Entry, error) {
	var e []Entry

	err := d.db.View(func(txn *badger.Txn) error {
		var opts badger.IteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			entry, err := getEntry(item)
			if err != nil {
				fmt.Println("ERROR ", err)
			}
			e = append(e, entry)
		}
		return nil
	})
	if err != nil {
		return e, err
	}
	return e, nil
}
