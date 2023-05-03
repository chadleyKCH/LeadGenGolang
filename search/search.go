package search

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

var (
	StateAbb  string
	DriverURL string
	Lead      string
)

func SearchThomasnet() {

	const (
		// These paths will be different on your system.
		// chromeDriverPath = "C:/Users/coleh/VS_CODES/LG_Other/chromedriver.exe"
		chromeDriverPath = "/usr/local/bin/chromedriver"
		port             = 4445
	)

	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
		selenium.Output(os.Stderr),
	}
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"goog:chromeOptions": map[string]interface{}{
			"args": []string{
				"--headless",
				"--disable-gpu",
				"--no-sandbox",
				"--disable-dev-shm-usage",
			},
		},
	}
	P := 4445
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", P))
	if err != nil {
		log.Fatal("COULDN'T START NEW SELENIUM REMOTE")
	}
	driver.Get("http://www.thomasnet.com/")
	time.Sleep(time.Second * 3)
	// Select the Search bar
	input, err := driver.FindElement(selenium.ByCSSSelector, "#homesearch > form > div > div > div.site-search__search-query-input-wrap.search-suggest-preview > input")
	if err != nil {
		log.Fatalf("FAILED TO FIND THE SEARCH BAR: %v", err)
	}
	//Enter lead from Lead_Template.xlsx into the search bar
	err = input.SendKeys(Lead)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 2)
	// Define the button selector and search for its element.
	search, err := driver.FindElement(selenium.ByCSSSelector, "#homesearch > form > div > button")
	if err != nil {
		panic(err)
	}
	//Click on the search button.
	err = search.Click()
	if err != nil {
		log.Fatalf("FAILED TO CLICK THE SEARCH BUTTON: %v", err)
	}
	time.Sleep(time.Second * 2)
	// Looks for iframe
	// iframeElem, err := driver.FindElement(selenium.ByCSSSelector, "iframe[src*='about:blank']")
	// if err != nil {
	// 	log.Fatalf("FAILED TO FIND THE IFRAME ELEMENT: %v", err)
	// }
	// // Switches to iframe
	// err = driver.SwitchFrame(iframeElem)
	// if err != nil {
	// 	log.Fatalf("FAILED TO SWITCH TO THE IFRAME: %v", err)
	// }
	// time.Sleep(time.Second * 2)
	// // Looks for accept button
	// acceptNoti, err := driver.FindElement(selenium.ByXPATH, "/html/body/appcues/cue/section/div/div[3]/div/div/div/div/div/a")
	// if err != nil {
	// 	log.Fatalf("FAILED TO FIND THE IFRAME ACCEPT BUTTON: %v", err)
	// }
	// // Clicks accept button
	// err = acceptNoti.Click()
	// if err != nil {
	// 	log.Fatalf("FAILED TO CLICK THE IFRAME ACCEPT BUTTON: %v", err)
	// }
	// // Switches back to main frame
	// switchDefault := driver.SwitchFrame(nil)
	// if switchDefault != nil {
	// 	log.Fatalf("FAILED TO SWITCH BACK TO MAIN FRAME: %v", switchDefault)
	// }
	// time.Sleep(time.Second * 2)

	// Checks if StateAbb is blank. If so, then return the driver URL to scrape.go
	if StateAbb == "" {
		DriverURL, _ = driver.CurrentURL()
		return
	} else {
		// Find the region dropdown
		regionDropdown, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > div > div.thm-custom-select.search-options-regions > a")
		if err != nil {
			log.Fatalf("FAILED TO FIND SELECT REGION DROPDOWN: %v", err)
		}
		// Click the region dropdown
		err = regionDropdown.Click()
		if err != nil {
			log.Fatalf("FAILED TO CLICK SELECT REGION DROPDOWN: %v", err)
		}
		time.Sleep(time.Second * 2)

		// Finds specified region
		regionSelect, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > div > div.thm-custom-select.search-options-regions > div [data-value="+StateAbb+"]")
		if err != nil {
			log.Fatalf("COULD NOT FIND THE stateabb SPECIFIED: %v", err)
		}
		// Clicks specified region
		err = regionSelect.Click()
		if err != nil {
			log.Fatalf("COULD NOT CLICK THE stateabb: %v", err)
		}
		// Select the region dropdown
		err = regionDropdown.Click()
		if err != nil {
			log.Fatalf("COULD NOT CLICK THE REGION DROPDOWN: %v", err)
		}
		time.Sleep(time.Second * 2)
		// Finds search button
		regionSearch, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > button")
		if err != nil {
			log.Fatalf("COULD NOT FIND SEARCH BUTTON AFTER THE REGION HAS BEEN SELECTED: %v", err)
		}
		// Clicks search
		err = regionSearch.Click()
		if err != nil {
			log.Fatalf("COULD NOT CLICK SEARCH BUTTON AFTER THE REGION HAS BEEN SELECTED: %v", err)
		}
		time.Sleep(time.Second * 2)

		results, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.interim-search-results.logged-out > section.network-search-results > div > div.network-search-results__primary > div > section > div")
		if err != nil {
			fmt.Printf("COULD NOT FIND NETWORK RESULT: %v", err)
		}
		if results != nil {
			err = results.Click()
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("ELEMENT NETWORK NOT FOUND", err)
		}
		// Finds the 'Located in' option
		LocIn, err := driver.FindElement(selenium.ByCSSSelector, "#main > div.filter-block.located-serving-card > ul > li:nth-child(1) > a")
		if err != nil {
			log.Fatalf("ERROR WITH 'LOCATED IN' PART OF CODE: %v", err)
		}
		time.Sleep(time.Second * 2)
		// If the 'Located in' option does not exist then it skips over...else: clicks located in
		if LocIn != nil {
			// Clicks 'Located in'
			err = LocIn.Click()
			if err != nil {
				log.Fatalf("COULD NOT CLICK ON 'LOCATED IN': %v", err)
			}
		} else {
			log.Fatalf("COULD NOT FIND 'LOCATED IN' ELEMENT: %v", err)
		}
		time.Sleep(time.Second * 2)

		DriverURL, _ = driver.CurrentURL()
	}

}
