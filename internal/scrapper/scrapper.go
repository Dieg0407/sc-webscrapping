package scrapper

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const url = "https://prod2.seace.gob.pe/seacebus-uiwd-pub/buscadorPublico/buscadorPublico.xhtml"
const startDateSelector = "tbBuscador:idFormBuscarProceso:dfechaInicio_input"
const endDateSelector = "tbBuscador:idFormBuscarProceso:dfechaFin_input"
const searchButtonSelector = "tbBuscador:idFormBuscarProceso:btnBuscarSelToken"
const retrievedRowsDataContainerSelector = ".ui-paginator-current"
const advancedSearchSelector = ".ui-fieldset-legend"
const tableDataSelector = "tbBuscador:idFormBuscarProceso:dtProcesos_data"
const nextPageButton = ".ui-paginator-next"
const previousPageButton = ".ui-paginator-prev"

func Start(date time.Time) {
	logger := log.New(os.Stderr, "[scrubber] ", log.LstdFlags)
	logger.Printf("Process initialized for date: %s\n", date.Format("2006-01-02"))

	// initialize the selenium
	service, err := selenium.NewChromeDriverService("chromedriver", 4444)
	if err != nil {
		logger.Printf("Failed to start the selenium service:\n%v", err)
		return
	}

	defer service.Stop()

	driver, err := setupDriver()
	if err != nil {
		logger.Printf("Failed to open the browser:\n%v", err)
		return
	}

	err = fillDates(driver, date)
	if err != nil {
		logger.Printf("Failed to fill dates:\n%v", err)
		return
	}

	// sleep for 10s
	time.Sleep(10 * time.Second)

	// Scroll to the bottom of the page using JavaScript
	_, err = driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
	if err != nil {
		logger.Printf("Error scrolling to the bottom: %v", err)
		return
	}

	recordsObtained, err := findTotalAmountOfRows(driver)
	if err != nil {
		logger.Printf("Failed to find total amount of rows:\n%v", err)
		return
	}

	logger.Printf("Total amount of rows obtained: %d\n", recordsObtained)

	if recordsObtained == 0 {
		logger.Printf("No records obtained\n")
		return
	}

	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		_, err := driver.FindElement(selenium.ByID, tableDataSelector)
		if err != nil {
			return false, err
		}
		return true, nil
	}, 10*time.Second)

	if err != nil {
		logger.Printf("Failed to wait for the page to load:\n%v", err)
		return
	}

	rowIdentifierFormat, err := extractRowIdentifierFormat(driver)
	if err != nil {
		logger.Printf("Failed to extract row identifier:\n%v", err)
		return
	}
	logger.Printf("Row identifier format extracted: %s\n", rowIdentifierFormat)

	printHeader()
	for i := 0; i < int(recordsObtained); i++ {
		logger.Println("Processing record: ", i+1)
		err = selectElement(driver, i, rowIdentifierFormat, logger)
		if err != nil {
			logger.Printf("Failed to select element:\n%v", err)
			break
		}

		logger.Println("Record processed successfully")
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
	advancedSearchButton, err := driver.FindElement(selenium.ByCSSSelector, advancedSearchSelector)
	if err != nil {
		return fmt.Errorf("couldn't obtain the advanced search button:\n%w", err)
	}

	err = advancedSearchButton.Click()
	if err != nil {
		return fmt.Errorf("couldn't click the advanced search button:\n%w", err)
	}

	time.Sleep(2 * time.Second)
	formattedDate := date.Format("02/01/2006")
	startDateSelector, err := driver.FindElement(selenium.ByID, startDateSelector)
	if err != nil {
		return fmt.Errorf("couldn't obtain the start date selector:\n%w", err)
	}
	err = startDateSelector.SendKeys(formattedDate)
	if err != nil {
		return fmt.Errorf("couldn't set the start date value:\n%w", err)
	}

	time.Sleep(2 * time.Second)
	endDateSelector, err := driver.FindElement(selenium.ByID, endDateSelector)
	if err != nil {
		return fmt.Errorf("couldn't obtain the end date selector:\n%w", err)
	}
	err = endDateSelector.SendKeys(formattedDate)
	if err != nil {
		return fmt.Errorf("couldn't set the end date value:\n%w", err)
	}

	button, err := driver.FindElement(selenium.ByID, searchButtonSelector)
	if err != nil {
		return fmt.Errorf("couldn't obtain the search button:\n%w", err)
	}
	err = button.Click()
	if err != nil {
		return fmt.Errorf("couldn't click the search button:\n%w", err)
	}

	return nil
}

