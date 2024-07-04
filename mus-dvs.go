package dvs

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
)

// New creates a new DVS.
func New[V any](reg com.Registry) DVS[V] {
	return DVS[V]{reg}
}

// DVS provides versioning support for the mus-go serializer.
type DVS[V any] struct {
	reg com.Registry
}

// MakeBSAndMarshal makes bs, migrates v to the version specified by dtm, and
// then marshals dtm + resulting v version.
//
// In addition to the created byte slice and the number of used bytes, it can
// also return ErrUnknownDTM or ErrWrongTypeVersion.
func (dvs DVS[V]) MakeBSAndMarshal(dtm com.DTM, v V) (bs []byte, n int,
	err error) {
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	return mver.MigrateCurrentAndMakeBSAndMarshal(v)
}

// ReliablyMarshal migrates v to the version specified by dtm, and then reliably
// (if bs is too small creates a new one) marshals dtm + resulting v version.
//
// In addition to the received or created byte slice and the number of
// used bytes, it can also return ErrUnknownDTM or ErrWrongTypeVersion.
func (dvs DVS[V]) ReliablyMarshal(dtm com.DTM, v V, bs []byte) (abs []byte,
	n int, err error) {
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	return mver.MigrateCurrentAndReliablyMarshal(v, bs)
}

// Unmarshal unmarshals dtm + data, and then migrates data to the version
// specified by dtm.
//
// In addition to dtm and migrated dataand the number of used bytes, it can
// also return ErrUnknownDTM or ErrWrongTypeVersion.
func (dvs DVS[V]) Unmarshal(bs []byte) (dtm com.DTM, v V, n int, err error) {
	dtm, n, err = dts.UnmarshalDTM(bs)
	if err != nil {
		return
	}
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	var n1 int
	v, n1, err = mver.UnmarshalAndMigrateOld(bs[n:])
	n += n1
	return
}

func (dvs DVS[V]) getMV(dtm com.DTM) (mver MigrationVersion[V], err error) {
	tver, err := dvs.reg.Get(dtm)
	if err != nil {
		return
	}
	mver, ok := tver.(MigrationVersion[V])
	if !ok {
		err = com.ErrWrongTypeVersion
		return
	}
	return
}
