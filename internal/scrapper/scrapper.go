package scrapper

import (
	"fmt"
	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const url = "https://prod2.seace.gob.pe/seacebus-uiwd-pub/buscadorPublico/buscadorPublico.xhtml"
const startDateSelector = "tbBuscador:idFormBuscarProceso:dfechaInicio_input"
const endDateSelector = "tbBuscador:idFormBuscarProceso:dfechaFin_input"
const searchButtonSelector = "tbBuscador:idFormBuscarProceso:btnBuscarSelToken"

func Start(date time.Time) {
	fmt.Printf("Process initialized for date: %s\n", date.Format("2006-01-02"))

	// initialize the selenium
	service, err := selenium.NewChromeDriverService("chromedriver", 4444)
	if err != nil {
		fmt.Printf("Failed to start the selenium service: %v", err)
		return
	}

	defer service.Stop()

	driver, err := setupDriver()
	if err != nil {
		fmt.Printf("Failed to open the browser: %v", err)
		return
	}

	err = fillDates(driver, date)
	if err != nil {
		fmt.Printf("Failed to fill dates: %v", err)
		return
	}

	// sleep for 10s
	time.Sleep(10 * time.Second)
	err = takeScreenshot(driver)
	if err != nil {
		fmt.Printf("Failed to take screenshot: %v", err)
	}
}

func setupDriver() (selenium.WebDriver, error) {
	capabilities := selenium.Capabilities{}
	arguments := []string{"--headless"}
	capabilities.AddChrome(chrome.Capabilities{Args: arguments})

	driver, err := selenium.NewRemote(capabilities, "")
	if err != nil {
		return nil, err
	}

	err = driver.Get(url)
	if err != nil {
		return nil, err
	}

	err = driver.ResizeWindow("", 1920, 1080)
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func fillDates(driver selenium.WebDriver, date time.Time) error {
	formattedDate := date.Format("02/01/2006")
	startDateSelector, err := driver.FindElement(selenium.ByID, startDateSelector)
	if err != nil {
		return err
	}
	startDateSelector.SendKeys(formattedDate)

	endDateSelector, err := driver.FindElement(selenium.ByID, endDateSelector)
	if err != nil {
		return err
	}
	endDateSelector.SendKeys(formattedDate)

	button, err := driver.FindElement(selenium.ByID, searchButtonSelector)
	if err != nil {
		return err
	}
	err = button.Click()
	if err != nil {
		return err
	}

	return nil
}

func takeScreenshot(driver selenium.WebDriver) error {
	screenshot, err := driver.Screenshot()
	if err != nil {
		return err
	}

	err = os.WriteFile("/tmp/screenshot.png", screenshot, 0644)
	if err != nil {
		return err
	}

	return nil
}
