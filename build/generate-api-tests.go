// +build ignore

package main

import (
	"strings"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "text/template"
    "go/build"
)

const subPackage = "github.com/32leaves/ruruku/pkg/api/v1/test"

type TestFunc struct {
    TestFunc string
    TestName string
}

type Cfg struct {
    Pkg string
    Funcs []TestFunc
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("usage <test-package> <src-file>")
        os.Exit(1)
    }
    tp := os.Args[1]
    srcfn := os.Args[2]

    p, err := build.Import(subPackage, "", build.FindOnly)
	if err != nil {
		panic(err)
	}

    set := token.NewFileSet()
    packs, err := parser.ParseDir(set, p.Dir, nil, 0)
    if err != nil {
        fmt.Println("Failed to parse package:", err)
        os.Exit(1)
    }

    funcs := make([]TestFunc, 0)
    for _, pack := range packs {
        for _, f := range pack.Files {
            for _, d := range f.Decls {
                if fn, isFn := d.(*ast.FuncDecl); isFn && strings.HasPrefix(fn.Name.String(), "RunTest") {
                    nme := fn.Name.String()
                    funcs = append(funcs, TestFunc{
                        TestName: nme[3:],
                        TestFunc: nme,
                    })
                }
            }
        }
    }

    tpl, err := template.New("fn").Parse(`
{{- $pkg := .Pkg -}}
package {{ .Pkg }}

// BEWARE: This file was auto-generated using go generate. DO NOT EDIT.

import (
    "testing"
    apitests "github.com/32leaves/ruruku/pkg/api/v1/test"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.SetLevel(log.WarnLevel)
}

{{ range .Funcs }}
func {{ .TestName }}(t *testing.T) {
    srv := newTestServer()
    apitests.{{ .TestFunc }}(t, srv)
}
{{ end -}}
`)
    if err != nil {
        panic(err)
    }

    f, err := os.OpenFile(fmt.Sprintf("gen_%s", srcfn), os.O_CREATE | os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    tpl.Execute(f, Cfg{
        Pkg: tp,
        Funcs: funcs,
    })
    f.Sync()
}
