package filewriter

import (
	"os"
	"os/exec"
)

// GoFunc represents a Go function in our generated file
type GoFunc struct {
	Name   string
	Params string
	Body   string
}

// WriteFunc writes a Go function to a file
func WriteGoFunc(file *os.File, f GoFunc) {
	file.WriteString("func " + f.Name + "(" + f.Params + ") {\n")
	file.WriteString(f.Body)
	file.WriteString("\n}\n\n")
}

// deletes and recreates a new file
func ResetFile(filepath string) *os.File {
	_ = os.Remove(filepath)
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	return file
}

// sets package name in a .go file
func SetPackageName(file *os.File, name string) {
	file.WriteString("package " + name + "\n\n")
}

// sets import statements in a .go file
func SetImports(file *os.File, imports []string) {
	file.WriteString("import (\n")
	for _, imp := range imports {
		file.WriteString("\"" + imp + "\"\n")
	}
	file.WriteString(")\n\n")
}

// runs gofmt on the generated file
func RunGoFmt(file *os.File) {
	_ = file.Close()
	cmd := exec.Command("gofmt", "-w", file.Name())
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
