# mus-dvs-go
mus-dvs-go provides data versioning support for the [mus-go](https://github.com/mus-format/mus-go) 
serializer. With mus-dvs-go we can do 2 things:
1. Marshal the current data version as if it was an old version.
2. Unmarshal the old data version as if it was the current version.

# What is this for?
In the client-server architecture, it would be great to use on the server side 
only the current version of data, despite the fact that the server have to 
support outdated clients - accept/send old versions of data from/to them.

A similar situation arises when working with a storage that stores both 
current and old versions of data. It would be great to always get only the 
current version from it.

mus-dvs-go could help us with this. Let's see how. But first, to learn why 
everything is organized this way, read the [Data Versioning section](https://github.com/mus-format/specification#data-versioning) from the MUS format specification. Also read the 
mus-dts-go module [documentation](https://github.com/mus-format/mus-dts-go).

# Tests
Test coverage is 100%.

# How To Use
```go
package main

// Suppose we have the following types (each of which indicates the current 
// version):
type Foo FooV2
type Bar BarV1

// , and following versions:
type FooV1 struct {...}
type FooV2 struct {...}
type BarV1 struct {...}

// , and DTS definitions for them.
var FooV1DTS = ...
var FooV2DTS = ...
var BarV1DTS = ...

// To create a versioning support for our types, we first need to create a 
// registry of all the versions we support. To do this, we use the dvs.Version 
// type. This type is actually quite simple, it contains DTS and migrate 
// functions, one of which migrates the old version to the current one, and the 
// other do the opposite - migrates current to old.
//
// PLEASE NOTE that the index of each version in the registry must be equal
// to its DTM, i.e:
//
//	registry[FooV1DTM] == dvs.Version[FooV1, Foo]
//	registry[FooV2DTM] == dvs.Version[FooV2, Foo]
//	registry[BarV1DTM] == dvs.Version[BarV1, Foo]
//
// , thanks to this we can get the necessary version from the registry very 
// quickly.
var registry = dvs.NewRegistry(
  []dvs.TypeVersion{
    // FooV1 version.
    dvs.Version[FooV1, Foo]{
      DTS: FooV1DTS,
      MigrateOld: func(t FooV1) (v Foo, err error) {...},
      MigrateCurrent: func(v Foo) (t FooV1, err error) {...},
    },
    // FooV2 version.
    dvs.Version[FooV2, Foo]{
      DTS: FooV2DTS,
      MigrateOld: func(t FooV2) (v Foo, err error) {...},
      MigrateCurrent: func(v Foo) (t FooV2, err error) {...},
    },
    // BarV1 version.
    dvs.Version[BarV1, Bar]{
      DTS: BarV1DTS,
      MigrateOld: func(t BarV1) (v Bar, err error) {...},
      MigrateCurrent: func(v Bar) (t BarV1, err error) {...},
    },
  },
)

// And finally we can create versioning support for our types. Please note that 
// we use a single registry for all the DVS types.
var FooDVS = dvs.New[Foo](registry)
var BarDVS = dvs.New[Bar](registry)

func main() {
  var (
    foo Foo
    bs  []byte
  )
  // 1. Marshal the current version as if it were an old version.
  foo = Foo{...}
  // Migrates the current version to FooV1, and then marshals it to bs.
  bs, _, _ = FooDVS.MakeBSAndMarshalMUS(FooV1DTM, foo)

  // We should find the FooV1 version in bs.
  fooV1, _, _ := FooV1DTS.UnmarshalMUS(bs)
  assert.EqualDeep(fooV1, FooV1{...})

  // 2. Unmarshal the old version as if it were the current version.
  // Fills bs with the FooV1 version.
  fooV1 = FooV1{...}
  bs = make([]byte, FooV1DTS.SizeMUS(fooV1))
  FooV1DTS.MarshalMUS(fooV1, bs)

  // Unmarshals the FooV1 version from bs, and then migrates it to the current 
  // version.
  dt, foo, _, _ := FooDVS.Unmarshal(bs)

  // We have to get the correct Foo.
  assert.Equal[com.DTM](dt, FooV1DTM)
  assert.EqualDeep(foo, Foo{...})
}
```
You can find the full code at [mus-examples-go](https://github.com/mus-format/mus-examples-go/tree/main/dvs).

In summary, to communicate with the server the client should:
```go
...
// Marshal data to the bs.
bs := make([]byte, FooV1DTS.SizeMUS(fooV1))
FooV1DTS.MarshalMUS(fooV1, bs)
// Send bs to the server.
...
```
, and the server should:
```go
// Receive bs from the client.
...
bs := ...
// Unmarshal the current version from bs and use it.
_, foo, _, err := FooDVS.UnmarshalMUS(bs[n:])
...
```
Also, if the client expects to receive data, it must transmit a DTM to the 
server. This DTM will be used on the server to send the version of the data 
required by the client:
```go
// Receive DTM from the client.
...
dmt := ...
// Marshal appropriate version.
bs, _ := FooDVS.MakeBSAndMarshal(dtm, foo)
// Send bs to the client.
... 
```
You can find the full code at [mus-examples-go](https://github.com/mus-format/mus-examples-go/tree/main/rest).