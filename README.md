mgo.ejson
=========

A simple package to Marshal/Unmarshal MongoDB's EJson in Go.

# ejson
--
    import "github.com/OneOfOne/mgo.ejson"


# Warning

This is alpha quality at best, and while it works for most ejson types, it is still *very* expermental.

## Usage

#### func  ConvertEjson

```go
func ConvertEjson(j []byte) (bson.M, error)
```
ConvertEjson takes in ejson and converts it to a clean **bson.M** object that can be
used with bson.Marshal.

#### func  Unmarshal

```go
func Unmarshal(j []byte, v interface{}) error
```
Unmarshal takes ejson and a pointer to a struct or a map, converts up the ejson
to a clean **bson.M** then sets **v**.

***warning*** the current implementation is rather slow.

How it works: ejson -> json.Unmarshal -> cleanup -> bson.Marshal -> bson.Unmarshal

