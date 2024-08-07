package dvs

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/mus-format/mus-go"
)

func TestVersion(t *testing.T) {
	var (
		ver = Version[FooV1, Foo]{
			DTS: FooV1DTS,
			MigrateOld: func(t FooV1) (v Foo, err error) {
				return Foo{num: t.num, str: "undefined"}, nil
			},
			MigrateCurrent: func(v Foo) (t FooV1, err error) {
				return FooV1{num: v.num}, nil
			},
		}
	)

	t.Run("MigrateCurrentAndReliablyMarshal should marshal data if it receives too big bs",
		func(t *testing.T) {
			var (
				foo           = Foo{num: 11}
				bs            = make([]byte, 1000)
				wantBS        = []byte{0, 22}
				wantN         = 2
				wantErr error = nil
			)
			testMigrateCurrentAndReliablyMarshal(ver, foo, bs, wantBS, wantN, wantErr,
				t)
		})

	t.Run("MigrateCurrentAndReliablyMarshal should marshal data if it receives too small bs",
		func(t *testing.T) {
			var (
				foo           = Foo{num: 11}
				bs            = []byte{}
				wantBS        = []byte{0, 22}
				wantN         = 2
				wantErr error = nil
			)
			testMigrateCurrentAndReliablyMarshal(ver, foo, bs, wantBS, wantN, wantErr,
				t)
		})

	t.Run("If MigrateCurrent fails with an error, MigrateCurrentAndReliablyMarshal should return it",
		func(t *testing.T) {
			var (
				wantBS  []byte = nil
				wantN          = 0
				wantErr        = errors.New("MigrateCurrent error")
				ver            = Version[FooV1, Foo]{
					MigrateCurrent: func(v Foo) (t FooV1, err error) {
						err = wantErr
						return
					},
				}
			)
			testMigrateCurrentAndReliablyMarshal(ver, Foo{}, []byte{},
				wantBS,
				wantN,
				wantErr,
				t)
		})

	t.Run("MigrateCurrentAndMakeBSAndMarshal should marshal data",
		func(t *testing.T) {
			var (
				foo           = Foo{num: 22}
				wantBS        = []byte{0, 44}
				wantN         = 2
				wantErr error = nil
			)
			testMigrateCurrentAndMakeBSAndMarshal(ver, foo, wantBS, wantN, wantErr, t)
		})

	t.Run("If MigrateCurrent fails with an error, MigrateCurrentAndMakeBSAndMarshal should return it",
		func(t *testing.T) {
			var (
				wantBS  []byte = nil
				wantN          = 0
				wantErr        = errors.New("MigrateCurrent error")
				ver            = Version[FooV1, Foo]{
					MigrateCurrent: func(v Foo) (t FooV1, err error) {
						err = wantErr
						return
					},
				}
			)
			testMigrateCurrentAndMakeBSAndMarshal(ver, Foo{}, wantBS, wantN,
				wantErr,
				t)
		})

	t.Run("UnmarshalAndMigrateOld should unmarshal data", func(t *testing.T) {
		var (
			bs            = []byte{22}
			wantFoo       = Foo{num: 11, str: "undefined"}
			wantN         = 1
			wantErr error = nil
		)
		testUnmarshalAndMigrateOld[FooV1, Foo](ver, bs, wantFoo, wantErr, wantN, t)
	})

	t.Run("If UnmarshalData fails with an error, UnmarshalAndMigrateOld should return it",
		func(t *testing.T) {
			var (
				bs            = []byte{}
				wantFoo       = Foo{}
				wantN         = 0
				wantErr error = mus.ErrTooSmallByteSlice
			)
			testUnmarshalAndMigrateOld[FooV1, Foo](ver, bs, wantFoo, wantErr, wantN,
				t)
		})

}

func testMigrateCurrentAndReliablyMarshal[T, V any](ver Version[T, V], v V,
	bs []byte,
	wantBS []byte,
	wantN int,
	wantErr error,
	t *testing.T,
) {
	bs, n, err := ver.MigrateCurrentAndReliablyMarshal(v, bs)
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

func testUnmarshalAndMigrateOld[T, V any](ver Version[T, V], bs []byte,
	wantV V,
	wantErr error,
	wantN int,
	t *testing.T,
) {
	v, n, err := ver.UnmarshalAndMigrateOld(bs)
	if err != wantErr {
		t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
	}
	if !reflect.DeepEqual(v, wantV) {
		t.Errorf("unexpected v, want '%v' actual '%v'", wantV, v)
	}
	if n != wantN {
		t.Errorf("unexpected n, want '%v' actual '%v'", wantN, n)
	}
}

func testMigrateCurrentAndMakeBSAndMarshal[T, V any](ver Version[T, V], v V,
	wantBS []byte,
	wantN int,
	wantErr error,
	t *testing.T,
) {
	bs, n, err := ver.MigrateCurrentAndMakeBSAndMarshal(v)
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
