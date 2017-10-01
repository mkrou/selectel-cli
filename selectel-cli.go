package main

import (
	"github.com/ernado/selectel/storage"
	"log"
	"strconv"
	"fmt"
	"bufio"
	"os"
	"strings"
	"io/ioutil"
)

var (
	user = os.Getenv("SELECTEL_USER")
	key  = os.Getenv("SELECTEL_PASS")
)

type File struct {
	Parent *File
	Path   []*File
	Name   string
	IsFile bool
}

func (p *File) Find(name string) *File {
	for _, pth := range p.Path {
		if pth.Name == name {
			return pth
		}
	}
	new := &File{
		Name:   name,
		Parent: p,
	}
	p.Path = append(p.Path, new)
	return new
}

func (p *File) Fullname() string {
	current := p.Parent
	fullname := p.Name
	for current != nil {
		fullname = current.Name + "/" + fullname
		current = current.Parent
	}
	return fullname
}
func (p *File) FindObject(api storage.ContainerAPI) storage.ObjectAPI {
	return api.Object(p.Fullname())
}
func (p *File) Select() *File {
	if p.IsFile {
		return nil
	}
	for index, pth := range p.Path {
		fmt.Println(strconv.Itoa(index) + ": " + pth.Name)
	}
	return p.Path[unswer(len(p.Path))]
}

var bucket *File = &File{
	Parent: nil,
	Path:   []*File{},
	Name:   "",
}

func main() {
	api, err := storage.New(user, key)
	if err != nil {
		log.Fatal(err)
	}
	containers, err := api.ContainersInfo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Select container:")
	for index, container := range containers {
		fmt.Println(strconv.Itoa(index) + ": " + container.Name)
	}
	container := api.Container(containers[unswer(len(containers))].Name)
	objects, err := container.ObjectsInfo()
	if err != nil {
		log.Fatal(err)
	}
	for _, object := range objects {
		paths := strings.Split(object.Name, "/")
		current := bucket
		for index, p := range paths {
			current = current.Find(p)
			if len(paths) == index+1 {
				current.IsFile = true
			}
		}
	}
	current := bucket
	var file *File
	for current != nil {
		file = current
		current = current.Select()
	}
	object := file.FindObject(container)
	res, err := object.Download()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Write a filename:")

	if err := ioutil.WriteFile(stringUnswer(), res, 755); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
}
func unswer(max int) int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	code, err := strconv.Atoi(text)
	if err != nil || code > max || code < 0 {
		fmt.Println("Wrong input, try again")
		return unswer(max)
	}
	return code
}
func stringUnswer() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
