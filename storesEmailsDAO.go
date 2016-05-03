package main

import "github.com/bobintornado/boltdb-boilerplate"

import (
	"github.com/goinggo/tracelog"
	"fmt"
)

var (
	storeBucket string = "storesEmails"
	defaultEmail string = "@"
)

func init() {
	// Init buckets and create file
	buckets := []string{storeBucket}

	err := boltdbboilerplate.InitBolt("./storesDatabase.boltdb", buckets)
	if err != nil {
		tracelog.Errorf(fmt.Errorf("Exception At..."), "storesEmailDAO", "init", "Can't init storesDatabase boltDB")
	}
}

func insertStore(store string) {
	// Put
	err := boltdbboilerplate.Put([]byte(storeBucket), []byte(store), []byte(defaultEmail))
	if err != nil {
		tracelog.Errorf(err, "storesEmailDAO", "insertStore", "Can't insert store " +
		store + " to boltDB")
	}
}

func deleteStore(store string) {
	// Delete
	err := boltdbboilerplate.Delete([]byte(storeBucket), []byte(store))
	if err != nil {
		tracelog.Errorf(err, "storesEmailDAO", "deleteStore", "Can't delete store " +
		store + " from boltDB")
	}
}

func insertAll(newStores []string) {
	//Get the list of stores, insert new, delete obsolete with empty contact email
	existingStores := getAllStores()

	newcomers, obsolete := difference(newStores, existingStores)

	if (newcomers != nil) {
		for i := range newcomers {
			insertStore(i)
		}
	}

	if (obsolete != nil) {
		for i := range obsolete {
			insertStore(i)
		}
	}
}

func getStore(store string) (string) {
	value := boltdbboilerplate.Get([]byte(storeBucket), []byte(store))

	return value
}

func insertStoreWithEmail(store string, email string) {
	err := boltdbboilerplate.Put([]byte(storeBucket), []byte(store), []byte(email))
	if err != nil {
		tracelog.Errorf(err, "storesEmailDAO", "insertStoreWithEmail", "Can't insert store " +
		store + " with email " + email + " to boltDB")
	} else {
		tracelog.Info("storesEmailDAO", "insertStoreWithEmail", "Update Store's  " + store +
		" email to " + email)
	}
}

func updateStoreWithEmail(store string, email string) {
	//Delete default pair-value for store if exists and insert store with email
	if (getStore(store) != nil) {
		deleteStore(store)
	}

	insertStoreWithEmail(store, email)
}

func getAllStores() []string {
	// Get all keys (stores)
	stores := boltdbboilerplate.GetAllKeys([]byte(storeBucket))
	tracelog.Info("storesEmailDAO", "getAllStores", "Extracted Stores : " + stores)
	return stores
}

func getAllStoresEmails() []string {
	// Get all values (emails)
	storesMap := boltdbboilerplate.GetAllKeyValues([]byte(storeBucket))
	tracelog.Info("storesEmailDAO", "getAllStores", "Extracted Stores : " + storesMap)

	emails := make([]string, 0, len(storesMap))

	for _, value := range storesMap {
		//remove empty or unknown emails
		if (value != nil && value != "@") {
			emails = append(emails, value)
		}
	}
	return emails
}

func closeDb() {
	// Close
	boltdbboilerplate.Close()
}


/*	// Exampless
	// Put
	err := boltdbboilerplate.Put([]byte(storeBucket), []byte("ownerKey"), []byte("username"))

	// Get owner
	value := boltdbboilerplate.Get([]byte(storeBucket), []byte("ownerKey"))

	// Delete
	err = boltdbboilerplate.Delete([]byte(storeBucket), []byte("ownerKey"))

	// Insert two key/value
	err = boltdbboilerplate.Put([]byte(storeBucket), []byte("key1"), []byte("value1"))
	err = boltdbboilerplate.Put([]byte(storeBucket), []byte("key2"), []byte("value2"))



	// Get all key/value pairs
	pairs := boltdbboilerplate.GetAllKeyValues([]byte(storeBucket))
	// pairs = [{Key:key1, Value:value1}, {Key: key2, Value:value2}]

	// Close
	boltdbboilerplate.Close()

*/

