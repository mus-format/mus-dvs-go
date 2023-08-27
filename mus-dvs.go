package dvs

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
)

// New creates a new DVS.
func New[V any](reg Registry) DVS[V] {
	return DVS[V]{reg}
}

// DVS provides versioning support for the mus-go serializer.
type DVS[V any] struct {
	reg Registry
}

// MakeBSAndMarshalMUS makes bs, migrates v to the version specified by dtm,
// marshals dtm and the resulting v version to the MUS format.
//
// In addition to the created byte slice, returns the number of used bytes and
// one of the ErrUnknownDTM or ErrWrongTypeVersion errors.
func (dvs DVS[V]) MakeBSAndMarshalMUS(dtm com.DTM, v V) (bs []byte, n int,
	err error) {
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	return mver.MigrateCurrentAndMakeBSAndMarshalMUS(v)
}

// ReliablyMarshalMUS migrates v to the version specified by dtm, reliably (if
// bs is too small creates a new one) marshals dtm and the resulting v version
// to the MUS format.
//
// In addition to the received or created byte slice, returns the number of
// used bytes and one of the ErrUnknownDTM or ErrWrongTypeVersion errors.
func (dvs DVS[V]) ReliablyMarshalMUS(dtm com.DTM, v V, bs []byte) (abs []byte,
	n int, err error) {
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	return mver.MigrateCurrentAndReliablyMarshalMUS(v, bs)
}

// UnmarshalMUS unmarshals dtm and data from the MUS format, migrates data to
// the version specified by dtm.
//
// In addition to dtm and migrated data, returns the number of used bytes and
// one of the ErrUnknownDTM or ErrWrongTypeVersion errors.
func (dvs DVS[V]) UnmarshalMUS(bs []byte) (dtm com.DTM, v V, n int, err error) {
	dtm, n, err = dts.UnmarshalDTMUS(bs)
	if err != nil {
		return
	}
	mver, err := dvs.getMV(dtm)
	if err != nil {
		return
	}
	var n1 int
	v, n1, err = mver.UnmarshalAndMigrateOldMUS(bs[n:])
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
		err = ErrWrongTypeVersion
		return
	}
	return
}
