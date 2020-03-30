package generator

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"text/template"
)

const fileName = "worksheet.tex"
const fileNamePdf = "worksheet.pdf"

func GeneratePDF(name, operation string) ([]byte, error) {
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
	err = fileContent.Execute(fileWriter, generateGroups(name, operation))
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

type exercise struct {
	A int
	B int
}

func genExercise() exercise {
	return exercise{A: rand.Intn(6) + 1, B: rand.Intn(6) + 1}
}

const N = 5

type exerciseGroup struct {
	Group1    []exercise
	Group2    []exercise
	Group3    []exercise
	Group4    []exercise
	Name      string
	Operation string
	Separator string
}

func valid(seen map[exercise]int, ex exercise, operation string) bool {
	switch operation {
	case "-":
		return ex.A >= ex.B && seen[ex] == 0
	case "/":
		return ex.A%ex.B == 0 && seen[ex] == 0
	}
	return seen[ex] == 0
}

func generateGroups(name, operation string) exerciseGroup {
	seen := make(map[exercise]int, 36)
	randomGroup := func() []exercise {
		result := make([]exercise, N)
		for i := range result {
			ex := genExercise()
			for !valid(seen, ex, operation){
				ex = genExercise()
			}
			seen[ex]++
			result[i] = ex
		}
		return result
	}
	separator := ""
	if operation == "/" {
		separator = "\\newpage"
	}
	group1 := randomGroup()
	group2 := make([]exercise, 0, 5)
	group3 := randomGroup()
	group4 := make([]exercise, 0, 5)

	if operation != "/" {
		group2 = randomGroup()
		group4 = randomGroup()
	}
	return exerciseGroup{
		Name:      name,
		Operation: operation,
		Separator: separator,
		Group1:    group1,
		Group2:    group2,
		Group3:    group3,
		Group4:    group4,
	}
}
