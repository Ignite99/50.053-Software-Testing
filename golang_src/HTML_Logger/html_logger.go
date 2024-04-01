package html_logger

import (
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
		headerFilePath: "./src/HTML_Logger/formats/header.html",
		footerFilePath: "./src/HTML_Logger/formats/footer.html",
		projectType:    projectType,
		outputFilePath: outputFilePath,
		outputFileName: outputFileName,
		outputFile:     outputFile,
	}
}

func (logger *HTMLLogger) CreateFile(headerFilePath string) error {
	headerFile, err := os.Open(headerFilePath)
	if err != nil {
		return fmt.Errorf("@HTMLLogger: Header file not found! %v", err)
	}
	defer headerFile.Close()

	outputFile, err := os.Create(logger.outputFilePath + logger.outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	logger.outputFile = outputFile
	defer outputFile.Close()

	headerBytes, err := ioutil.ReadAll(headerFile)
	if err != nil {
		return fmt.Errorf("error reading header file: %v", err)
	}
	if _, err := outputFile.Write(headerBytes); err != nil {
		return fmt.Errorf("error writing header content to output file: %v", err)
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
	logger.outputFile.WriteString(fmt.Sprintf(`<th style="%s">%d</th>\n`, style, logger.rowNum))

	for i := 0; i < rowSize; i++ {
		logger.outputFile.WriteString(fmt.Sprintf(`<th style="%s">%s</th>\n`, style, row[i]))
	}
	logger.outputFile.WriteString("</tr>\n")
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
