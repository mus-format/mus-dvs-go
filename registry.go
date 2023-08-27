package dvs

import com "github.com/mus-format/common-go"

// TypeVersion defines a type version for Registry.
type TypeVersion any

// Creates a new Registry.
func NewRegistry(versions []TypeVersion) Registry {
	r := Registry{versions: make([]TypeVersion, len(versions))}
	copy(r.versions, versions)
	return r
}

// Registry contains all supported type versions, each of which should implement
// the MigrationVersion interface.
type Registry struct {
	versions []TypeVersion
}

// Get returns the type version by DTM.
//
// Returns ErrUnknownDTM if specified DTM is not in Registry.
func (r Registry) Get(dtm com.DTM) (tv TypeVersion, err error) {
	i := int(dtm)
	if i >= 0 && i < len(r.versions) {
		tv = r.versions[i]
	} else {
		err = ErrUnknownDTM
	}
	return
}
