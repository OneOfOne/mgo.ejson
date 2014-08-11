package ejson

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Normalize takes in an ejson map and normalizes it to regular bson.
func Normalize(m map[string]interface{}) error {
	for key := range m {
		if sm, ok := m[key].(map[string]interface{}); ok {
			if d, ok := sm["$date"].(float64); ok {
				tm := int64(d * 0.001) // convert the milliseconds to seconds
				t := time.Unix(tm, 0)
				m[key] = &t
			} else if id, ok := sm["$oid"].(string); ok {
				var oid bson.ObjectId
				if bson.IsObjectIdHex(id) {
					oid = bson.ObjectIdHex(id)
				} else {
					oid = bson.ObjectId(id)
				}
				m[key] = &oid
			} else if b, ok := sm["$undefined"].(bool); ok {
				if b {
					m[key] = nil
				}
			} else if ref, ok := sm["$ref"].(string); ok {
				if id, ok := sm["$id"].(string); ok {
					var oid bson.ObjectId
					if bson.IsObjectIdHex(id) {
						oid = bson.ObjectIdHex(id)
					} else {
						oid = bson.ObjectId(id)
					}
					dbref := mgo.DBRef{
						Collection: ref,
						Id:         oid,
					}
					m[key] = &dbref
				} else {
					return fmt.Errorf(`ejson: "%s" expected a {"$ref": "string", "$id": "hex"} got "%+v"`, key, sm)
				}
			} else if bin, ok := sm["$binary"].(string); ok {
				b64, err := base64.StdEncoding.DecodeString(bin)
				if err != nil {
					return fmt.Errorf(`ejson: "%s" expected valid a base64 string in $binary, got: %v`, key, err)
				}
				if t, ok := sm["$type"].(string); ok {
					ti, err := strconv.ParseUint(t, 16, 8)
					if err != nil {
						return fmt.Errorf(`ejson: "%s" expected a valid hex byte in $type, got: %v`, key, err)
					}
					bin := bson.Binary{
						Kind: byte(ti),
						Data: b64,
					}
					m[key] = &bin
				}
			} else {
				if err := Normalize(sm); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//Unmarshal takes ejson and a pointer to a struct or a map,
// converts up the ejson to a clean **bson.M** then sets **v**.
// ***warning*** the current implementation is rather slow.
//
//How it works: ejson -> json.Unmarshal -> cleanup -> bson.Marshal -> bson.Unmarshal
func Unmarshal(j []byte, v interface{}) error {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		return err
	}
	if err := Normalize(m); err != nil {
		return err
	}
	js, err := bson.Marshal(m)
	if err != nil {
		return err
	}
	return bson.Unmarshal(js, v)
}
