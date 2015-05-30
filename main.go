package main

import (
	"fmt"
	"github.com/mefellows/mirror/command"
	_ "github.com/mefellows/mirror/filesystem/fs"
	_ "github.com/mefellows/mirror/filesystem/remote"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	cli := cli.NewCLI("mirror", Version)
	cli.Args = os.Args[1:]
	cli.Commands = command.Commands

	exitStatus, err := cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(exitStatus)

	/*
			//org := map[string]string{"foo": "bar"}

			network := new(bytes.Buffer)
			enc := gob.NewEncoder(network)
			dec := gob.NewDecoder(network)
			var tree fsys.FileTree
			treeimpl := fsys.StdFileSystemTree{}
			treeimpl.StdFile = fsys.File{FileName: "test"}
			tree = &treeimpl

			store(&tree)

			var loadedTree fsys.FileTree
			load(&loadedTree)

			fmt.Println(loadedTree) // bar
			fmt.Println("------------------------")

			err = enc.Encode(&tree)
			if err != nil {
				fmt.Printf("encode error:", err)
			}
			//x := fsys.StdFileSystemTree{}
			//gob.Register(x)

			//var tree2 fsys.FileTree
			var tree2 interface{}
			//var tree2 fsys.StdFileSystemTree
			err = dec.Decode(&tree2)
			if err != nil {
				fmt.Printf("decode error:", err)
			}
			fmt.Printf("Result: %v\n", tree2)
			fmt.Printf("Reflection value %v\n", reflect.ValueOf(tree2))
			fmt.Printf("Reflection type %v\n", reflect.TypeOf(tree2))

		}
		func store(data interface{}) {
			m := new(bytes.Buffer)
			enc := gob.NewEncoder(m)

			err := enc.Encode(data)
			if err != nil {
				panic(err)
			}

			err = ioutil.WriteFile("dep_data", m.Bytes(), 0600)
			if err != nil {
				panic(err)
			}
		}

		func load(e interface{}) {
			n, err := ioutil.ReadFile("dep_data")
			if err != nil {
				panic(err)
			}

			p := bytes.NewBuffer(n)
			dec := gob.NewDecoder(p)

			err = dec.Decode(e)
			if err != nil {
				panic(err)
			}
		}
	*/
}