func findTotalAmountOfRows(driver selenium.WebDriver) (int64, error) {
	retrievedRowsData, err := driver.FindElement(selenium.ByCSSSelector, retrievedRowsDataContainerSelector)
	if err != nil {
		return 0, fmt.Errorf("couldn't obtain the retrieved rows container:\n%s", err)
	}

	text, err := retrievedRowsData.Text()
	if err != nil {
		return 0, fmt.Errorf("couldn't obtain text data from the retrieved rows container:\n%s", err)
	}

	parts := strings.Fields(text)
	if len(parts) < 8 {
		return 0, fmt.Errorf("couldn't extract the total amount of rows from the container:\n%s", text)
	}

	part := parts[8]
	total, err := strconv.ParseInt(part, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse total amount of retrived rows:\n%s", err)
	}

	return total, nil
}

func extractRowIdentifierFormat(driver selenium.WebDriver) (string, error) {
	tableData, err := driver.FindElement(selenium.ByID, tableDataSelector)
	if err != nil {
		return "", fmt.Errorf("couldn't obtain the table data:\n%s", err)
	}

	rows, err := tableData.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		return "", fmt.Errorf("couldn't obtain the rows from the table data:\n%s", err)
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("no rows found in the table data")
	}

	columns, err := rows[0].FindElements(selenium.ByTagName, "td")
	if err != nil {
		return "", fmt.Errorf("couldn't obtain the columns from the row:\n%s", err)
	}

	if len(columns) < 13 {
		return "", fmt.Errorf("not enough columns found in the row, the format may have changed!")
	}

	actions, err := columns[12].FindElements(selenium.ByTagName, "a")
	if err != nil {
		return "", fmt.Errorf("couldn't obtain the actions from the row:\n%s", err)
	}
	if len(actions) < 2 {
		return "", fmt.Errorf("not enough actions found in the row, the format may have changed!")
	}

	goToElementAction := actions[1]
	attribute, err := goToElementAction.GetAttribute("id")
	if err != nil {
		return "", fmt.Errorf("couldn't obtain the id attribute from the action:\n%s", err)
	}

	return strings.Replace(attribute, ":0:", ":%d:", 1), nil
}

func takeScreenshot(driver selenium.WebDriver, path string) error {
	screenshot, err := driver.Screenshot()
	if err != nil {
		return fmt.Errorf("driver failed to take screenshot:\n%s", err)
	}

	err = os.WriteFile(path, screenshot, 0644)
	if err != nil {
		return fmt.Errorf("failed to write screen shot: \n%s", err)
	}

	return nil
}

func selectElement(driver selenium.WebDriver, id int, rowIdentifierFormat string, logger *log.Logger) error {
	err := goToPage(driver, calculatePageNumber(id), logger)
	if err != nil {
		return fmt.Errorf("failed to go to the page %d:\n%s", calculatePageNumber(id), err)
	}

	formattedId := fmt.Sprintf(rowIdentifierFormat, id)
	element, err := driver.FindElement(selenium.ByID, formattedId)
	if err != nil {
		return fmt.Errorf("failed to obtain the element with id %d and raw id '%s':\n%s", id, formattedId, err)
	}

	err = element.Click()
	if err != nil {
		return fmt.Errorf("failed to click the element with id %d and raw id '%s':\n%s", id, formattedId, err)
	}

	err = driver.WaitWithTimeout(waitForDetailsPageToLoad, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for the details page to load:\n%s", err)
	}

	// extract information
	err = extractData(driver, id)
	if err != nil {
		return fmt.Errorf("failed to extract data:\n%s", err)
	}

	// go back
	element, err = driver.FindElement(selenium.ByXPATH, "//button[span[text()='Regresar']]")
	if err != nil {
		return fmt.Errorf("failed to obtain the go back button:\n%s", err)
	}
	err = element.Click()
	if err != nil {
		return fmt.Errorf("failed to click the go back button:\n%s", err)
	}

	err = driver.WaitWithTimeout(waitForMainPageToLoad, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for the main page to load:\n%s", err)
	}

	return nil
}

