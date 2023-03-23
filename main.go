package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"lead-generator/storage"
	"os"

	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/blank"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/genExports"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/scrape"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/search"
	"github.com/tealeg/xlsx"
)

var (
	CONTAINER_NAME, CONTAINER_URL, ACCOUNT, ACCESS_KEY string
	RUN_ID                                             string
	leadsFile                                          = "Leads.csv"
)

func main() {
	fmt.Println("Starting...")
	S, err := storage.NewBlobStorageConn(CONTAINER_URL, ACCOUNT, ACCESS_KEY, CONTAINER_NAME)
	if err != nil {
		fmt.Printf("Failed to Connect Blob Storage Error: %s\n", err.Error())
		os.Exit(1)
	}
	excel_body := bytes.Buffer{}
	excel_body_file := fmt.Sprintf("%s_excel_body", RUN_ID)
	if err = S.Download(excel_body_file, &excel_body); err != nil {
		fmt.Printf("Failed to download Excel Body File: %s from Blob Storage Error: %s\n", excel_body_file, err.Error())
		os.Exit(1)
	}

	excelBytes := excel_body.Bytes()
	excelFile, err := xlsx.OpenBinary(excelBytes)
	if err != nil {
		fmt.Printf("Failed to open excel file: %s", err.Error())
		os.Exit(1)
	}

	sheet := excelFile.Sheets[0]

	for _, cell := range sheet.Rows[0].Cells {
		blank.Header = append(blank.Header, cell.Value)
	}

	var errur error
	scrape.File, errur = os.Create(leadsFile)
	if errur != nil {
		fmt.Printf("Can't create Leads.csv File: %s\n", errur.Error())
	}
	scrape.Writer = csv.NewWriter(scrape.File)

	headers := []string{
		"Company Name",
		"Company Location",
		"Company Type",
		"Company Description",
	}

	scrape.Writer.Write(headers)
	scrape.Writer.Flush()
	defer scrape.File.Close()

	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]

		search.Lead = row.Cells[0].Value
		search.StateAbb = row.Cells[1].Value

		switch {
		case search.Lead == "" && search.StateAbb != "":
			genExports.GenExports()
			continue
		case search.Lead != "" && search.StateAbb == "":
			search.SearchThomasnet()
			scrape.ScrapeWebsite()
			continue
		case search.Lead != "" && search.StateAbb != "":
			switch search.StateAbb {
			case "TX", "TX - N", "TX - S":
				blank.TXstate()
			case "CA", "CA - N", "CA - S":
				blank.CAstate()
			case "MA", "MA - E", "MA - W":
				blank.MAstate()
			case "NJ", "NJ - N", "NJ - S":
				blank.NJstate()
			case "NY", "NY - M", "NY - U":
				blank.NYstate()
			case "OH", "OH - N", "OH - S":
				blank.OHstate()
			case "PA", "PA - E", "PA - W":
				blank.PAstate()
			default:
				search.SearchThomasnet()
				scrape.ScrapeWebsite()
			}
			continue
		case search.Lead == "" && search.StateAbb == "":
			return
		}

	}

	// Open the leadsFile and handle any errors that occur
	outputFile, err := os.Open(leadsFile)
	if err != nil {
		fmt.Println(err)
	}

	// Defer closing the outputFile until the function returns
	defer outputFile.Close()

	var outputBuffer bytes.Buffer
	// Copy the contents of the outputFile into the outputBuffer and handle any errors that occur
	if _, err := io.Copy(&outputBuffer, outputFile); err != nil {
		fmt.Println(err)
		return
	}

	// Upload the leadsFile to Azure using the S.Upload method and handle any errors that occur
	if err := S.Upload(leadsFile, &outputBuffer); err != nil {
		fmt.Printf("Failed to upload output file to Azure: %s\n", err.Error())
		return
	}

}
