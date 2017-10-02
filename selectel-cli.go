package main

import (
	"log"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"github.com/earlcherry/selectel-cli/storage"
	"github.com/earlcherry/selectel-cli/cli"
)



var (
	user = os.Getenv("SELECTEL_USER")
	key  = os.Getenv("SELECTEL_PASS")
)

var root *storage.File = &storage.File{
	Parent: nil,
	Files:  []*storage.File{},
	Name:   "",
}

func main() {
	api, err := storage.New(user, key)
	if err != nil {
		log.Fatal(err)
	}
	var containers storage.Containers
	containers, err = api.ContainersInfo()
	if err != nil {
		log.Fatal(err)
	}
	container := containers.Select(api)
	objects, err := container.ObjectsInfo()
	if err != nil {
		log.Fatal(err)
	}
	for _, object := range objects {
		paths := strings.Split(object.Name, "/")
		current := root
		for index, p := range paths {
			current = current.Find(p, object)
			if len(paths) == index+1 {
				current.IsFile = true
			}
		}
	}
	current := root
	var file *storage.File
	for current != nil {
		file = current
		current = current.Select()
	}
	object := file.FindObject(container)
	res, err := object.Download()
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(cli.StringUnswer("Write filename: "), res, 755); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
}

