package dvs

// MigrateOld migrates the old data version to the current one.
type MigrateOld[T, V any] func(t T) (v V, err error)

// MigrateCurrent migrates the current data version to the old one.
type MigrateCurrent[V, T any] func(v V) (t T, err error)
