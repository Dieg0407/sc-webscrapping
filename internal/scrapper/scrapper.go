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
const selectionProceduresOption = "/html/body/div[3]/div/div[1]/ul/li[2]/a"

func Start(date time.Time) {
	logger := log.New(os.Stderr, "[scrapper] ", log.LstdFlags)
	logger.Printf("Proceso inicializado para la fecha: %s\n", date.Format("2006-01-02"))

	// Inicializar Selenium
	service, err := selenium.NewChromeDriverService("chromedriver", 4444)
	if err != nil {
		logger.Printf("Error al iniciar el servicio de Selenium:\n%v", err)
		return
	}

	defer service.Stop()

	driver, err := setupDriver()
	if err != nil {
		logger.Printf("Error al abrir el navegador:\n%v", err)
		return
	}
	err = selectSelectionProceduresOption(driver)
	if err != nil {
		logger.Printf("Error al seleccionar la opción de procedimientos de selección:\n%v", err)
		return
	}

	err = fillDates(driver, date)
	if err != nil {
		logger.Printf("Error al rellenar las fechas:\n%v", err)
		return
	}

	// Esperar 10 segundos
	time.Sleep(10 * time.Second)

	// Desplazarse hasta el final de la página usando JavaScript
	_, err = driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
	if err != nil {
		logger.Printf("Error al desplazarse al final: %v", err)
		return
	}

	recordsObtained, err := findTotalAmountOfRows(driver)
	if err != nil {
		logger.Printf("Error al encontrar la cantidad total de filas:\n%v", err)
		return
	}

	logger.Printf("Cantidad total de filas obtenidas: %d\n", recordsObtained)

	if recordsObtained == 0 {
		logger.Printf("No se obtuvieron registros\n")
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
		logger.Printf("Error al esperar a que la página se cargue:\n%v", err)
		return
	}

	rowIdentifierFormat, err := extractRowIdentifierFormat(driver)
	if err != nil {
		logger.Printf("Error al extraer el identificador de fila:\n%v", err)
		return
	}
	logger.Printf("Formato de identificador de fila extraído: %s\n", rowIdentifierFormat)

	printHeader()
	for i := 0; i < int(recordsObtained); i++ {
		logger.Println("Procesando registro: ", i+1)
		err = selectElement(driver, i, rowIdentifierFormat, logger)
		if err != nil {
			logger.Printf("Error al seleccionar el elemento:\n%v", err)
			break
		}

		logger.Println("Registro procesado correctamente")
	}

	logger.Println("Proceso finalizado")
}

func selectSelectionProceduresOption(driver selenium.WebDriver) error {
	option, err := driver.FindElement(selenium.ByXPATH, selectionProceduresOption)
	if err != nil {
		return fmt.Errorf("no se pudo encontrar la opción de procedimientos de selección:\n%w", err)
	}
	err = option.Click()
	if err != nil {
		return fmt.Errorf("no se pudo hacer clic en la opción de procedimientos de selección:\n%w", err)
	}
	time.Sleep(2 * time.Second)
	return nil
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
		return fmt.Errorf("no se pudo obtener el botón de búsqueda avanzada:\n%w", err)
	}

	err = advancedSearchButton.Click()
	if err != nil {
		return fmt.Errorf("no se pudo hacer clic en el botón de búsqueda avanzada:\n%w", err)
	}

	time.Sleep(2 * time.Second)
	formattedDate := date.Format("02/01/2006")
	startDateSelector, err := driver.FindElement(selenium.ByID, startDateSelector)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el selector de fecha de inicio:\n%w", err)
	}
	err = startDateSelector.SendKeys(formattedDate)
	if err != nil {
		return fmt.Errorf("no se pudo establecer el valor de la fecha de inicio:\n%w", err)
	}

	time.Sleep(2 * time.Second)
	endDateSelector, err := driver.FindElement(selenium.ByID, endDateSelector)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el selector de fecha de fin:\n%w", err)
	}
	err = endDateSelector.SendKeys(formattedDate)
	if err != nil {
		return fmt.Errorf("no se pudo establecer el valor de la fecha de fin:\n%w", err)
	}

	button, err := driver.FindElement(selenium.ByID, searchButtonSelector)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el botón de búsqueda:\n%w", err)
	}
	err = button.Click()
	if err != nil {
		return fmt.Errorf("no se pudo hacer clic en el botón de búsqueda:\n%w", err)
	}

	return nil
}

