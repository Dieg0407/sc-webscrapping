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
	entity, err := extractEntity(driver)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract entity:\n%s", err)
	}
	objectType, err := extractObjectType(driver)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract object type:\n%s", err)
	}

	description, err := extractDescription(driver)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract description:\n%s", err)
	}

	logger.Printf("%d|%s|%s|%s|%s\n", id, entity, nomenclature, objectType, description)
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

func extractEntity(driver selenium.WebDriver) (string, error) {
	entityInformationTable, err := driver.FindElement(selenium.ByID, "tbFicha:j_idt68")
	if err != nil {
		return "", err
	}

	rows, err := entityInformationTable.FindElements(selenium.ByCSSSelector, ".ui-widget-content")
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

func extractObjectType(driver selenium.WebDriver) (string, error) {
	objectTypeTable, err := driver.FindElement(selenium.ByID, "tbFicha:j_idt92")
	if err != nil {
		return "", err
	}
	rows, err := objectTypeTable.FindElements(selenium.ByCSSSelector, ".ui-widget-content")
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

func extractDescription(driver selenium.WebDriver) (string, error) {
	table, err := driver.FindElement(selenium.ByID, "bb")
	if err != nil {
		return "", err
	}

	information, err := table.FindElements(selenium.ByTagName, "span")
	if err != nil {
		return "", err
	}

	if len(information) < 1 {
		return "", fmt.Errorf("no description found")
	}

	span := information[0]
	description, err := span.Text()

	if err != nil {
		return "", err
	}

	return description, nil
}
