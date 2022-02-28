package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/hnlq715/struct2interface"
	"github.com/spf13/cobra"
)

func main() {
	var (
		dir, CopyDocs, PkgName, IfaceComment, IfaceName, Comment, Output string
		copyDocs, InPackage, CopyTypeDoc                                 bool
	)

	root := &cobra.Command{
		Use: "struct2interface",
		Run: func(cmd *cobra.Command, args []string) {
			// Workaround because jessevdk/go-flags doesn't support default values for boolean flags
			copyDocs = CopyDocs == "true"

			if IfaceComment == "" {
				IfaceComment = fmt.Sprintf("%s ...", IfaceName)
			}

			if Comment == "" {
				Comment = "Code generated by struct2interface; DO NOT EDIT."
			}

			filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					log.Fatal(err)
				}

				if d.IsDir() {
					return nil
				}

				if strings.HasPrefix(filepath.Base(path), "interface_") {
					return nil
				}
				if strings.HasPrefix(filepath.Base(path), "mock_") {
					return nil
				}
				if !strings.HasSuffix(filepath.Base(path), ".go") {
					return nil
				}

				result, err := struct2interface.Make([]string{path}, Comment, PkgName, IfaceName, IfaceComment, copyDocs, CopyTypeDoc)
				if err != nil {
					log.Fatal("struct2interface.Make failed,", err.Error(), path)
				}

				if len(result) == 0 {
					return nil
				}

				if InPackage {
					dir := filepath.Dir(path)
					Output = filepath.Join(dir, "interface_"+filepath.Base(path))
				}

				if Output == "" {
					fmt.Println(string(result))
				} else {
					ioutil.WriteFile(Output, result, 0644)
				}
				return nil
			})

		},
	}
	type cmdlineArgs struct {
		Dir          string `short:"d" long:"dir" description:"Go source file to read, either dir or glob"  `
		IfaceName    string `short:"i" long:"iface" description:"Name of the generated interface" `
		PkgName      string `short:"p" long:"pkg" description:"Package name for the generated interface" required:"true"`
		IfaceComment string `short:"y" long:"iface-comment" description:"Comment for the interface, default is '// <iface> ...'"`
		InPackage    string `short:"n" long:"inpackage" description:"Write interface into the same package"`

		CopyDocs string `long:"doc" description:"Copy docs from methods" option:"true" option:"false" default:"true"`
		copyDocs bool

		CopyTypeDoc bool   `short:"D" long:"type-doc" description:"Copy type doc from struct"`
		Comment     string `short:"c" long:"comment" description:"Append comment to top, default is '// Code generated by ifacemaker; DO NOT EDIT.'"`
		Output      string `short:"o" long:"output" description:"Output file name. If not provided, result will be printed to stdout."`
	}

	root.Flags().StringVarP(&dir, "dir", "d", "", "Go source file to read, either dir or glob")
	root.Flags().StringVarP(&IfaceName, "iface", "i", "", "Name of the generated interface")
	root.Flags().StringVarP(&PkgName, "pkg", "p", "", "Package name for the generated interface")
	root.Flags().StringVarP(&IfaceComment, "iface-comment", "y", "", "Comment for the interface, default is '// <iface> ...'\"")
	root.Flags().BoolVarP(&InPackage, "inpackage", "P", false, "Write interface into the same package")
	root.Flags().StringVarP(&CopyDocs, "doc", "", "true", "Copy docs from methods")
	root.Flags().BoolVarP(&CopyTypeDoc, "type-doc", "D", true, "Copy type doc from struct")
	root.Flags().StringVarP(&Comment, "comment", "c", "", "Append comment to top, default is '// Code generated by ifacemaker; DO NOT EDIT.'\"")
	root.Flags().StringVarP(&Output, "output", "o", "", "Output file name. If not provided, result will be printed to stdout.")
	root.Execute()

}
