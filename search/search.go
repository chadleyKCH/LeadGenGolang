package search

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var (
	StateAbb  string
	DriverURL string
	Lead      string
)

func SearchThomasnet() {
	pathToChrome := "C:/Users/coleh/LeadGen GOLANG_2/chromedriver.exe" //The path to chromedriver. //set to chrome blob

	service, err := selenium.NewChromeDriverService(pathToChrome, 4444) //Initialize driver service with chromedriver.
	if err != nil {
		panic(err)
	}
	defer service.Stop()
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		// "--headless",
	}})
	driver, err := selenium.NewRemote(caps, "") //Create a new session of Google Chrome.
	if err != nil {
		panic(err)
	}
	driver.Get("https://www.thomasnet.com/")
	time.Sleep(time.Second * 4)
	// Select the Search bar
	input, err := driver.FindElement(selenium.ByCSSSelector, "#homesearch > form > div > div > div.site-search__search-query-input-wrap.search-suggest-preview > input")
	if err != nil {
		panic(err)
	}
	//Enter lead from Lead_Template.xlsx into the search bar
	//CHECK IF LEAD IS NULL IF SO THEN PULL GEN EXPORTS.
	err = input.SendKeys(Lead)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 1)
	// Define the button selector and search for its element.
	search, err := driver.FindElement(selenium.ByCSSSelector, "#homesearch > form > div > button")
	if err != nil {
		panic(err)
	}
	//Click on the search button.
	err = search.Click()
	if err != nil {
		log.Fatalf("Error with search: %v", err)
	}
	time.Sleep(time.Second * 2)
	// Looks for iframe
	iframeElem, err := driver.FindElement(selenium.ByCSSSelector, "iframe[src*='about:blank']")
	if err != nil {
		log.Fatalf("Failed to find iframe element: %v", err)
	}
	// Switches to iframe
	err = driver.SwitchFrame(iframeElem)
	if err != nil {
		log.Fatalf("Failed to switch to iframe: %v", err)
	}
	time.Sleep(time.Second * 1)
	// Looks for accept button
	acceptNoti, err := driver.FindElement(selenium.ByXPATH, "/html/body/appcues/cue/section/div/div[3]/div/div/div/div/div/a")
	if err != nil {
		log.Fatalf("Error with finding the iframe accept button: %v", err)
	}
	// Clicks accept button
	err = acceptNoti.Click()
	if err != nil {
		log.Fatalf("Error clicking the iframe alert: %v", err)
	}
	// Switches back to main frame
	switchDefault := driver.SwitchFrame(nil)
	if switchDefault != nil {
		log.Fatalf("Error with switching the frame back: %v", switchDefault)
	}
	time.Sleep(time.Second * 1)

	// Checks if StateAbb is blank. If so, then return the driver URL to scrape.go
	if StateAbb == "" {
		DriverURL, _ = driver.CurrentURL()
		return
	} else {
		// Find the region dropdown
		regionDropdown, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > div > div.thm-custom-select.search-options-regions > a")
		if err != nil {
			log.Fatalf("Failed to find select region dropdown: %v", err)
		}
		// Click the region dropdown
		err = regionDropdown.Click()
		if err != nil {
			log.Fatalf("Could not click select region dropdown: %v", err)
		}
		time.Sleep(time.Second * 1)
		// Finds specified region
		regionSelect, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > div > div.thm-custom-select.search-options-regions > div [data-value="+StateAbb+"]")
		if err != nil {
			log.Fatalf("Could not select the region %v", err)
		}
		// Clicks specified region
		err = regionSelect.Click()
		if err != nil {
			log.Fatalf("Could not click the selected region %v", err)
		}
		// Select the region dropdown
		err = regionDropdown.Click()
		if err != nil {
			log.Fatalf("Could not click select region dropdown: %v", err)
		}
		time.Sleep(time.Second * 1)
		// Finds search button
		regionSearch, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > button")
		if err != nil {
			log.Fatalf("Could not select the region: %v", err)
		}
		// Clicks search
		err = regionSearch.Click()
		if err != nil {
			log.Fatalf("Could not click the selected region %v", err)
		}
		time.Sleep(time.Second * 1)
		results, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.interim-search-results.logged-out > section.network-search-results > div > div.network-search-results__primary > div > section > div")
		if err != nil {
			fmt.Println("Error finding Network Result ", err)
		}
		if results != nil {
			err = results.Click()
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Element Network Result not found", err)
		}
		// Finds the 'Located in' option
		LocIn, err := driver.FindElement(selenium.ByCSSSelector, "#main > div.filter-block.located-serving-card > ul > li:nth-child(1) > a")
		if err != nil {
			log.Printf("Error with 'Located In' portion of the code: %v", err)
		}
		// If the 'Located in' option does not exist then it skips over...else: clicks located in
		if LocIn != nil {
			// Clicks 'Located in'
			err = LocIn.Click()
			if err != nil {
				log.Printf("Error clicking on 'Located In': %v", err)
			}
		} else {
			log.Printf("Element 'Located In' not found.")
		}
		time.Sleep(time.Second * 2)

		DriverURL, _ = driver.CurrentURL()
	}

}
