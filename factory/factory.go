package factory

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"sync/atomic"
)

var (
	TagName    = "factory"
	emptyValue = reflect.Value{}
)

type (
	Args interface {
		Instance() any
		Parent() Args
		Context() context.Context
		pipeline(int) *pipeline
	}
	Formatter func(any) (any, error)
	Generator func(Args) (any, error)
	Factory   struct {
		numField         int
		curIdx           int
		isPtr            bool
		model            any
		rt               reflect.Type
		rv               *reflect.Value
		attrGens         []*attrGenerator
		orderingAttrGens []*attrGenerator
		nameIndexMap     map[string]int // pair for attribute name and field index.
		onCreate         func(Args) error
	}
)

// NewFactory returns a new factory for specified model class
// Each generator is applied in the order in which they are declared
func NewFactory(model any) *Factory {
	fa := &Factory{}
	fa.model = model
	fa.nameIndexMap = make(map[string]int)

	fa.init()
	return fa
}

func (fa *Factory) wrapWithFormatter(gen Generator, formatters ...Formatter) Generator {
	return func(a Args) (any, error) {
		ret, err := gen(a)
		if err != nil {
			return nil, err
		}
		for _, f := range formatters {
			ret, err = f(ret)
			if err != nil {
				return nil, err
			}
		}
		return ret, nil
	}
}

func (fa *Factory) Attr(name string, gen Generator, formatters ...Formatter) *Factory {
	return fa.fillAttrGen(nil, name, fa.wrapWithFormatter(gen, formatters...))
}

func (fa *Factory) SeqInt(name string, gen func(int) (any, error), formatters ...Formatter) *Factory {
	seq := int64(0)
	genFunc := func(krgs Args) (any, error) {
		new := atomic.AddInt64(&seq, 1)
		return gen(int(new))
	}
	return fa.fillAttrGen(nil, name, fa.wrapWithFormatter(genFunc, formatters...))
}

func (fa *Factory) SeqInt64(name string, gen func(int64) (any, error), formatters ...Formatter) *Factory {
	seq := int64(0)
	genFunc := func(args Args) (any, error) {
		new := atomic.AddInt64(&seq, 1)
		return gen(new)
	}
	return fa.fillAttrGen(nil, name, fa.wrapWithFormatter(genFunc, formatters...))
}

func (fa *Factory) SeqString(name string, gen func(string) (any, error), formatters ...Formatter) *Factory {
	seq := int64(0)
	genFunc := func(args Args) (any, error) {
		new := atomic.AddInt64(&seq, 1)
		return gen(strconv.FormatInt(new, 10))
	}
	return fa.fillAttrGen(nil, name, fa.wrapWithFormatter(genFunc, formatters...))
}

func (fa *Factory) SubFactory(name string, sub *Factory, formatters ...Formatter) *Factory {
	genFunc := func(args Args) (any, error) {
		pipeline := args.pipeline(fa.numField)
		ret, err := sub.create(args.Context(), nil, pipeline.Next(args))
		if err != nil {
			return nil, err
		}
		return ret, nil
	}
	return fa.fillAttrGen(nil, name, fa.wrapWithFormatter(genFunc, formatters...))
}

func (fa *Factory) SubSliceFactory(name string, sub *Factory, getSize func() int, formatters ...Formatter) *Factory {
	idx := fa.checkIdx(name)
	tp := fa.rt.Field(idx).Type
	genFunc := func(args Args) (any, error) {
		size := getSize()
		pipeline := args.pipeline(fa.numField)
		sv := reflect.MakeSlice(tp, size, size)
		for i := 0; i < size; i++ {
			ret, err := sub.create(args.Context(), nil, pipeline.Next(args))
			if err != nil {
				return nil, err
			}
			sv.Index(i).Set(reflect.ValueOf(ret))
		}
		return sv.Interface(), nil
	}
	return fa.fillAttrGen(&idx, name, fa.wrapWithFormatter(genFunc, formatters...))
}

func (fa *Factory) SubRecursiveFactory(name string, sub *Factory, getLimit func() int, formatters ...Formatter) *Factory {
	idx := fa.checkIdx(name)
	genFunc := func(args Args) (any, error) {
		pl := args.pipeline(fa.numField)
		if !pl.stacks.Has(idx) {
			pl.stacks.Set(idx, getLimit())
		}
		if pl.stacks.Next(idx) {
			ret, err := sub.create(args.Context(), nil, pl.Next(args))
			if err != nil {
				return nil, err
			}
			return ret, nil
		}
		return nil, nil
	}
	return fa.fillAttrGen(&idx, name, fa.wrapWithFormatter(genFunc, formatters...))
}

func (fa *Factory) SubRecursiveSliceFactory(name string, sub *Factory, getSize, getLimit func() int, formatters ...Formatter) *Factory {
	idx := fa.checkIdx(name)
	tp := fa.rt.Field(idx).Type
	genFunc := func(args Args) (any, error) {
		pl := args.pipeline(fa.numField)
		if !pl.stacks.Has(idx) {
			pl.stacks.Set(idx, getLimit())
		}
		if pl.stacks.Next(idx) {
			size := getSize()
			sv := reflect.MakeSlice(tp, size, size)
			for i := 0; i < size; i++ {
				ret, err := sub.create(args.Context(), nil, pl.Next(args))
				if err != nil {
					return nil, err
				}
				sv.Index(i).Set(reflect.ValueOf(ret))
			}
			return sv.Interface(), nil
		}
		return nil, nil
	}
	return fa.fillAttrGen(&idx, name, fa.wrapWithFormatter(genFunc, formatters...))
}

