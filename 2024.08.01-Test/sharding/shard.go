// Package naive provides a "naive" algorithm for solving the test problem
package shard

import (
	"crypto/sha1"
	"sync"
	"unsafe"
)

type Shard struct {
	sync.RWMutex                 // Compose from sync.RWMutex
	m            map[uint32]bool // m contains the shard's data
}

type ShardedMap []*Shard // ShardedMap is a *Shards slice

func NewShardedMap(nshards int) ShardedMap {
	shards := make([]*Shard, nshards) // Initialize a *Shards slice

	for i := 0; i < nshards; i++ {
		shard := make(map[uint32]bool)
		shards[i] = &Shard{m: shard}
	}

	return shards // A ShardedMap IS a *Shards slice!
}

func (m ShardedMap) getShardIndex(key uint32) int {
	a := (*[4]byte)(unsafe.Pointer(&key))[:]
	checksum := sha1.Sum(a)   // Use Sum from "crypto/sha1"
	hash := int(checksum[17]) // Pick an arbitrary byte as the hash
	return hash % len(m)      // Mod by len(m) to get index
}

func (m ShardedMap) getShared(key uint32) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

func (m ShardedMap) Get(key uint32) bool {
	shard := m.getShared(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardedMap) Set(key uint32, value bool) {
	shard := m.getShared(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

func (m ShardedMap) Len() int {
	res := 0

	for _, shard := range m {
		res += len(shard.m)
	}

	return res
}
