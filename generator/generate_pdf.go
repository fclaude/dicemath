package generator

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
)

const fileName = "worksheet.tex"
const fileNamePdf = "worksheet.pdf"

func GeneratePDF(name string) ([]byte, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "worksheet")
	if err != nil {
		return nil, err
	}
	//defer os.RemoveAll(dir)

	// the template is hardcoded, can ignore the error
	fileContent, err := template.New("sheet").Delims("[[", "]]").Parse(sheetTemplate)
	if err != nil {
		return nil, err
	}

	filePath := path.Join(dir, fileName)
	fileWriter, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	err = fileContent.Execute(fileWriter, generateGroups(name))
	fileWriter.Close()
	if err != nil {
		return nil, err
	}

	err = exec.Command("pdflatex", "-aux-directory="+dir, "-output-directory="+dir, filePath).Run()
	if err != nil {
		return nil, err
	}

	pdf := path.Join(dir, fileNamePdf)
	return ioutil.ReadFile(pdf)
}
