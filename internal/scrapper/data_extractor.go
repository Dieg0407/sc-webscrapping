package scrapper

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

const nomenclatureXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[2]/td/table/tbody/tr[1]/td[2]"
const entityXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[6]/td/table/tbody/tr[1]/td[2]"
const objectTypeXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[9]/td/table/tbody/tr[1]/td[2]"
const valueXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[9]/td/table/tbody/tr[3]/td[2]/span[1]"
const currencyXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[9]/td/table/tbody/tr[3]/td[2]/span[2]"
const winnerElement = "tbFicha:idGridLstItems:0:dtParticipantes_data"

func extractData(driver selenium.WebDriver, id int) (ExtractedInformation, error) {
	stdout := log.New(os.Stdout, "", 0)
	stderr := log.New(os.Stderr, "[data-extractor] ", 0)

	nomenclature, err := extractTextByXPath(driver, nomenclatureXPath)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract nomenclature:\n%s", err)
	}
	entity, err := extractTextByXPath(driver, entityXPath)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract entity:\n%s", err)
	}
	objectType, err := extractTextByXPath(driver, objectTypeXPath)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract object type:\n%s", err)
	}
	value, err := extractTextByXPath(driver, valueXPath)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract value:\n%s", err)
	}
	currency, err := extractTextByXPath(driver, currencyXPath)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract currency:\n%s", err)
	}

	description, err := extractDescription(driver)
	if err != nil {
		return ExtractedInformation{}, fmt.Errorf("failed to extract description:\n%s", err)
	}
	hasWinner, winnerData, err := extractWinner(driver)

	if hasWinner {
		winnerName := winnerData[0]
		mype := winnerData[1]
		jungle := winnerData[2]

		stdout.Printf("%d|%s|%s|%s|%s|%s|%s|%s|%s|%s\n",
			id,
			entity,
			nomenclature,
			objectType,
			description,
			value,
			currency,
			winnerName,
			mype,
			jungle,
		)
	} else {
		stderr.Printf("Process with id %d and description %s has no winner\n", id+1, description)
	}

	return ExtractedInformation{}, nil
}

func extractTextByXPath(driver selenium.WebDriver, xpath string) (string, error) {
	element, err := driver.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		return "", err
	}
	text, err := element.Text()
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
	legends, err := driver.FindElements(selenium.ByTagName, "legend")
	if err != nil {
		return "", err
	}

	for _, legend := range legends {
		text, err := legend.Text()
		if err != nil {
			return "", err
		}
		if !strings.Contains(text, "Ver listado") {
			continue
		}

		err = legend.Click()
		if err != nil {
			return "", err
		}

		break
	}

	time.Sleep(1 * time.Second)

	table, err := driver.FindElement(selenium.ByID, "tbFicha:idGridLstItems_content")
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

func extractValue(driver selenium.WebDriver) (string, error) {
	valueTable, err := driver.FindElement(selenium.ByID, "tbFicha:j_idt93")
	if err != nil {
		return "", err
	}
	rows, err := valueTable.FindElements(selenium.ByCSSSelector, ".ui-widget-content")
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

func extractWinner(driver selenium.WebDriver) (bool, []string, error) {
	winnerTable, err := driver.FindElement(selenium.ByID, winnerElement)
	if err != nil {
		return false, nil, err
	}
	columns, err := winnerTable.FindElements(selenium.ByTagName, "td")
	if len(columns) == 1 {
		return false, nil, nil
	}

	winner, err := columns[0].Text()
	if err != nil {
		return false, nil, err
	}

	mype, err := columns[1].Text()
	if err != nil {
		return false, nil, err
	}

	jungle, err := columns[2].Text()
	if err != nil {
		return false, nil, err
	}

	return true, []string{winner, mype, jungle}, nil
}
