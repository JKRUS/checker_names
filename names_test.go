package names

import (
	"errors"
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckModuleNames(t *testing.T) {
	msgFileName := "the name contains invalid characters. (Use only \"a-z\", \"_\", \"0-9\")"
	msgDirName := "the name contains invalid characters. (Use only \"a-z\")"
	msgPackageName := msgDirName
	msgDirAndPackageName := "the package name must match the name of the target directory"
	var tests = []struct {
		inputDirName     string
		inputPackageName string
		inputFileName    string
		want             error
	}{
		//test dir name
		//correct
		{"dirname", "dirname", "filename.go", nil},
		//not correct
		{"Dirname", "Dirname", "filename.go", errors.New(msgDirName)},
		{"dirName", "dirName", "filename.go", errors.New(msgDirName)},
		{"dir_name", "dir_name", "filename.go", errors.New(msgDirName)},
		{"dirname123", "dirname123", "filename.go", errors.New(msgDirName)},
		{"dir_name123", "dir_name123", "filename.go", errors.New(msgDirName)},
		{"dirИмя", "dirИмя", "filename.go", errors.New(msgDirName)},
		{"dirn$ame", "dirn$ame", "filename.go", errors.New(msgDirName)},

		//test package name
		//correct
		{"dirname", "dirname", "filename.go", nil},
		//not correct
		{"Dirname", "Dirname", "filename.go", errors.New(msgPackageName)},
		{"dirName", "dirName", "filename.go", errors.New(msgPackageName)},
		{"dir_name", "dir_name", "filename.go", errors.New(msgPackageName)},
		{"dirname123", "dirname123", "filename.go", errors.New(msgPackageName)},
		{"dir_name123", "dir_name123", "filename.go", errors.New(msgPackageName)},
		{"dirИмя", "dirИмя", "filename.go", errors.New(msgPackageName)},
		{"dirn$ame", "dirn$ame", "filename.go", errors.New(msgPackageName)},

		//test dir + package name
		//correct
		{"dirname", "dirname", "filename.go", nil},
		//not correct
		{"dirname", "packagename", "filename.go", errors.New(msgDirAndPackageName)},

		//test file name
		//correct
		{"dirname", "dirname", "filename.go", nil},
		{"dirname", "dirname", "file_name.go", nil},
		{"dirname", "dirname", "filename123.go", nil},
		{"dirname", "dirname", "file_name_123.go", nil},
		//not correct
		{"dirname", "dirname", "Filename.go", errors.New(msgFileName)},
		{"dirname", "dirname", "fileName.go", errors.New(msgFileName)},
		{"dirname", "dirname", "Имяfile.go", errors.New(msgFileName)},
		{"dirname", "dirname", "file&name.go", errors.New(msgFileName)},
		{"dirname", "dirname", "file.name.go", errors.New(msgFileName)},
	}

	for _, test := range tests {
		got := checkPackageNames(test.inputDirName, test.inputPackageName, test.inputFileName)
		var s1, s2 string
		if got == nil {
			s1 = ""
		} else {
			s1 = got.Error()
		}
		if test.want == nil {
			s2 = ""
		} else {
			s2 = test.want.Error()
		}
		if s1 != s2 {
			t.Errorf("checkModuleNames(%q,%q,%q) = %v, need - %v", test.inputDirName,
				test.inputPackageName, test.inputFileName, got, test.want)
		}
	}
}

func TestCheckIdentName(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset
	src := `package foo

	import (
	"fmt"
	"time"
	"flag"
	"time2"
	)
	var (
	//correct
	varname string
	varName string
	Varname string
	var_name string
	varname123 string
	v string
	//not correct
	varИмя string
	)
	func bar(переменнаяName []string) {
	for _, f := range переменнаяName{

		fmt.Println(time.Now())
	}`
	msgVarName := "the variable name in line contains invalid characters. (Use only \"a-z\", \"A-Z\", \"0-9\" \"_\")"
	var (
		testErrors = []string{msgVarName, msgVarName, msgVarName, ""}
		Errs       = make([]string, 4)
		ok         = false
	)

	fileAst, _ := parser.ParseFile(fset, "", src, parser.AllErrors)
	Errs, ok = checkIdentName(fileAst, fset)

	if !ok {
		for i, er := range Errs {
			if er != testErrors[i] {
				t.Errorf("checkIdentName() = %v, need - %v", er, testErrors[i])
			}
		}
	}

}
