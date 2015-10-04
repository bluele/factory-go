package factory

import (
	"reflect"
	"strconv"
	"sync/atomic"
)

var (
	TagName = "factory"
)

type Factory struct {
	model        interface{}
	numField     int
	rt           reflect.Type
	rv           *reflect.Value
	attrGens     []*attrGenerator
	nameIndexMap map[string]int // pair for attribute name and field index.
	isPtr        bool
}

type Args interface {
	Instance() interface{}
	pipeline(int) *Pipeline
}

type argsStruct struct {
	rv *reflect.Value
	pl *Pipeline
}

func (args *argsStruct) Instance() interface{} {
	return args.rv.Interface()
}

func (args *argsStruct) pipeline(num int) *Pipeline {
	if args.pl == nil {
		return NewPipeline(num)
	}
	return args.pl
}

type Stacks []*int64

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

type Pipeline struct {
	stacks Stacks
}

func NewPipeline(size int) *Pipeline {
	return &Pipeline{stacks: make(Stacks, size)}
}

func (pl *Pipeline) Copy() *Pipeline {
	npl := &Pipeline{}
	npl.stacks = make(Stacks, len(pl.stacks))
	for i, sptr := range pl.stacks {
		if sptr != nil {
			stack := *sptr
			npl.stacks[i] = &stack
		}
	}
	return npl
}

func NewFactory(model interface{}) *Factory {
	fa := &Factory{}
	fa.model = model
	fa.nameIndexMap = make(map[string]int)

	fa.init()
	return fa
}

type attrGenerator struct {
	genFunc func(Args) (interface{}, error)
	key     string
	value   interface{}
	isNil   bool
}

func (fa *Factory) init() {
	rt := reflect.TypeOf(fa.model)
	rv := reflect.ValueOf(fa.model)

	fa.isPtr = rt.Kind() == reflect.Ptr

	if fa.isPtr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	fa.numField = rv.NumField()

	for i := 0; i < fa.numField; i++ {
		tf := rt.Field(i)
		vf := rv.Field(i)
		ag := &attrGenerator{}

		if tf.Type.Kind() == reflect.Ptr && vf.IsNil() {
			ag.isNil = true
		} else {
			ag.value = vf.Interface()
		}

		attrName := getAttrName(tf, TagName)
		ag.key = attrName
		fa.nameIndexMap[attrName] = i
		fa.attrGens = append(fa.attrGens, ag)
	}

	fa.rt = rt
	fa.rv = &rv
}

func (fa *Factory) Attr(name string, gen func(Args) (interface{}, error)) *Factory {
	idx := fa.checkIdx(name)
	fa.attrGens[idx].genFunc = gen
	return fa
}

func (fa *Factory) SeqInt(name string, gen func(int) (interface{}, error)) *Factory {
	idx := fa.checkIdx(name)
	var seq int64 = 1
	fa.attrGens[idx].genFunc = func(args Args) (interface{}, error) {
		defer atomic.AddInt64(&seq, 1)
		return gen(int(seq))
	}
	return fa
}

func (fa *Factory) SeqString(name string, gen func(string) (interface{}, error)) *Factory {
	idx := fa.checkIdx(name)
	var seq int64 = 1
	fa.attrGens[idx].genFunc = func(args Args) (interface{}, error) {
		defer atomic.AddInt64(&seq, 1)
		return gen(strconv.FormatInt(seq, 10))
	}
	return fa
}

func (fa *Factory) SubFactory(name string, sub *Factory) *Factory {
	idx := fa.checkIdx(name)
	fa.attrGens[idx].genFunc = func(args Args) (interface{}, error) {
		return sub.Create()
	}
	return fa
}

func (fa *Factory) SubSliceFactory(name string, sub *Factory, getSize func() int) *Factory {
	idx := fa.checkIdx(name)
	tp := fa.rt.Field(idx).Type
	fa.attrGens[idx].genFunc = func(args Args) (interface{}, error) {
		size := getSize()
		sv := reflect.MakeSlice(tp, size, size)
		for i := 0; i < size; i++ {
			ret := sub.MustCreate()
			sv.Index(i).Set(reflect.ValueOf(ret))
		}
		return sv.Interface(), nil
	}
	return fa
}

func (fa *Factory) SubRecursiveFactory(name string, sub *Factory, getLimit func() int) *Factory {
	idx := fa.checkIdx(name)
	fa.attrGens[idx].genFunc = func(args Args) (interface{}, error) {
		pipeline := args.pipeline(fa.numField)
		if !pipeline.stacks.Has(idx) {
			pipeline.stacks.Set(idx, getLimit())
		}
		if pipeline.stacks.Next(idx) {
			ret, err := sub.create(nil, pipeline.Copy())
			if err != nil {
				return nil, err
			}
			return ret, nil
		}
		return nil, nil
	}
	return fa
}

func (fa *Factory) SubRecursiveSliceFactory(name string, sub *Factory, getSize, getLimit func() int) *Factory {
	idx := fa.checkIdx(name)
	tp := fa.rt.Field(idx).Type
	fa.attrGens[idx].genFunc = func(args Args) (interface{}, error) {
		pipeline := args.pipeline(fa.numField)
		if !pipeline.stacks.Has(idx) {
			pipeline.stacks.Set(idx, getLimit())
		}
		if pipeline.stacks.Next(idx) {
			size := getSize()
			sv := reflect.MakeSlice(tp, size, size)
			for i := 0; i < size; i++ {
				ret, err := sub.create(nil, pipeline.Copy())
				if err != nil {
					return nil, err
				}
				sv.Index(i).Set(reflect.ValueOf(ret))
			}
			return sv.Interface(), nil
		}
		return nil, nil
	}
	return fa
}

func (fa *Factory) checkIdx(name string) int {
	idx, ok := fa.nameIndexMap[name]
	if !ok {
		panic("No such atrribute name: " + name)
	}
	return idx
}

func (fa *Factory) Create() (interface{}, error) {
	return fa.CreateWithOption(nil)
}

func (fa *Factory) CreateWithOption(opt map[string]interface{}) (interface{}, error) {
	return fa.create(opt, nil)
}

func (fa *Factory) MustCreate() interface{} {
	return fa.MustCreateWithOption(nil)
}

func (fa *Factory) MustCreateWithOption(opt map[string]interface{}) interface{} {
	inst, err := fa.CreateWithOption(opt)
	if err != nil {
		panic(err)
	}
	return inst
}

func (fa *Factory) create(opt map[string]interface{}, pl *Pipeline) (interface{}, error) {
	inst := reflect.New(fa.rt).Elem()

	args := &argsStruct{}
	args.pl = pl
	if fa.isPtr {
		addr := inst.Addr()
		args.rv = &addr
	} else {
		args.rv = &inst
	}

	for i := 0; i < fa.numField; i++ {
		if v, ok := opt[fa.attrGens[i].key]; ok {
			inst.Field(i).Set(reflect.ValueOf(v))
		} else {
			ag := fa.attrGens[i]
			if ag.genFunc == nil {
				if !ag.isNil {
					inst.Field(i).Set(reflect.ValueOf(ag.value))
				}
			} else {
				v, err := ag.genFunc(args)
				if err != nil {
					return nil, err
				}
				if v != nil {
					inst.Field(i).Set(reflect.ValueOf(v))
				}
			}
		}
	}

	if fa.isPtr {
		inst = inst.Addr()
	}

	return inst.Interface(), nil
}