func waitForDetailsPageToLoad(wd selenium.WebDriver) (bool, error) {
	_, err := wd.FindElement(selenium.ByXPATH, "//button[span[text()='Regresar']]")
	if err != nil {
		return false, fmt.Errorf("failed to obtain the go to the details page:\n%s", err)
	}

	return true, nil
}

func waitForMainPageToLoad(wd selenium.WebDriver) (bool, error) {
	_, err := wd.FindElement(selenium.ByCSSSelector, advancedSearchSelector)
	if err != nil {
		return false, fmt.Errorf("failed to obtain the go back to the main page:\n%s", err)
	}

	return true, nil
}

func calculatePageNumber(id int) int {
	page := id / 15
	page++

	return page
}

func goToPage(driver selenium.WebDriver, page int, logger *log.Logger) error {
	paginator, err := driver.FindElements(selenium.ByCSSSelector, ".ui-paginator-page")
	if err != nil {
		return fmt.Errorf("failed to obtain the paginator elements:\n%s", err)
	}

	activePage := 0
	for _, element := range paginator {
		classNames, err := element.GetAttribute("class")
		if err != nil {
			return fmt.Errorf("failed to obtain the class names:\n%s", err)
		}

		if !strings.Contains(classNames, "ui-state-active") {
			continue
		}

		text, err := element.Text()
		if err != nil {
			return fmt.Errorf("failed to obtain the text from the element:\n%s", err)
		}

		activePage, err = strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("failed to parse the active page number:\n%s", err)
		}
		break
	}

	if activePage == page {
		return nil
	}
	if activePage < page {
		// move forward
		logger.Println("moving forward")
		err = clickNextPage(driver)
		if err != nil {
			return fmt.Errorf("failed to click the next page button:\n%s", err)
		}
		_, err = driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
		if err != nil {
			return err
		}
		return goToPage(driver, page, logger)
	}

	logger.Println("moving backward")
	err = clickPreviousPage(driver)
	if err != nil {
		return fmt.Errorf("failed to click the previous page button:\n%s", err)
	}
	return goToPage(driver, page, logger)
}

func clickNextPage(driver selenium.WebDriver) error {
	nextPage, err := driver.FindElement(selenium.ByCSSSelector, nextPageButton)
	if err != nil {
		return fmt.Errorf("failed to obtain the next page button:\n%s", err)
	}

	classNames, err := nextPage.GetAttribute("class")
	if err != nil {
		return fmt.Errorf("failed to obtain the class names:\n%s", err)
	}

	if strings.Contains(classNames, "ui-state-disabled") {
		return fmt.Errorf("the next page button is disabled")
	}

	err = nextPage.Click()
	if err != nil {
		return fmt.Errorf("failed to click the next page button:\n%s", err)
	}
	time.Sleep(5 * time.Second)
	return nil
}

func clickPreviousPage(driver selenium.WebDriver) error {
	previousPage, err := driver.FindElement(selenium.ByCSSSelector, previousPageButton)
	if err != nil {
		return fmt.Errorf("failed to obtain the previous page button:\n%s", err)
	}

	classNames, err := previousPage.GetAttribute("class")
	if err != nil {
		return fmt.Errorf("failed to obtain the class names:\n%s", err)
	}

	if strings.Contains(classNames, "ui-state-disabled") {
		return fmt.Errorf("the previous page button is disabled")
	}

	err = previousPage.Click()
	if err != nil {
		return fmt.Errorf("failed to click the previous page button:\n%s", err)
	}
	time.Sleep(5 * time.Second)
	return nil
}
