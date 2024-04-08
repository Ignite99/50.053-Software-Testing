package html_logger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type HTMLLogger struct {
	columnNum      int
	rowNum         int
	headerFilePath string
	footerFilePath string
	projectType    string
	outputFilePath string
	outputFileName string
	outputFile     *os.File
}

// Constructor
func NewHTMLLogger(outputFilePath string, outputFileName string, projectType string, outputFile *os.File) *HTMLLogger {
	return &HTMLLogger{
		rowNum:         0,
		columnNum:      0,
		headerFilePath: "./HTML_Logger/formats/header.html",
		footerFilePath: "./HTML_Logger/formats/footer.html",
		projectType:    projectType,
		outputFilePath: outputFilePath,
		outputFileName: outputFileName,
		outputFile:     outputFile,
	}
}

func (logger *HTMLLogger) CreateFile() error {
	headerFilePath := logger.headerFilePath

	fmt.Println("Header file path " + headerFilePath)

	headerFile, err := os.Open(headerFilePath)
	if err != nil {
		return errors.New("[ERROR] header file path failed to open")
	}

	outputFile, err := os.Create(logger.outputFilePath + logger.outputFileName)
	if err != nil {
		return errors.New("[ERROR] creating output file failed")
	}
	logger.outputFile = outputFile

	headerBytes, err := ioutil.ReadAll(headerFile)
	if err != nil {
		return errors.New("[ERROR] reading header file failed")
	}

	if _, err := outputFile.Write(headerBytes); err != nil {
		return errors.New("[ERROR] writing header content to output file")
	}

	switch logger.projectType {
	case "DJANGO":
		outputFile.WriteString("<p><b>Fuzz Target: </b>Django</p>")
	case "COAP":
		outputFile.WriteString("<p><b>Fuzz Target: </b>CoAP</p>")
	case "BLE":
		outputFile.WriteString("<p><b>Fuzz Target: </b>BLE</p>")
	}

	// TODO - display any other overall stats

	return nil
}

func (logger *HTMLLogger) AddText(style string, text string) {
	logger.outputFile.WriteString(fmt.Sprintf("<p style=\"%s;\">%s</p>\n", style, text))
}

func (logger *HTMLLogger) CreateTableHeadings(style string, columnNames []string) {
	logger.columnNum = len(columnNames)
	logger.outputFile.WriteString(" <table>\n")
	logger.outputFile.WriteString(fmt.Sprintf(" <tr style=\"%s;\">\n", style))
	logger.outputFile.WriteString("<th>Test no.</th>\n")

	for i := 0; i < logger.columnNum; i++ {
		logger.outputFile.WriteString(fmt.Sprintf("<th>%s</th>\n", columnNames[i]))
	}

	logger.outputFile.WriteString("</tr>\n")
}

func (logger *HTMLLogger) AddRow(row []string) {
	rowSize := len(row)

	logger.rowNum++

	if rowSize != logger.columnNum {
		fmt.Println("@HTMLLogger: Invalid number of columns!")
		return
	}

	logger.outputFile.WriteString("<tr>\n")
	logger.outputFile.WriteString(fmt.Sprintf("<th>%d</th>\n", logger.rowNum))

	for i := 0; i < rowSize; i++ {
		logger.outputFile.WriteString(fmt.Sprintf("<th>%s</th>\n", row[i]))
	}
	logger.outputFile.WriteString("</tr>\n")
}

func (logger *HTMLLogger) AddRowWithStyle(style string, row []string) {
	rowSize := len(row)

	logger.rowNum++

	if rowSize != logger.columnNum {
		fmt.Println("@HTMLLogger: Invalid number of columns!")
		return
	}

	logger.outputFile.WriteString("<tr>\n")
	logger.outputFile.WriteString(fmt.Sprintf(`<th style="%s">%d</th>`, style, logger.rowNum))

	for i := 0; i < rowSize; i++ {
		logger.outputFile.WriteString(fmt.Sprintf(`<th style="%s">%s</th>`, style, row[i]))
	}
	logger.outputFile.WriteString("</tr>")
}

// Close with footer (at least thats what i think it does - nic)
func (logger *HTMLLogger) CloseFile(footerFilePath string) error {
	footerFile, err := os.Open(footerFilePath)
	if err != nil {
		return fmt.Errorf("@HTMLLogger: Footer file not found! %v", err)
	}
	defer footerFile.Close()

	footerBytes, err := ioutil.ReadAll(footerFile)
	if err != nil {
		return fmt.Errorf("error reading footer file: %v", err)
	}
	if _, err := logger.outputFile.Write(footerBytes); err != nil {
		return fmt.Errorf("error writing footer content to output file: %v", err)
	}

	if err := logger.outputFile.Close(); err != nil {
		return fmt.Errorf("error closing output file: %v", err)
	}

	return nil
}
