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
	"text/tabwriter"
	"time"
)

var (
	user = os.Getenv("SELECTEL_USER")
	key  = os.Getenv("SELECTEL_PASS")
)

type File struct {
	Parent   *File
	Files    []*File
	Name     string
	Size     uint64
	Modified time.Time
	Type     string
	IsFile   bool
}

func (p *File) Find(name string,object storage.ObjectInfo) *File {
	for _, ptn := range p.Files {
		if ptn.Name == name {
			return ptn
		}
	}
	new := &File{
		Name:     name,
		Parent:   p,
		Size:     object.Size,
		Modified: object.LastModified,
		Type:     object.ContentType,
	}
	p.Files = append(p.Files, new)
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
	FileTable(p.Files)
	return p.Files[unswer(len(p.Files),"Select object: ")]
}

var bucket *File = &File{
	Parent: nil,
	Files:  []*File{},
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
	ContainerTable(containers)
	container := api.Container(containers[unswer(len(containers),"Select container: ")].Name)
	objects, err := container.ObjectsInfo()
	if err != nil {
		log.Fatal(err)
	}
	for _, object := range objects {
		paths := strings.Split(object.Name, "/")
		current := bucket
		for index, p := range paths {
			current = current.Find(p,object)
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
	if err := ioutil.WriteFile(stringUnswer("Write filename: "), res, 755); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
}
func unswer(max int,caption string) int {
	fmt.Print(caption)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	code, err := strconv.Atoi(text)
	if err != nil || code > max || code < 0 {
		fmt.Println("Wrong input, try again")
		return unswer(max,caption)
	}
	return code - 1
}
func stringUnswer(caption string) string {
	fmt.Print(caption)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
func FileTable(files []*File) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	table := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Println()
	fmt.Fprintf(table, format, "#", "Filename", "Size", "Modified", "Type")
	fmt.Fprintf(table, format, "-", "--------", "----", "--------", "----")
	for index, file := range files {
		fmt.Fprintf(table, format, index+1, file.Name, file.Size, file.Modified, file.Type)
	}

	table.Flush()
	fmt.Println()
}
func ContainerTable(containers []storage.ContainerInfo) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	table := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Println()
	fmt.Fprintf(table, format, "#", "Container", "Type", "Objects", "Size")
	fmt.Fprintf(table, format, "-", "---------", "----", "-------", "----")
	for index, container := range containers {
		fmt.Fprintf(table, format, index+1, container.Name, container.Type, container.ObjectCount, container.BytesUsed)
	}

	table.Flush()
	fmt.Println()
}
