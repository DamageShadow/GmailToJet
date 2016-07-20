package main

// import
import (
	"github.com/bobintornado/boltdb-boilerplate"
	"log"
	"fmt"
)

var buckets = []string{"ownerBucket", "sensors"}

func init() {
	err := boltdbboilerplate.InitBolt("./database.boltdb", buckets)
	if err != nil {
		log.Fatal("Can't init boltDB")
	}

}

func main() {
	// Put
	err := boltdbboilerplate.Put([]byte("ownerBucket"), []byte("ownerKey"), []byte("username"))
	if err != nil {
		log.Fatal("Can't init boltDB")
	}
	// Get owner
	value := boltdbboilerplate.Get([]byte("ownerBucket"), []byte("ownerKey"))

	fmt.Println(value)

	// Delete
	err = boltdbboilerplate.Delete([]byte("ownerBucket"), []byte("ownerKey"))

	// Insert two key/value
	err = boltdbboilerplate.Put([]byte("sensors"), []byte("key1"), []byte("value1"))
	err = boltdbboilerplate.Put([]byte("sensors"), []byte("key2"), []byte("value2"))

	// Get all keys
	keys := boltdbboilerplate.GetAllKeys([]byte("sensors"))
	// keys = [key1, key2]
	fmt.Println(keys)

	// Get all key/value pairs
	pairs := boltdbboilerplate.GetAllKeyValues([]byte("sensors"))
	// pairs = [{Key:key1, Value:value1}, {Key: key2, Value:value2}]

	fmt.Println(pairs)

	// Close
	boltdbboilerplate.Close()

}