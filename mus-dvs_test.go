package dvs

import (
	"bytes"
	"reflect"
	"testing"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

const (
	FooV1DTM com.DTM = iota
	FooV2DTM
	BarV1DTM
	BarV2DTM
)

// -----------------------------------------------------------------------------
type FooV1 struct {
	num int
}

func MarshalFooV1(foo FooV1, bs []byte) (n int) {
	return varint.MarshalInt(foo.num, bs)
}

func UnmarshalFooV1(bs []byte) (foo FooV1, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(bs)
	return
}

func SizeFooV1(foo FooV1) (size int) {
	return varint.SizeInt(foo.num)
}

var FooV1DTS = dts.New[FooV1](FooV1DTM,
	mus.MarshallerFn[FooV1](MarshalFooV1),
	mus.UnmarshallerFn[FooV1](UnmarshalFooV1),
	mus.SizerFn[FooV1](SizeFooV1))

type FooV2 struct {
	num int
	str string
}

func MarshalFooV2(foo FooV2, bs []byte) (n int) {
	n = varint.MarshalInt(foo.num, bs)
	n += ord.MarshalString(foo.str, nil, bs[n:])
	return
}

func UnmarshalFooV2(bs []byte) (foo FooV2, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(bs)
	if err != nil {
		return
	}
	var n1 int
	foo.str, n1, err = ord.UnmarshalString(nil, bs[n:])
	n += n1
	return
}

func SizeFooV2(foo FooV2) (size int) {
	size = varint.SizeInt(foo.num)
	return size + ord.SizeString(foo.str, nil)
}

var FooV2DTS = dts.New[FooV2](FooV1DTM,
	mus.MarshallerFn[FooV2](MarshalFooV2),
	mus.UnmarshallerFn[FooV2](UnmarshalFooV2),
	mus.SizerFn[FooV2](SizeFooV2))

// -----------------------------------------------------------------------------
type BarV1 struct {
	num int
}

func MarshalBarV1(bar BarV1, bs []byte) (n int) {
	return varint.MarshalInt(bar.num, bs)
}

func UnmarshalBarV1(bs []byte) (bar BarV1, n int, err error) {
	bar.num, n, err = varint.UnmarshalInt(bs)
	return
}

func SizeBarV1(bar BarV1) (size int) {
	return varint.SizeInt(bar.num)
}

var BarV1DTS = dts.New[BarV1](BarV1DTM,
	mus.MarshallerFn[BarV1](MarshalBarV1),
	mus.UnmarshallerFn[BarV1](UnmarshalBarV1),
	mus.SizerFn[BarV1](SizeBarV1))

type BarV2 struct {
	num int
	str string
}

func MarshalBarV2(bar BarV2, bs []byte) (n int) {
	n = varint.MarshalInt(bar.num, bs)
	n += ord.MarshalString(bar.str, nil, bs[n:])
	return
}

func UnmarshalBarV2(bs []byte) (bar BarV2, n int, err error) {
	bar.num, n, err = varint.UnmarshalInt(bs)
	if err != nil {
		return
	}
	var n1 int
	bar.str, n1, err = ord.UnmarshalString(nil, bs[n:])
	n += n1
	return
}

func SizeBarV2(bar BarV2) (size int) {
	size = varint.SizeInt(bar.num)
	return size + ord.SizeString(bar.str, nil)
}

var BarV2DTS = dts.New[BarV2](BarV1DTM,
	mus.MarshallerFn[BarV2](MarshalBarV2),
	mus.UnmarshallerFn[BarV2](UnmarshalBarV2),
	mus.SizerFn[BarV2](SizeBarV2))

// -----------------------------------------------------------------------------
type Foo FooV2
type Bar BarV2

func TestDVS(t *testing.T) {
	reg := com.NewRegistry([]com.TypeVersion{
		Version[FooV1, Foo]{
			DTS: FooV1DTS,
			MigrateOld: func(t FooV1) (v Foo, err error) {
				v.num = t.num
				v.str = "undefined"
				return
			},
			MigrateCurrent: func(v Foo) (t FooV1, err error) {
				t.num = v.num
				return
			},
		},
		Version[FooV2, Foo]{
			DTS: FooV2DTS,
			MigrateOld: func(t FooV2) (v Foo, err error) {
				return Foo(t), nil
			},
			MigrateCurrent: func(v Foo) (t FooV2, err error) {
				return FooV2(v), nil
			},
		},
		Version[BarV1, Bar]{
			DTS: BarV1DTS,
			MigrateOld: func(t BarV1) (v Bar, err error) {
				v.num = t.num
				v.str = "undefined"
				return
			},
			MigrateCurrent: func(v Bar) (t BarV1, err error) {
				t.num = v.num
				return
			},
		},
		Version[BarV2, Bar]{
			DTS: BarV2DTS,
			MigrateOld: func(t BarV2) (v Bar, err error) {
				return Bar(t), nil
			},
			MigrateCurrent: func(v Bar) (t BarV2, err error) {
				return BarV2(v), nil
			},
		},
	})

	fooDVS := New[Foo](reg)
	barDVS := New[Bar](reg)

	t.Run("MakeBSAndMarshal should marshal data, if it receives a correct registered data type",
		func(t *testing.T) {
			var (
				foo           = Foo{num: 5, str: "hello world"}
				wantBS        = []byte{0, 10}
				wantN         = 2
				wantErr error = nil
			)
			testMakeBSAndMarshal[Foo](fooDVS, FooV1DTM, foo, wantBS, wantN,
				wantErr,
				t)
		})

	t.Run("MakeBSAndMarshal should return ErrUnknownDTM, if it receives an unknown data type",
		func(t *testing.T) {
			var (
				wantBS  []byte = nil
				wantN          = 0
				wantErr error  = com.ErrUnknownDTM
			)
			testMakeBSAndMarshal[Foo](fooDVS, 5, Foo{}, wantBS, wantN, wantErr, t)
		})

	t.Run("MakeBSAndMarshal should return ErrWrongTypeVersion, if it receives an index of wrong type version",
		func(t *testing.T) {
			var (
				wantBS  []byte = nil
				wantN          = 0
				wantErr error  = com.ErrWrongTypeVersion
			)
			testMakeBSAndMarshal[Foo](fooDVS, BarV1DTM, Foo{}, wantBS, wantN,
				wantErr,
				t)
		})

	t.Run("ReliablyMarshal should marshal data, if it receives a correct registered data type and too big bs",
		func(t *testing.T) {
			var (
				foo           = Foo{num: 5, str: "hello world"}
				bs            = make([]byte, 1000)
				wantBS        = []byte{0, 10}
				wantN         = 2
				wantErr error = nil
			)
			testReliablyMarshal[Foo](fooDVS, 0, foo, bs, wantBS, wantN, wantErr, t)
		})

	t.Run("ReliablyMarshal should marshal data, if it receives a correct registered data type and too small bs",
		func(t *testing.T) {
			var (
				foo           = Foo{num: 5, str: "hello world"}
				wantBS        = []byte{0, 10}
				wantN         = 2
				wantErr error = nil
			)
			testReliablyMarshal[Foo](fooDVS, 0, foo, []byte{}, wantBS, wantN, wantErr,
				t)
		})

	t.Run("ReliablyMarshal should return ErrUnknownDTM, if it receives an unknown data type",
		func(t *testing.T) {
			var (
				wantBS  []byte = nil
				wantN          = 0
				wantErr error  = com.ErrUnknownDTM
			)
			testReliablyMarshal[Foo](fooDVS, 5, Foo{}, []byte{}, wantBS, wantN,
				wantErr,
				t)
		})

	t.Run("ReliablyMarshal should return ErrWrongTypeVersion, if it receives an index of wrong type version",
		func(t *testing.T) {
			var (
				wantBS  []byte = nil
				wantN          = 0
				wantErr error  = com.ErrWrongTypeVersion
			)
			testReliablyMarshal[Foo](fooDVS, 2, Foo{}, []byte{}, wantBS, wantN,
				wantErr,
				t)
		})

	t.Run("Unmarshal should unmarshal data, if there is a correct registered data type in the bs",
		func(t *testing.T) {
			var (
				FooV1 = FooV1{num: 5}
				bs    = func() []byte {
					bs := make([]byte, FooV1DTS.Size(FooV1))
					FooV1DTS.Marshal(FooV1, bs)
					return bs
				}()
				wantDT        = FooV1DTM
				wantFoo       = Foo{num: FooV1.num, str: "undefined"}
				wantN         = 2
				wantErr error = nil
			)
			testUnmarshal[Foo](fooDVS, bs, wantDT, wantFoo, wantN, wantErr, t)
		})

	t.Run("Unmarshal should return ErrUnknownDTM, if there is an unknown data type in the bs",
		func(t *testing.T) {
			var (
				bs              = []byte{10}
				wantDT  com.DTM = 5
				wantFoo         = Foo{}
				wantN           = 1
				wantErr error   = com.ErrUnknownDTM
			)
			testUnmarshal[Foo](fooDVS, bs, wantDT, wantFoo, wantN, wantErr, t)
		})

	t.Run("Unmarshal should return ErrWrongTypeVersion, if there is an wrong data type in the bs",
		func(t *testing.T) {
			var (
				bs              = []byte{4}
				wantDT  com.DTM = 2
				wantFoo         = Foo{}
				wantN           = 1
				wantErr error   = com.ErrWrongTypeVersion
			)
			testUnmarshal[Foo](fooDVS, bs, wantDT, wantFoo, wantN, wantErr, t)
		})

	t.Run("If UnmarshalDataType fails with an error, Unmarshal should return it",
		func(t *testing.T) {
			var (
				bs              = []byte{}
				wantDT  com.DTM = 0
				wantFoo         = Foo{}
				wantN           = 0
				wantErr error   = mus.ErrTooSmallByteSlice
			)
			testUnmarshal[Foo](fooDVS, bs, wantDT, wantFoo, wantN, wantErr, t)
		})

	t.Run("We should be able to use same registry for several DVS",
		func(t *testing.T) {
			var wantErr error = nil
			_, _, err := fooDVS.MakeBSAndMarshal(1, Foo{})
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
			_, _, err = barDVS.MakeBSAndMarshal(3, Bar{})
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
		})

}

func testUnmarshal[V any](dvs DVS[V], bs []byte, wantDT com.DTM,
	wantFoo Foo,
	wantN int,
	wantErr error,
	t *testing.T,
) {
	dtm, v, n, err := dvs.Unmarshal(bs)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if dtm != wantDT {
		t.Errorf("unexpected dtm, want '%v' actual '%v'", wantDT, dtm)
	}
	if !reflect.DeepEqual(v, wantFoo) {
		t.Errorf("unexpected v, want '%v' actual '%v'", wantFoo, v)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
}

func testReliablyMarshal[V any](dvs DVS[V], dtm com.DTM, v V, bs []byte,
	wantBS []byte,
	wantN int,
	wantErr error,
	t *testing.T,
) {
	bs, n, err := dvs.ReliablyMarshal(dtm, v, bs)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if !bytes.Equal(bs[:n], wantBS) {
		t.Errorf("unexpected bs, want '%v' actual '%v'", wantBS, bs)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
}

func testMakeBSAndMarshal[V any](dvs DVS[V], dtm com.DTM, v V, wantBS []byte,
	wantN int,
	wantErr error,
	t *testing.T,
) {
	bs, n, err := dvs.MakeBSAndMarshal(dtm, v)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if !bytes.Equal(bs, wantBS) {
		t.Errorf("unexpected bs, want '%v' actual '%v'", wantBS, bs)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
}
