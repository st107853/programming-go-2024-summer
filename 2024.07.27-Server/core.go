package main

import (
	"errors"
	"strconv"
	"sync"
)

var ErrorNoSuchId = errors.New("no such id")

type album struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var store = struct {
	sync.RWMutex
	m      map[uint64]*album
	lastId uint64
}{m: make(map[uint64]*album)}

func Post(title, artist, price string) (uint64, error) {
	n, _ := strconv.ParseFloat(price, 64)
	a := album{Title: title, Artist: artist, Price: n}

	store.Lock()
	store.lastId++
	store.m[store.lastId] = &a
	store.Unlock()

	return store.lastId, nil
}

func Get(id string) (*album, error) {
	idx, _ := strconv.ParseUint(id, 10, 64)

	store.RLock()
	value, ok := store.m[idx]
	store.RUnlock()

	if !ok {
		return nil, ErrorNoSuchId
	}

	return value, nil
}

func Delete(id uint64) error {
	store.Lock()
	delete(store.m, id)
	store.Unlock()

	return nil
}
