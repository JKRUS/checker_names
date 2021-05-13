package names

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"secret_path/checkers"
	"secret_path/maps"
	"strings"
)

type packageNamesChecker struct {
	packageImportPath maps.String
}

func NewPackageNamesChecker() checkers.Checker {
	return &packageNamesChecker{
		packageImportPath: make(maps.String),
	}
}

func (c *packageNamesChecker) Setup(dir string, _ string) error {
	dir = dir + "/"
	c.packageImportPath.Add(dir, filepath.Base(dir))
	return nil
}

func (c *packageNamesChecker) Check(fileName string) checkers.Messages {
	var output checkers.Messages
	if strings.HasSuffix(fileName, "_test.go") {
		return nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.AllErrors)
	if err != nil {
		return output.Append(fileName, fmt.Sprintf("Error: FileName checker: %s", err),
			checkers.NoLine, checkers.NoColumn)
	}
	currentPackageName := file.Name.String()
	if currentPackageName == "main" {
		return nil
	}
	pathRealFileName, realFileName := filepath.Split(fileName)
	wantPackageName := c.packageImportPath[pathRealFileName]

	err = checkPackageNames(wantPackageName, currentPackageName, realFileName)
	if err != nil {
		return output.Append(fileName, err.Error(), checkers.NoLine, checkers.NoColumn)
	}
	errors, ok := checkIdentName(file, fset)
	if !ok {
		for _, er := range errors {
			return output.Append(fileName, er, checkers.NoLine, checkers.NoColumn)
		}
	}
	return output
}

func checkPackageNames(packageDirName string, packageName string, fileName string) error {
	err := checkDirName(packageDirName)
	if err != nil {
		return err
	}
	err = checkPackageName(packageName)
	if err != nil {
		return err
	}
	if packageDirName != packageName {
		return fmt.Errorf("the package name %q must match the name of the target directory%q",
			packageName, packageDirName)
		//for test
		//return fmt.Errorf("the package name must match the name of the target directory")
	}
	err = checkFileName(fileName)
	if err != nil {
		return err
	}
	return nil
}

func checkDirName(s string) error {
	varMatcher := regexp.MustCompile(`[^a-z]`)
	if len(varMatcher.FindAllString(s, -1)) > 0 {
		return fmt.Errorf("the name %q contains invalid characters. (Use only \"a-z\")", s)
		//for test
		//return fmt.Errorf("the name contains invalid characters. (Use only \"a-z\")")

	}
	return nil
}

func checkPackageName(s string) error {
	varMatcher := regexp.MustCompile(`[^a-z]`)
	if len(varMatcher.FindAllString(s, -1)) > 0 {
		return fmt.Errorf("the name %q contains invalid characters. (Use only \"a-z\")", s)
		//for test
		//return fmt.Errorf("The name contains invalid characters. (Use only \"a-z\")")
	}
	return nil
}

func checkFileName(s string) error {
	//delete the suffix ".go"
	s = strings.TrimSuffix(s, ".go")
	varMatcher := regexp.MustCompile(`[^a-z_0-9]`)
	if len(varMatcher.FindAllString(s, -1)) > 0 {
		return fmt.Errorf("the name %q contains invalid characters. (Use only \"a-z\", \"_\", \"0-9\")", s)
		//for test
		//return fmt.Errorf("the name contains invalid characters. (Use only \"a-z\", \"_\", \"0-9\")")
	}
	return nil
}

func checkIdentName(file *ast.File, fset *token.FileSet) ([]string, bool) {
	var (
		errors     []string
		err        = true
		varMatcher = regexp.MustCompile(`\W`)
	)

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if len(varMatcher.FindAllString(x.String(), -1)) > 0 {
				errors = append(errors, fmt.Sprintf("the variable name %q in line %v contains invalid characters. (Use only \"a-z\", \"A-Z\", \"0-9\", \"_\")",
					x.Name, fset.Position(n.Pos())))
				//for test
				//errors = append(errors,
				//fmt.Sprintf("the variable name in line contains invalid characters. (Use only \"a-z\", \"A-Z\", \"0-9\", \"_\")"))
				err = false
			}
		}
		return true
	})
	return errors, err
}
