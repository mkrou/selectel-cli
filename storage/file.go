package storage

import (
	"github.com/ernado/selectel/storage"
	"time"
	"github.com/earlcherry/selectel-cli/cli"
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

func (p *File) Rows() []cli.Tabler {
	var result []cli.Tabler
	for _, file := range p.Files {
		result = append(result, file)
	}
	return result
}
func (p *File) Find(name string, object storage.ObjectInfo) *File {
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
func (p *File) Row(n int) []interface{} {
	return []interface{}{n, p.Name, p.Size, p.Modified, p.Type,}
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
	selected:=cli.Table(
		"%v\t%v\t%v\t%v\t%v\t\n",
		[]interface{}{"#", "Filename", "Size", "Modified", "Type"},
		[]interface{}{"-", "--------", "----", "--------", "----"},
		p.Rows(),
	)

	return p.Files[selected]
}
