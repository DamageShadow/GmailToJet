package main

import (
	"log"

	"github.com/boltdb/bolt"
	"encoding/json"
	"fmt"
)

type Response struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func mainBolt() {
	db, err := bolt.Open("../db/blog.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	post := &Response{
		Page : 1,
		Fruits:[]string{"apple", "peach", "pear"}, }

/*	encoded, err := json.Marshal(post)
	fmt.Println(string(encoded))

	res := Response{}

	if err := json.Unmarshal(encoded, &res); err != nil {
		panic(err)
	}

	fmt.Println(res)*/


	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("posts"))
		if err != nil {
			return err
		}
		encoded, err := json.Marshal(post)
		if err != nil {
			return err
		}
		return b.Put([]byte("post1"), encoded)
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		v := b.Get([]byte("post1"))
		res := Response{}
		json.Unmarshal(v, &res)
		fmt.Println(res)
		return nil
	})


	db.Close()
}