func findTotalAmountOfRows(driver selenium.WebDriver) (int64, error) {
	retrievedRowsData, err := driver.FindElements(selenium.ByCSSSelector, retrievedRowsDataContainerSelector)
	if err != nil {
		return 0, fmt.Errorf("no se pudo obtener el contenedor de filas recuperadas:\n%s", err)
	}
	var parts []string = []string{}
	var text string
	for _, element := range retrievedRowsData {
		text, err := element.Text()
		if err != nil {
			return 0, fmt.Errorf("no se pudo obtener el texto del contenedor de filas recuperadas:\n%s", err)
		}
		if strings.Contains(text, "Mostrando") {
			parts = strings.Fields(text)
			break
		}
	}

	if len(parts) < 8 {
		return 0, fmt.Errorf("no se pudo extraer la cantidad total de filas del contenedor:\n%s", text)
	}

	part := parts[8]
	total, err := strconv.ParseInt(part, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error al analizar la cantidad total de filas recuperadas:\n%s", err)
	}

	return total, nil
}

func extractRowIdentifierFormat(driver selenium.WebDriver) (string, error) {
	tableData, err := driver.FindElement(selenium.ByID, tableDataSelector)
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener los datos de la tabla:\n%s", err)
	}

	rows, err := tableData.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener las filas de los datos de la tabla:\n%s", err)
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("no se encontraron filas en los datos de la tabla")
	}

	columns, err := rows[0].FindElements(selenium.ByTagName, "td")
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener las columnas de la fila:\n%s", err)
	}

	if len(columns) < 13 {
		return "", fmt.Errorf("no se encontraron suficientes columnas en la fila, ¡el formato puede haber cambiado!")
	}

	actions, err := columns[12].FindElements(selenium.ByTagName, "a")
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener las acciones de la fila:\n%s", err)
	}
	if len(actions) < 2 {
		return "", fmt.Errorf("no se encontraron suficientes acciones en la fila, ¡el formato puede haber cambiado!")
	}

	goToElementAction := actions[1]
	attribute, err := goToElementAction.GetAttribute("id")
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el atributo id de la acción:\n%s", err)
	}

	return strings.Replace(attribute, ":0:", ":%d:", 1), nil
}

func selectElement(driver selenium.WebDriver, id int, rowIdentifierFormat string, logger *log.Logger) error {
	err := goToPage(driver, calculatePageNumber(id), logger)
	if err != nil {
		return fmt.Errorf("error al ir a la página %d:\n%s", calculatePageNumber(id), err)
	}

	formattedId := fmt.Sprintf(rowIdentifierFormat, id)
	element, err := driver.FindElement(selenium.ByID, formattedId)
	if err != nil {
		return fmt.Errorf("error al obtener el elemento con id %d e id sin formato '%s':\n%s", id, formattedId, err)
	}

	err = element.Click()
	if err != nil {
		return fmt.Errorf("error al hacer clic en el elemento con id %d e id sin formato '%s':\n%s", id, formattedId, err)
	}

	err = driver.WaitWithTimeout(waitForDetailsPageToLoad, 30*time.Second)
	if err != nil {
		return fmt.Errorf("error al esperar a que se cargue la página de detalles:\n%s", err)
	}

	// Extraer información
	err = extractData(driver, id)
	if err != nil {
		return fmt.Errorf("error al extraer datos:\n%s", err)
	}

	// Regresar
	element, err = driver.FindElement(selenium.ByXPATH, "//button[span[text()='Regresar']]")
	if err != nil {
		return fmt.Errorf("error al obtener el botón Regresar:\n%s", err)
	}
	err = element.Click()
	if err != nil {
		return fmt.Errorf("error al hacer clic en el botón Regresar:\n%s", err)
	}

	err = driver.WaitWithTimeout(waitForMainPageToLoad, 30*time.Second)
	if err != nil {
		return fmt.Errorf("error al esperar a que se cargue la página principal:\n%s", err)
	}

	return nil
}

