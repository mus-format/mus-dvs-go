package dvs

import "errors"

// ErrUnknownDTM happens when DTM received from bs is not in Registry.
var ErrUnknownDTM = errors.New("unknown DTM")

// ErrWrongTypeVersion happens when type version from Registry is not
// MigrationVersion.
var ErrWrongTypeVersion = errors.New("TypeVersion is not MigrationVersion")
