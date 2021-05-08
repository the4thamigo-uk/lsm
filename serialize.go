package lsm

import (
	"encoding/gob"
	"github.com/the4thamigo-uk/lsm/sortedmap"
	"io"
)

func WriteMap(w io.Writer, m sortedmap.Map) error {
	enc := gob.NewEncoder(w)
	it := m.Iter("")
	for {
		key, val, ok := it()
		if !ok {
			break
		}
		if err := enc.Encode(key); err != nil {
			return err
		}
		if err := enc.Encode(&val); err != nil {
			return err
		}
	}
	return nil
}

func ReadMap(r io.Reader) (sortedmap.Map, error) {
	var m sortedmap.Map
	dec := gob.NewDecoder(r)

	for {
		var key string
		if err := dec.Decode(&key); err != nil {
			if err == io.EOF {
				return m, nil
			}
			return m, err
		}
		var val interface{}
		if err := dec.Decode(&val); err != nil {
			return m, err
		}
		m.Add(key, val)
	}
}
