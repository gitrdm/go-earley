package bitmatrix

import (
	"github.com/patrickhuber/go-earley/internal/bitarray"
	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
)

type BitMatrix struct {
	matrix *bitarray.BitArray
	bits   int
}

func New(bits int) *BitMatrix {
	return &BitMatrix{
		matrix: bitarray.New(bits * bits),
		bits:   bits,
	}
}

func (bm *BitMatrix) Get(x int, y int) (bool, error) {
	return bm.get(x, y).Deconstruct()
}

func (bm *BitMatrix) get(x int, y int) types.Result[bool] {
	index := x + y*bm.bits
	return result.New(bm.matrix.Get(index))
}

func (bm *BitMatrix) Set(x int, y int, value bool) error {
	_, err := bm.set(x, y, value).Deconstruct()
	return err
}

func (bm *BitMatrix) set(x int, y int, value bool) types.Result[struct{}] {
	index := x + y*bm.bits
	return result.New(struct{}{}, bm.matrix.Set(index, value))
}

func (bm *BitMatrix) Len() int {
	return bm.bits
}

func TransitiveClosure(bm *BitMatrix) (*BitMatrix, error) {
	return transitiveClosure(bm).Deconstruct()
}

func transitiveClosure(bm *BitMatrix) (res types.Result[*BitMatrix]) {
	defer handle.Error(&res)
	clone := Clone(bm)

	for k := 0; k < clone.Len(); k++ {
		for j := 0; j < clone.Len(); j++ {
			for i := 0; i < clone.Len(); i++ {
				ij := clone.get(i, j).Unwrap()
				ik := clone.get(i, k).Unwrap()
				kj := clone.get(k, j).Unwrap()
				clone.set(i, j, ij || ik && kj).Unwrap()
			}
		}
	}
	return result.Ok(clone)
}

func Clone(bm *BitMatrix) *BitMatrix {
	clone := bitarray.Clone(bm.matrix)
	return &BitMatrix{
		matrix: clone,
		bits:   bm.bits,
	}
}
