package main

import (
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/urfave/cli"
)

var (
	db *leveldb.DB
)

func main() {
	app := cli.NewApp()
	app.Name = "ldbkv"
	app.Usage = "read key values from leveldb"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p, path",
			Usage: "path to leveldb database",
		},
	}
	app.Before = initStore
	app.Commands = []cli.Command{
		{
			Name:   "get",
			Usage:  "get a value for a key",
			Action: get,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "k, key",
					Usage: "string value of key to print",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "list the keys in the database",
			Action: list,
		},
	}

	// Run the CLI program
	app.Run(os.Args)
}

func initStore(c *cli.Context) (err error) {
	if db, err = leveldb.OpenFile(c.String("path"), nil); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func get(c *cli.Context) (err error) {
	// Ensure the database gets closed
	defer db.Close()

	var data []byte
	if data, err = db.Get([]byte(c.String("key")), nil); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// TODO: probably have to unmarshal here
	fmt.Println(string(data))

	return nil
}

func list(c *cli.Context) (err error) {
	// Ensure the database gets closed
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Println(string(iter.Key()))
	}

	return nil
}
