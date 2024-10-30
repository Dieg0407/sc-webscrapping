package scrapper

import (
	"fmt"
	"log"
	"os"

	"github.com/tebeka/selenium"
)

type ExtractedInformation struct {
	id           int
	entity       string
	nomenclature string
	objectType   string
	description  string
	value        string
	currency     string
	winner       string
	mype         string
	jungle       string
}

func extractData(driver selenium.WebDriver, id int) (ExtractedInformation, error) {
	logger := log.New(os.Stdout, "", 0)

	nomenclature, err := extractNomenclature(driver)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract nomenclature:\n%s", err)
	}

	logger.Printf("%d|%s\n", id, nomenclature)
	return ExtractedInformation{}, nil
}

func extractNomenclature(driver selenium.WebDriver) (string, error) {
	generalInformationTable, err := driver.FindElement(selenium.ByID, "tbFicha:j_idt25")
	if err != nil {
		return "", err
	}

	rows, err := generalInformationTable.FindElements(selenium.ByCSSSelector, ".ui-widget-content")
	if err != nil {
		return "", err
	}

	row := rows[0]
	headerAndValue, err := row.FindElements(selenium.ByTagName, "td")
	if err != nil {
		return "", err
	}

	value := headerAndValue[1]
	text, err := value.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}
