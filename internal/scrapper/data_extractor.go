package scrapper

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

const nomenclatureXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[2]/td/table/tbody/tr[1]/td[2]"
const entityXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[6]/td/table/tbody/tr[1]/td[2]"
const objectTypeXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[9]/td/table/tbody/tr[1]/td[2]"
const valueXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[9]/td/table/tbody/tr[3]/td[2]/span[1]"
const currencyXPath = "/html/body/div[3]/div/div/div/div/form/table[2]/tbody/tr[1]/td[1]/table/tbody/tr/td/fieldset/div/table/tbody/tr[9]/td/table/tbody/tr[3]/td[2]/span[2]"
const winnerElement = "tbFicha:idGridLstItems:0:dtParticipantes_data"
const printTemplate = "%s;\"%s\";\"%s\";%s;\"%s\";%s;%s;\"%s\";%s;%s\n"

func printHeader() {
	stdout := log.New(os.Stdout, "", 0)
	stdout.Printf(printTemplate,
		"Identificador",
		"Entidad",
		"Nomenclarura",
		"Objecto",
		"Descripción",
		"Valor",
		"Moneda",
		"Ganador",
		"Es MYPE",
		"Es Selva",
	)
}

func extractData(driver selenium.WebDriver, id int) error {
	stdout := log.New(os.Stdout, "", 0)
	stderr := log.New(os.Stderr, "[extractor-datos] ", 0)

	nomenclature, err := extractTextByXPath(driver, nomenclatureXPath)
	if err != nil {
		return fmt.Errorf("error al extraer la nomenclatura:\n%s", err)
	}
	entity, err := extractTextByXPath(driver, entityXPath)
	if err != nil {
		return fmt.Errorf("error al extraer la entidad:\n%s", err)
	}
	objectType, err := extractTextByXPath(driver, objectTypeXPath)
	if err != nil {
		return fmt.Errorf("error al extraer el tipo de objeto:\n%s", err)
	}
	value, err := extractTextByXPath(driver, valueXPath)
	if err != nil {
		return fmt.Errorf("error al extraer el valor:\n%s", err)
	}
	currency, err := extractTextByXPath(driver, currencyXPath)
	if err != nil {
		return fmt.Errorf("error al extraer la moneda:\n%s", err)
	}

	description, err := extractDescription(driver)
	if err != nil {
		return fmt.Errorf("error al extraer la descripción:\n%s", err)
	}
	hasWinner, winnerData, err := extractWinner(driver)

	if hasWinner {
		winnerName := winnerData[0]
		mype := winnerData[1]
		jungle := winnerData[2]

		stderr.Printf("El proceso con id %d y descripción %s tiene un ganador\n", id+1, description)
		stdout.Printf(printTemplate,
			fmt.Sprintf("%d", id+1),
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
		stderr.Printf("El proceso con id %d y descripción %s no tiene ganador\n", id+1, description)
	}

	return nil
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
