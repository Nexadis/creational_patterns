package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type config struct {
	addr    string
	dburi   string
	timeout int64
	logFile string
}

func (c *config) String() string {
	return fmt.Sprintf("config {addr:%s, dburi:%s, timeout:%d, logFile:%s}", c.addr, c.dburi, c.timeout, c.logFile)
}

type opt func(c *config)

func SetAddr(addr string) opt {
	return func(c *config) {
		c.addr = addr
	}
}

func SetDBUri(dburi string) opt {
	return func(c *config) {
		c.dburi = dburi
	}
}

func NewConfig(opts ...opt) *config {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Node — интерфейс с методом Clone.
type Node interface {
	Clone() Node
	WriteEntry(w io.Writer, opts WriteOpts) error
}

type File struct {
	Name string
}

type WriteOpts struct {
	// Level определяет уровень вложенности файла
	// для соответствующего сдвига при выводе
	Level int
}

// WriteEntry выводит имя файла с нужным сдвигом.
func (f *File) WriteEntry(w io.Writer, opts WriteOpts) error {
	_, err := fmt.Fprintf(w, "%s%s\n", strings.Repeat("    ", opts.Level), f.Name)
	return err
}

// Clone возвращает копию файла.
func (f *File) Clone() Node {
	// возвращаем указатель на File, так как этот тип поддерживает интерфейс Node
	return &File{
		Name: f.Name,
	}
}

type Folder struct {
	File     // вложенная структура File
	Children []Node
}

// WriteEntry выводит имя директории и её содержимое.
func (f *Folder) WriteEntry(w io.Writer, opts WriteOpts) error {
	err := f.File.WriteEntry(w, opts)
	if err != nil {
		return err
	}

	opts.Level += 1
	for _, v := range f.Children {
		err := v.WriteEntry(w, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

// Clone возвращает копию директории.
func (f *Folder) Clone() Node {
	clone := &Folder{
		Children: make([]Node, len(f.Children)),
	}
	clone.Name = f.Name
	for i, v := range f.Children {
		clone.Children[i] = v.Clone()
	}
	return clone
}

func (f *Folder) String() string {
	var sb strings.Builder
	f.WriteEntry(&sb, WriteOpts{Level: 0})
	return sb.String()
}

type singleInstance struct{}

var (
	singleton *singleInstance
	once      sync.Once
)

func getSingleton() *singleInstance {
	once.Do(
		func() {
			fmt.Println("Initialize singleton")
			singleton = &singleInstance{}
		},
	)
	return singleton
}

func main() {
	border("Singleton")
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("Singleton: %p\n", getSingleton())
		}(i)
	}
	time.Sleep(100 * time.Millisecond)
	border("Prototype")
	folder := &Folder{
		Children: []Node{
			&File{Name: "file1"},
			&File{Name: "file2"},
			&Folder{
				Children: []Node{
					&File{Name: "file3"},
					&File{Name: "file4"},
				},
				File: File{Name: "subfolder"},
			},
		},
		File: File{Name: "root"},
	}
	fmt.Printf("Original:\n%s\n", folder)

	clone := folder.Clone()
	fmt.Printf("Clone:\n%s\n", clone)
	border("FuncOpts")
	config := NewConfig(SetAddr("Some addr"), SetDBUri("db uri"))
	fmt.Println(config)
}

func border(name string) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\t", name)
	fmt.Println(strings.Repeat("=", 80))
}
