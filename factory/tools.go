package factory

import (
	"context"
	"reflect"
	"sync/atomic"
)

type (
	Stacks   []*int64
	pipeline struct {
		stacks Stacks
		parent Args
	}
	attrGenerator struct {
		genFunc  func(Args) (any, error)
		key      string
		value    any
		isNil    bool
		isFilled bool
	}
	argsStruct struct {
		ctx context.Context
		rv  *reflect.Value
		pl  *pipeline
	}
)

func newPipeline(size int) *pipeline {
	return &pipeline{stacks: make(Stacks, size)}
}

// Instance returns a object to which the generator declared just before is applied
func (args *argsStruct) Instance() any {
	return args.rv.Interface()
}

// Parent returns a parent argument if current factory is a subfactory of parent
func (args *argsStruct) Parent() Args {
	if args.pl == nil {
		return nil
	}
	return args.pl.parent
}

func (args *argsStruct) Context() context.Context {
	return args.ctx
}

func (args *argsStruct) UpdateContext(ctx context.Context) {
	args.ctx = ctx
}

func (st *Stacks) Size(idx int) int64 {
	return *(*st)[idx]
}

// Set method is not goroutine safe.
func (st *Stacks) Set(idx, val int) {
	var ini int64 = 0
	(*st)[idx] = &ini
	atomic.StoreInt64((*st)[idx], int64(val))
}

func (st *Stacks) Push(idx, delta int) {
	atomic.AddInt64((*st)[idx], int64(delta))
}

func (st *Stacks) Pop(idx, delta int) {
	atomic.AddInt64((*st)[idx], -int64(delta))
}

func (st *Stacks) Next(idx int) bool {
	st.Pop(idx, 1)
	return *(*st)[idx] >= 0
}

func (st *Stacks) Has(idx int) bool {
	return (*st)[idx] != nil
}

func (pl *pipeline) Next(args Args) *pipeline {
	npl := &pipeline{}
	npl.parent = args
	npl.stacks = make(Stacks, len(pl.stacks))
	for i, sptr := range pl.stacks {
		if sptr != nil {
			stack := *sptr
			npl.stacks[i] = &stack
		}
	}
	return npl
}

func (args *argsStruct) pipeline(num int) *pipeline {
	if args.pl == nil {
		return newPipeline(num)
	}
	return args.pl
}