func waitForDetailsPageToLoad(wd selenium.WebDriver) (bool, error) {
	_, err := wd.FindElement(selenium.ByXPATH, "//button[span[text()='Regresar']]")
	if err != nil {
		return false, fmt.Errorf("error al obtener la página de detalles:\n%s", err)
	}

	return true, nil
}

func waitForMainPageToLoad(wd selenium.WebDriver) (bool, error) {
	_, err := wd.FindElement(selenium.ByCSSSelector, advancedSearchSelector)
	if err != nil {
		return false, fmt.Errorf("error al obtener la página principal:\n%s", err)
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
		return fmt.Errorf("error al obtener los elementos del paginador:\n%s", err)
	}

	activePage := 0
	for _, element := range paginator {
		classNames, err := element.GetAttribute("class")
		if err != nil {
			return fmt.Errorf("error al obtener los nombres de clase:\n%s", err)
		}

		if !strings.Contains(classNames, "ui-state-active") {
			continue
		}

		text, err := element.Text()
		if err != nil {
			return fmt.Errorf("error al obtener el texto del elemento:\n%s", err)
		}

		activePage, err = strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("error al analizar el número de página activo:\n%s", err)
		}
		break
	}

	if activePage == page {
		return nil
	}
	if activePage < page {
		// Avanzar
		logger.Println("avanzando")
		err = clickNextPage(driver)
		if err != nil {
			return fmt.Errorf("error al hacer clic en el botón de siguiente página:\n%s", err)
		}
		_, err = driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
		if err != nil {
			return err
		}
		return goToPage(driver, page, logger)
	}

	logger.Println("retrocediendo")
	err = clickPreviousPage(driver)
	if err != nil {
		return fmt.Errorf("error al hacer clic en el botón de página anterior:\n%s", err)
	}
	return goToPage(driver, page, logger)
}

func clickNextPage(driver selenium.WebDriver) error {
	nextPageButton, err := findActiveButton(driver, nextPageButton)
	if err != nil {
		return fmt.Errorf("error al encontrar el botón activo de página siguiente:\n%s", err)
	}
	err = nextPageButton.Click()
	if err != nil {
		return fmt.Errorf("error al hacer clic en el botón de página siguiente:\n%s", err)
	}

	time.Sleep(5 * time.Second)
	return nil
}

func clickPreviousPage(driver selenium.WebDriver) error {
	previousPageButton, err := findActiveButton(driver, previousPageButton)
	if err != nil {
		return fmt.Errorf("error al encontrar el botón activo de página anterior:\n%s", err)
	}

	err = previousPageButton.Click()
	if err != nil {
		return fmt.Errorf("error al hacer clic en el botón de página anterior:\n%s", err)
	}
	time.Sleep(5 * time.Second)
	return nil
}

func findActiveButton(driver selenium.WebDriver, buttonSelector string) (selenium.WebElement, error) {
	buttons, err := driver.FindElements(selenium.ByCSSSelector, buttonSelector)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los botones:\n%s", err)
	}

	for _, button := range buttons {
		classNames, err := button.GetAttribute("class")
		if err != nil && err.Error() != "stale element not found" {
			return nil, fmt.Errorf("error al obtener los nombres de clase:\n%s", err)
		}
		if !strings.Contains(classNames, "ui-state-disabled") {
			return button, nil
		}
	}

	return nil, fmt.Errorf("no se encontró el botón activo")
}
