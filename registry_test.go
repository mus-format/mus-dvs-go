package dvs

import (
	"reflect"
	"testing"
)

func TestRegistry(t *testing.T) {
	var (
		ver1     = Version[FooV1, Foo]{}
		ver2     = Version[FooV2, Foo]{}
		versions = []TypeVersion{ver1, ver2}
		reg      = NewRegistry(versions)
	)

	t.Run("Get should return type version", func(t *testing.T) {
		ver, err := reg.Get(0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ver, ver1) {
			t.Errorf("unexpected ver, want '%v' actual '%v'", ver1, ver)
		}
		ver, err = reg.Get(1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ver, ver2) {
			t.Errorf("unexpected ver, want '%v' actual '%v'", ver1, ver)
		}
	})

	t.Run("Get should return ErrUnknownDTM if receives too big index",
		func(t *testing.T) {
			ver, err := reg.Get(5)
			if err != ErrUnknownDTM {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrUnknownDTM,
					err)
			}
			if ver != nil {
				t.Errorf("unexpected ver, want '%v' actual '%v'", nil, ver)
			}
		})

	t.Run("Get should return ErrUnknownDTM if receives negative index",
		func(t *testing.T) {
			ver, err := reg.Get(-1)
			if err != ErrUnknownDTM {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrUnknownDTM,
					err)
			}
			if ver != nil {
				t.Errorf("unexpected ver, want '%v' actual '%v'", nil, ver)
			}
		})

}
