package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/blank"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/excel"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/genExports"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/scrape"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/search"
	"github.com/tealeg/xlsx"
)

var (
	accountName   = "leadgeneratorstorageblob"
	accountKey    = "rOjnqusNKyYwLdVeK+v0eKnxN0lyRi5vbPQ0csRDFNWOnpwiutTU0Zg/m2pMn1AyCKcdJIyiJbut+AStkcj6Rg=="
	containerName = "lead-generator"
)

func main() {
	//downloads the excel file from AZSB
	excel.Excel()

	//opens the excel file to read data.
	file, err := xlsx.OpenFile("Lead_Template.xlsx")
	if err != nil {
		fmt.Println("Can't opent the excel file", err.Error())
	}

	//set active sheet
	sheet := file.Sheets[0]

	//iterates over row zero to retrieve all the headers
	for _, cell := range sheet.Rows[0].Cells {
		blank.Header = append(blank.Header, cell.Value)
	}

	var errur error
	scrape.File, errur = os.Create("Leads.csv")
	if errur != nil {
		fmt.Println("Can't create Leads.csv file", errur.Error())
	}
	scrape.Writer = csv.NewWriter(scrape.File)

	headers := []string{
		"Company Name",
		"Company Description",
		"Company Location",
	}

	scrape.Writer.Write(headers)
	scrape.Writer.Flush()
	defer scrape.File.Close()

	//Iterate over each row in the sheet, starting from index 0
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

	// create a new background context
	ctx := context.Background()
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Printf("UPLOADING %s TO %s CONTAINER\n", scrape.BlobName, containerName)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++")

	// flush the writer
	scrape.Writer.Flush()

	// move the file cursor to beginning of the file
	scrape.File.Seek(0, 0)

	// create a buffer and read all contents from the file into it
	buffer := bytes.Buffer{}
	buffer.ReadFrom(scrape.File)

	// create a shared key credential account for azure storage and handle errors
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Println("Can't connect to azure blob", err.Error())
	}

	// Create a pipeline using credentials and set pipeline options.
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Parse the string url and convert it into a url type. Use that URL to create a ServiceURL struct.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	serviceURL := azblob.NewServiceURL(*u, p)

	// Using the service URL, create or reference the named ContainerURL object, and perform operations on the container.
	containerURL := serviceURL.NewContainerURL(containerName)

	// Create a new BlockBlobURL object (interface for working with block blobs).
	blobURL := containerURL.NewBlockBlobURL(scrape.BlobName)

	// Upload the buffer's content to the specified block blob.
	_, erroor := azblob.UploadBufferToBlockBlob(ctx, buffer.Bytes(), blobURL, azblob.UploadToBlockBlobOptions{})
	if erroor != nil {
		fmt.Println("Can't upload the buffer to block blob", erroor.Error())
	}

	// Check if any errors occurred during writing process to the file
	err = scrape.Writer.Error()
	if err != nil {
		fmt.Println("Can't write to file", err.Error())
	}

	//prints out the final populated struct data to console.
	// fmt.Println(excel.Data)
	// fmt.Println(d)
	a := os.Remove(genExports.GeneralExports)
	if a != nil {
		log.Fatal(a)
	}
	b := os.Remove("Lead_Template.xlsx")
	if b != nil {
		log.Fatal(b)
	}
	time.Sleep(time.Second * 2)
	c := os.Remove(scrape.BlobName)
	if c != nil {
		log.Fatal(c)
	}
}
