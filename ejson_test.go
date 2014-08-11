package ejson

import (
	"testing"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const j = `{"_id":{"$oid":"53c2ab5e4291b17b666d742a"},"last_seen_at":{"$date":1405266782008},"display_name":{"$undefined":true},
"ref":{"$ref":"col2", "$id":"53c2ab5e4291b17b666d742b"},
"d":1405266782008}`

type TestS struct {
	Id          bson.ObjectId `bson:"_id"`
	LastSeenAt  *time.Time    `bson:"last_seen_at"`
	DisplayName *string       `bson:"display_name,omitempty"`
	Ref         mgo.DBRef     `bson:"ref"`
	D           uint64        `bson:"d"`
}

func TestUnmarshal(t *testing.T) {
	var ts TestS
	if err := Unmarshal([]byte(j), &ts); err != nil {
		t.Fatal(err)
	}
	if ts.Id != bson.ObjectIdHex("53c2ab5e4291b17b666d742a") {
		t.Fatal("Unexpected ts.Id")
	}
	if ts.Ref.Id != bson.ObjectIdHex("53c2ab5e4291b17b666d742b") {
		t.Fatal("Unexpected ts.Ref.Id")
	}
}
