package main

import (
	"encoding/xml"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := cli.NewApp()
	app.Name = "wadoku2mongo"
	app.Usage = "Converts a wadoku XML dump into a mongodb collection"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Nils Hasenbanck",
			Email: "nils@hasenbanck.de",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "data/wadoku.xml",
			Usage: "The path to the file of the wadoku XML dump",
		},
		cli.StringFlag{
			Name:  "mongodb",
			Value: "127.0.0.1",
			Usage: "The server connection string to the mongodb server",
		},
	}

	app.Action = func(c *cli.Context) {
		export(c.String("file"), c.String("mongodb"))
	}

	app.Run(os.Args)
}

func export(file string, connection string) {
	var wadokufile *os.File
	var xmldata []byte
	var err error

	if wadokufile, err = os.Open(file); err != nil {
		log.Fatal("Can't open wadoku XML file: " + err.Error())
	}
	defer wadokufile.Close()

	if xmldata, err = ioutil.ReadAll(wadokufile); err != nil {
		log.Fatal("Can't read wadoku XML file: " + err.Error())
	}

	dict := XMLDict{}

	if err = xml.Unmarshal([]byte(xmldata), &dict); err != nil {
		log.Fatal("Can't unmarshal xmldata: " + err.Error())
		return
	}

	if err = saveIntoMongo(dict, connection); err != nil {
		log.Fatal("Can't save entries into mongodb: " + err.Error())
	}
}
