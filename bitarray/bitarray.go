package bitarray

import (
	"fmt"
	"math"
)

type BitArray struct {
	array []uint64
	bits  int
}

func New(bits int) *BitArray {
	if bits == 0 {
		bits = 1
	}
	size := bits / 64
	if bits%64 != 0 {
		size += 1
	}
	array := make([]uint64, size)
	return &BitArray{
		array: array,
		bits:  bits,
	}
}

func (ba *BitArray) Get(index int) (bool, error) {
	if err := ba.check(index); err != nil {
		return false, err
	}

	bucket := index / ba.bits
	bit := index % 64
	var mask uint64 = 1 << bit
	return (ba.array[bucket] & mask) != 0, nil
}

func (ba *BitArray) Set(index int, value bool) error {
	if err := ba.check(index); err != nil {
		return err
	}
	bucket := index / ba.bits
	bit := index % 64

	var i uint64 = 0
	if value {
		i = 1
	}
	var mask uint64 = i << bit
	ba.array[bucket] |= mask
	return nil
}

func (ba *BitArray) SetAll(b bool) {
	var v uint64 = 0
	if b {
		v = math.MaxUint32
	}
	for i := 0; i < len(ba.array); i++ {
		ba.array[i] = v
	}
}

func (ba *BitArray) Len() int {
	return ba.bits
}

func Clone(ba *BitArray) *BitArray {
	array := make([]uint64, len(ba.array))
	copy(array, ba.array)
	return &BitArray{
		array: array,
		bits:  ba.bits,
	}
}

func (ba *BitArray) check(index int) error {
	if index < 0 {
		return fmt.Errorf("index is less than zero")
	}
	if index >= ba.bits {
		return fmt.Errorf("index is greater than size %d", ba.bits)
	}
	return nil
}
