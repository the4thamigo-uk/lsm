package lsm

import (
	"encoding/gob"
	"fmt"
	"github.com/the4thamigo-uk/multicloser"
	"os"
)

type (
	Database struct {
		path   string
		index  map[string]int64
		rdFile *os.File
		apFile *os.File
		mc     *multicloser.MultiCloser
	}
	record struct {
		Key string
		Val string
	}
)

func New(path string) (*Database, error) {
	var mc multicloser.MultiCloser

	fn := path + "/.db"

	apFile, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Could not open file for append '%s' : %w", fn, err)
	}
	mc.Defer(apFile.Close)

	rdFile, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("Could not open file for read '%s' : %w", fn, err)
	}
	mc.Defer(rdFile.Close)

	return &Database{
		path:   path,
		rdFile: rdFile,
		apFile: apFile,
		mc:     &mc,
	}, nil
}

func (d *Database) Set(key, val string) error {
	off, err := d.write(key, val)
	if err != nil {
		return err
	}
	d.addIndex(key, off)
	return nil
}

func (d *Database) Get(key string) (string, error) {
	off, ok := d.getIndex(key)
	if !ok {
		return "", fmt.Errorf("no record with key '%v'", key)
	}
	_, err := d.rdFile.Seek(off, 0)
	if err != nil {
		return "", err
	}
	dec := gob.NewDecoder(d.rdFile)
	var rec record
	err = dec.Decode(&rec)
	if err != nil {
		return "", fmt.Errorf("could not decode value with key '%v' : %w", key, err)
	}

	if rec.Key != key {
		return "", fmt.Errorf("record found with incorrect key '%v' != '%v' : %w", rec.Key, key, err)
	}
	return rec.Val, nil
}

func (d *Database) Close(del bool) error {
	err := d.mc.Close()
	if err != nil {
		return err
	}
	if del {
		return os.Remove(d.apFile.Name())
	}
	return nil
}

func (d *Database) addIndex(key string, off int64) {
	if d.index == nil {
		d.index = map[string]int64{}
	}
	d.index[key] = off
}

func (d *Database) getIndex(key string) (int64, bool) {
	if d.index == nil {
		return -1, false
	}
	off, ok := d.index[key]
	return off, ok
}
func (d *Database) write(key, val string) (int64, error) {
	rec := record{
		Key: key,
		Val: val,
	}
	st, err := d.apFile.Stat()
	if err != nil {
		return -1, fmt.Errorf("could not stat file : %w", err)
	}
	off := st.Size()
	enc := gob.NewEncoder(d.apFile)
	err = enc.Encode(rec)
	if err != nil {
		return -1, fmt.Errorf("could not encode value '%v' : %w", rec, err)
	}
	return off, nil
}
