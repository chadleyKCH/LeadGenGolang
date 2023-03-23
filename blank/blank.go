package blank

import (
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/scrape"
	"github.com/chadleyKCH/LeadGenGolang/LeadGen_GOLANG/src/search"
)

var (
	Header []string
)

func TXstate() {
	if search.StateAbb == "TX" {
		search.StateAbb = "NT"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
		search.StateAbb = "GT"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "TX - N" {
		search.StateAbb = "NT"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "TX - S" {
		search.StateAbb = "GT"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}

func CAstate() {
	if search.StateAbb == "CA" {
		search.StateAbb = "CN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
		search.StateAbb = "CS"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "CA - N" {
		search.StateAbb = "CN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "CA - S" {
		search.StateAbb = "CS"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}

func MAstate() {
	if search.StateAbb == "MA" {
		search.StateAbb = "EM"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
		search.StateAbb = "WM"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "MA - E" {
		search.StateAbb = "EM"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "MA - W" {
		search.StateAbb = "WM"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}

func NJstate() {
	if search.StateAbb == "NJ" {
		search.StateAbb = "JN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
		search.StateAbb = "JS"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "NJ - N" {
		search.StateAbb = "JN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "NJ - S" {
		search.StateAbb = "JS"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}

func NYstate() {
	if search.StateAbb == "NY" {
		search.StateAbb = "DN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
		search.StateAbb = "UN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "NY - M" {
		search.StateAbb = "DN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "NY - U" {
		search.StateAbb = "UN"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}

func OHstate() {
	if search.StateAbb == "OH" {
		search.StateAbb = "NO"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
		search.StateAbb = "SO"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "OH - N" {
		search.StateAbb = "NO"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "OH - S" {
		search.StateAbb = "SO"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}

func PAstate() {
	if search.StateAbb == "PA" {
		search.StateAbb = "EP"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
		search.StateAbb = "WP"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "PA - E" {
		search.StateAbb = "EP"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
	if search.StateAbb == "PA - W" {
		search.StateAbb = "WP"
		search.SearchThomasnet()
		scrape.ScrapeWebsite()
	}
}
