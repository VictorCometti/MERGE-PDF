package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

const UNIDOC_LICENSE_API_KEY = "af4fe092a5d482884df745b0da2912334ed591cf150c3720ea70f3ee13fd5d37"

func init() {
	err := license.SetMeteredKey(UNIDOC_LICENSE_API_KEY)
	if err != nil {
		common.Log.Error("Erro: %v", err)
	}
}

func main() {
	if len(os.Args) < 4 {
		os.Exit(0)
	}

	outputPath := ""
	inputPaths := []string{}

	for i, arg := range os.Args {
		if i == 0 {
			continue
		} else if i == 1 {
			outputPath = arg
			continue
		}
		inputPaths = append(inputPaths, arg)
	}

	err := margePdf(inputPaths, outputPath)

	if err != nil {
		common.Log.Error("Erro: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Mesclagem completa, olhe seu árquivo: : %s\n", outputPath)

	deleteAfterMarge(inputPaths)
}

func margePdf(inputPaths []string, outputPath string) error {
	pdfWriter := model.NewPdfWriter()

	for _, inputPath := range inputPaths {
		pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
		if err != nil {
			common.Log.Error("Erro: %v", err)
			return err
		}
		defer f.Close()

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			common.Log.Error("Erro: %v", err)
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				common.Log.Error("Erro: %v", err)
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				common.Log.Error("Erro: %v", err)
				return err
			}
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		common.Log.Error("Erro: %v", err)
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		common.Log.Error("Erro: %v", err)
		return err
	}
	return nil
}

func deleteAfterMarge(inputPaths []string) {
	for _, inputPath := range inputPaths {
		err := os.Remove(inputPath)
		if err != nil {
			common.Log.Error("Erro: %v", err)
		}
	}
}
