package scrape

import (
	"encoding/csv"
	"fmt"
	"os"

	"lead-generator/search"

	"github.com/gocolly/colly"
)

type company_chrome struct {
	companyName, companyDesc, companyLoc, companyType string
}

var (
	companys []company_chrome
	Writer   *csv.Writer
	File     *os.File
	BlobName = "Leads.csv"
)

// COMPLETED
// Scrapes the specified lead on Thomasnet and then returns the results

func ScrapeWebsite() {
	// Use colly.NewCollector to create a new Collector instance
	c := colly.NewCollector()

	// OnRequest is a function handler that is run when a new request is made
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	// Initialize the CSV writer
	Writer = csv.NewWriter(File)

	// Flush any buffered data in the CSV writer
	Writer.Flush()

	// OnHTML is a function handler that is called whenever element(s) matching a selector are found
	c.OnHTML("div.profile-card", func(e *colly.HTMLElement) {
		// Create a variable called company of type company_card
		company := company_chrome{}
		company.companyName = e.ChildText("h2.profile-card__title")
		company.companyLoc = e.ChildText("span.profile-card__location a")
		company.companyType = e.ChildText("span[data-content=\"Company Type\"]")
		company.companyDesc = e.ChildText("div.profile-card__body-text p")

		// Append the newly created company variable to the existing companys slice
		companys = append(companys, company)
		// Write the values stored inside the most recently created company variable to the CSV file
		if len(companys) != 0 {
			leads := []string{
				company.companyName,
				company.companyLoc,
				company.companyType,
				company.companyDesc,
			}
			Writer.Write(leads)
			Writer.Flush()
		}
	})

	// Visit a specified website (in this case search.DriverURL) via the collector
	c.Visit(search.DriverURL)
}
