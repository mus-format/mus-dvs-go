package dvs

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
)

// MigrationVersion represents a generic type version for Registry that can
// be migrated.
//
// It contains methods to support all mus-dvs-go functionality.
type MigrationVersion[V any] interface {
	MigrateCurrentAndReliablyMarshal(v V, bs []byte) (abs []byte, n int,
		err error)
	MigrateCurrentAndMakeBSAndMarshal(v V) (bs []byte, n int, err error)
	UnmarshalAndMigrateOld(bs []byte) (v V, n int, err error)
}

// Version is an implementation of the MigrationVersion interface.
type Version[T any, V any] struct {
	DTS            dts.DTS[T]
	MigrateOld     com.MigrateOld[T, V]
	MigrateCurrent com.MigrateCurrent[V, T]
}

func (ver Version[T, V]) MigrateCurrentAndReliablyMarshal(v V,
	bs []byte) (
	abs []byte, n int, err error) {
	t, err := ver.MigrateCurrent(v)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			abs, n = ver.makeBSAndMarshal(t)
		}
	}()
	n = ver.marshal(t, bs)
	abs = bs
	return
}

func (ver Version[T, V]) MigrateCurrentAndMakeBSAndMarshal(v V) (
	bs []byte, n int, err error) {
	t, err := ver.MigrateCurrent(v)
	if err != nil {
		return
	}
	bs, n = ver.makeBSAndMarshal(t)
	return
}

func (ver Version[T, V]) UnmarshalAndMigrateOld(bs []byte) (v V, n int,
	err error) {
	t, n, err := ver.DTS.UnmarshalData(bs)
	if err != nil {
		return
	}
	v, err = ver.MigrateOld(t)
	return
}

func (ver Version[T, V]) makeBSAndMarshal(t T) (bs []byte, n int) {
	bs = make([]byte, ver.DTS.Size(t))
	n = ver.marshal(t, bs)
	return
}

func (ver Version[T, V]) marshal(t T, bs []byte) (n int) {
	return ver.DTS.Marshal(t, bs)
}
