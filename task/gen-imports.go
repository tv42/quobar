// +build task

// Generate side effect only import statements, usually used for
// registering plugins.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"go/build"

	"github.com/kisielk/gotool"
)

var (
	genOutput  = flag.String("o", "", "output path")
	genPackage = flag.String("package", os.Getenv("GOPACKAGE"), "Go package name")
)

var gen = template.Must(template.New("gen").Parse(`package {{.Package}}

import (
{{range .Imports}}{{"\t"}}_ "{{.}}"
{{end}})
`))

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s -o PATH PACKAGE..\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func expandPackages(spec []string) ([]string, error) {
	// expand "..."
	paths := gotool.ImportPaths(spec)

	var r []string
	for _, path := range paths {
		pkg, err := build.Import(path, ".", 0)
		if _, ok := err.(*build.NoGoError); ok {
			// directory with no Go source files in it
			continue
		}
		if err != nil {
			return nil, err
		}
		if pkg.ImportPath == "" {
			return nil, fmt.Errorf("no import path found: %v", path)
		}
		r = append(r, pkg.ImportPath)
	}
	return r, nil
}

// avoidCycles removes import paths that point to the directory
// containing filename. This lets the generated file be placed in a
// package it is referring to via importpath/...
func avoidCycles(imports []string, filename string) ([]string, error) {
	curDir, err := filepath.Abs(".")
	if err != nil {
		return imports, err
	}
	pkg, err := build.Import(".", path.Dir(path.Join(curDir, filename)), build.FindOnly)
	if err != nil {
		return imports, err
	}
	r := imports[:0]
	for _, p := range imports {
		if p != pkg.ImportPath {
			r = append(r, p)
		}
	}
	return r, nil
}

func process(dst string, imports []string) error {
	dir := filepath.Dir(dst)
	tmp, err := ioutil.TempFile(dir, "temp-gen-import-all-")
	if err != nil {
		return err
	}
	closed := false
	removed := false
	defer func() {
		if !closed {
			// silence errcheck
			_ = tmp.Close()
		}
		if !removed {
			// silence errcheck
			_ = os.Remove(tmp.Name())
		}
	}()

	imports, err = expandPackages(imports)
	if err != nil {
		return fmt.Errorf("listing packages: %v", err)
	}

	imports, err = avoidCycles(imports, *genOutput)
	if err != nil {
		return fmt.Errorf("trying to avoid cycles: %v", err)
	}

	type state struct {
		Package string
		Imports []string
	}
	s := state{
		Package: *genPackage,
		Imports: imports,
	}
	if err := gen.Execute(tmp, s); err != nil {
		return fmt.Errorf("template error: %v", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("cannot write temp file: %v", err)
	}
	closed = true

	if err := os.Rename(tmp.Name(), *genOutput); err != nil {
		return fmt.Errorf("cannot finalize file: %v", err)
	}
	removed = true

	return nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if *genOutput == "" {
		flag.Usage()
		os.Exit(2)
	}
	if *genPackage == "" {
		log.Fatal("$GOPACKAGE must be set or -package= passed")
	}

	if err := process(*genOutput, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