// OnCreate registers a callback on object creation.
// If callback function returns error, object creation is failed.
func (fa *Factory) OnCreate(cb func(Args) error) *Factory {
	fa.onCreate = cb
	return fa
}

func (fa *Factory) Create() (any, error) {
	return fa.CreateWithOption(nil)
}

func (fa *Factory) CreateWithOption(opt map[string]any) (any, error) {
	return fa.create(context.Background(), opt, nil)
}

func (fa *Factory) CreateWithContext(ctx context.Context) (any, error) {
	return fa.create(ctx, nil, nil)
}

func (fa *Factory) CreateWithContextAndOption(ctx context.Context, opt map[string]any) (any, error) {
	return fa.create(ctx, opt, nil)
}

func (fa *Factory) MustCreate() any {
	return fa.MustCreateWithOption(nil)
}

func (fa *Factory) MustCreateWithOption(opt map[string]any) any {
	return fa.MustCreateWithContextAndOption(context.Background(), opt)
}

func (fa *Factory) MustCreateWithContextAndOption(ctx context.Context, opt map[string]any) any {
	inst, err := fa.CreateWithContextAndOption(ctx, opt)
	if err != nil {
		panic(err)
	}
	return inst
}

/*
Bind values of a new objects to a pointer to struct.

ptr: a pointer to struct
*/
func (fa *Factory) Construct(ptr any) error {
	return fa.ConstructWithOption(ptr, nil)
}

/*
Bind values of a new objects to a pointer to struct with option.

ptr: a pointer to struct
opt: attibute values
*/
func (fa *Factory) ConstructWithOption(ptr any, opt map[string]any) error {
	return fa.ConstructWithContextAndOption(context.Background(), ptr, opt)
}

/*
Bind values of a new objects to a pointer to struct with context and option.

ctx: context object
ptr: a pointer to struct
opt: attibute values
*/
func (fa *Factory) ConstructWithContextAndOption(ctx context.Context, ptr any, opt map[string]any) error {
	pt := reflect.TypeOf(ptr)
	if pt.Kind() != reflect.Ptr {
		return errors.New("ptr should be pointer type.")
	}
	pt = pt.Elem()
	if pt.Name() != fa.modelName() {
		return errors.New("ptr type should be " + fa.modelName())
	}

	inst := reflect.ValueOf(ptr).Elem()
	_, err := fa.build(ctx, &inst, pt, opt, nil)
	return err
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
	fa.orderingAttrGens = make([]*attrGenerator, fa.numField)

	for i := 0; i < fa.numField; i++ {
		tf := rt.Field(i)
		vf := rv.Field(i)
		ag := &attrGenerator{}

		if !vf.CanSet() || (tf.Type.Kind() == reflect.Ptr && vf.IsNil()) {
			ag.isNil = true
		} else {
			ag.value = vf.Interface()
		}

		attrName := getAttrName(tf, TagName)
		ag.key = attrName
		ag.isFilled = false
		fa.nameIndexMap[attrName] = i
		fa.attrGens = append(fa.attrGens, ag)
	}

	fa.rt = rt
	fa.rv = &rv
}

func (fa *Factory) modelName() string {
	return fa.rt.Name()
}

func (fa *Factory) fillAttrGen(idx *int, name string, gen func(Args) (any, error)) *Factory {
	if idx == nil {
		i := fa.checkIdx(name)
		idx = &i
	}
	fa.attrGens[*idx].genFunc = gen
	fa.attrGens[*idx].isFilled = true
	orderingIdx := fa.getOrderingIdx()
	fa.orderingAttrGens[orderingIdx] = fa.attrGens[*idx]
	return fa
}

func (fa *Factory) checkIdx(name string) int {
	idx, ok := fa.nameIndexMap[name]
	if !ok {
		panic("No such attribute name: " + name)
	}
	return idx
}
func (fa *Factory) getOrderingIdx() int {
	idx := fa.curIdx
	if fa.curIdx < fa.numField-1 {
		fa.curIdx += 1
	}
	return idx
}

func (fa *Factory) fillMissingAttr(ctx context.Context) {
	for _, attr := range fa.attrGens {
		if !attr.isFilled {
			fa.fillAttrGen(nil, attr.key, attr.genFunc)
		}
	}
}

func (fa *Factory) build(ctx context.Context, inst *reflect.Value, tp reflect.Type, opt map[string]any, pl *pipeline) (any, error) {
	args := &argsStruct{}
	args.pl = pl
	args.ctx = ctx
	if fa.isPtr {
		addr := (*inst).Addr()
		args.rv = &addr
	} else {
		args.rv = inst
	}

	fa.fillMissingAttr(ctx)

	for _, attr := range fa.orderingAttrGens {
		if v, ok := opt[attr.key]; ok {
			inst.FieldByName(attr.key).Set(reflect.ValueOf(v))
		} else {
			if attr.genFunc == nil {
				if !attr.isNil {
					inst.FieldByName(attr.key).Set(reflect.ValueOf(attr.value))
				}
			} else {
				v, err := attr.genFunc(args)
				if err != nil {
					return nil, err
				}
				if v != nil {
					inst.FieldByName(attr.key).Set(reflect.ValueOf(v))
				}
			}
		}
	}

	for k, v := range opt {
		setValueWithAttrPath(inst, tp, k, v)
	}

	if fa.onCreate != nil {
		if err := fa.onCreate(args); err != nil {
			return nil, err
		}
	}

	if fa.isPtr {
		return (*inst).Addr().Interface(), nil
	}
	return inst.Interface(), nil
}

func (fa *Factory) create(ctx context.Context, opt map[string]any, pl *pipeline) (any, error) {
	inst := reflect.New(fa.rt).Elem()
	return fa.build(ctx, &inst, fa.rt, opt, pl)
}
