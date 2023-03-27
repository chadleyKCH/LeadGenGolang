package genExports

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"

	"lead-generator/scrape"
	"lead-generator/search"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/tealeg/xlsx"
)

var (
	accountName    = "leadgeneratorstorageblob"
	accountKey     = "rOjnqusNKyYwLdVeK+v0eKnxN0lyRi5vbPQ0csRDFNWOnpwiutTU0Zg/m2pMn1AyCKcdJIyiJbut+AStkcj6Rg=="
	containerName  = "lead-generator"
	GeneralExports = "GeneralExports.xlsx"
)

func GenExports() {
	//gets shared key credentials which is used to sign in to Azure services instead of the account access keys.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	//builds ablob storage pipeline object for getting data from Azure blob storage.
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	ctx := context.Background()

	//creates a new URL structure for use in the Blob service URLs (for the specified storage account).
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	// creates a service URL object pointing to the service's specified blob location.
	serviceURL := azblob.NewServiceURL(*u, p)
	//builds the container's URL using the ServiceURL object.
	containerURL := serviceURL.NewContainerURL(containerName)
	// gets the specified blob's URL.
	blobURL := containerURL.NewBlobURL(GeneralExports)

	//downloads the specified block blob using one operation and returns its content as an io.ReadCloser.
	responseBody, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}

	//gets a stream to read data from and retries based on specified criteria.
	bodyStream := responseBody.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	//creating an empty buffer value to append downloaded data.
	downloadedData := bytes.Buffer{}

	// reads downloaded data into a buffer.
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		fmt.Println(err.Error())
	}
	//opens or create a new file to copy downloaded data from body to a local machine.
	file, err := os.OpenFile("GeneralExports.xlsx", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}

	//writes the downloaded data byte.buffer to a local file, creating it if it doesn't exist.
	_, err = file.Write(downloadedData.Bytes())
	if err != nil {
		fmt.Println(err.Error())
	}
	// NEED TO ITERATE THROUGH THE GENERAL EXPORTS AND SET LEAD TO THAT THEN SEARCH

	genExports, err := xlsx.OpenFile(GeneralExports)
	if err != nil {
		fmt.Println(err.Error())
	}

	sheet := genExports.Sheets[0]

	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if row.Cells[0].Value == "" {
			break
		}

		search.Lead = row.Cells[0].Value
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}
