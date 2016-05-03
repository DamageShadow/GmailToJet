package main

import "github.com/bobintornado/boltdb-boilerplate"

import (
	"github.com/goinggo/tracelog"
	"fmt"
	"strings"

)

var (
	storeBucket string = "storesEmails"
	defaultEmail string = "@"
	dbFile string = "./db/storesDatabase.boltdb"
)

func init() {
	// Init buckets and create file
	buckets := []string{storeBucket}
	/*

	f, err := os.Create(dbFile)
	if err != nil {
		panic(err)
	}

	err = os.Chmod(dbFile, 0777)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	f.Close()
	*/

	err := boltdbboilerplate.InitBolt("./db/storesDatabase.boltdb", buckets)
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
		for _, i := range newcomers {
			insertStore(i)
		}
	}

	if (obsolete != nil) {
		for _, i := range obsolete {
			insertStore(i)
		}
	}
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

	if (getStoreEmail(store) != "") {
		deleteStore(store)
	}

	insertStoreWithEmail(store, email)
}

func getStoreEmail(store string) (string) {
	value := boltdbboilerplate.Get([]byte(storeBucket), []byte(store))

	return string(value)
}

func getAllStores() []string {
	var stores []string
	// Get all keys (stores)
	storesBytes := boltdbboilerplate.GetAllKeys([]byte(storeBucket))

	for _, i := range storesBytes {
		stores = append(stores, string(i))
	}

	tracelog.Info("storesEmailDAO", "getAllStores", "Extracted Stores : " + strings.Join(stores[:], ","))
	return stores
}

func getAllStoresEmails() []string {
	// Get all values (emails)
	storesMap := boltdbboilerplate.GetAllKeyValues([]byte(storeBucket))

	emails := make([]string, 0, len(storesMap))

	tracelog.Info("storesEmailDAO", "getAllStores", "Extracted Stores Emails : " + strings.Join(emails[:], ","))

	for _, value := range emails {
		//remove empty or unknown emails
		if (value != "" && value != "@") {
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

