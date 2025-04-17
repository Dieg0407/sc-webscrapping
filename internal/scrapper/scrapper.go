package scrapper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	// URLs y selectores
	url                                = "https://prod2.seace.gob.pe/seacebus-uiwd-pub/buscadorPublico/buscadorPublico.xhtml"
	startDateSelector                  = "tbBuscador:idFormBuscarProceso:dfechaInicio_input"
	endDateSelector                    = "tbBuscador:idFormBuscarProceso:dfechaFin_input"
	searchButtonSelector               = "tbBuscador:idFormBuscarProceso:btnBuscarSelToken"
	retrievedRowsDataContainerSelector = ".ui-paginator-current"
	advancedSearchSelector             = ".ui-fieldset-legend"
	tableDataSelector                  = "tbBuscador:idFormBuscarProceso:dtProcesos_data"
	nextPageButton                     = ".ui-paginator-next"
	previousPageButton                 = ".ui-paginator-prev"
	selectionProceduresTabButton       = "/html/body/div[3]/div/div[1]/ul/li[2]"
	selectionProceduresTabID           = "tbBuscador:tab1"

	// Constantes de tiempo
	initialWaitTime     = 2 * time.Second
	searchWaitTime      = 10 * time.Second
	pageLoadTimeout     = 10 * time.Second
	elementWaitTimeout  = 30 * time.Second
	pageNavigationDelay = 5 * time.Second

	// Mensajes de error
	errIniciarServicioSelenium = "Error al iniciar el servicio de Selenium"
	errAbrirNavegador          = "Error al abrir el navegador"
	errEncontrarTab            = "Error al encontrar el tab de procedimientos de selección"
	errHacerClicTab            = "Error al hacer clic en el tab de procedimientos de selección"
	errObtenerTab              = "Error al obtener el tab de procedimientos de selección"
	errRellenarFechas          = "Error al rellenar las fechas"
	errDesplazarsePagina       = "Error al desplazarse al final de la página"
	errEncontrarFilas          = "Error al encontrar la cantidad total de filas"
	errEsperarCargaPagina      = "Error al esperar a que la página se cargue"
	errExtraerIdentificador    = "Error al extraer el identificador de fila"
	errSeleccionarElemento     = "Error al seleccionar el elemento"
	errNoRegistros             = "No se obtuvieron registros"
	errProcesarRegistro        = "Error al procesar el registro"
	screenshotsDir             = "screenshots"
)

func Start(date time.Time) {
	logger := log.New(os.Stderr, "[scrapper] ", log.LstdFlags)
	logger.Printf("Proceso inicializado para la fecha: %s\n", date.Format("2006-01-02"))

	service, err := selenium.NewChromeDriverService("chromedriver", 4444)
	if err != nil {
		logger.Printf("%s:\n%v", errIniciarServicioSelenium, err)
		return
	}
	defer service.Stop()

	driver, err := setupDriver()
	if err != nil {
		logger.Printf("%s:\n%v", errAbrirNavegador, err)
		return
	}

	if err := inicializarTabProcedimientos(driver, logger); err != nil {
		return
	}

	tab, err := getSelectionProcessTab(driver)
	if err != nil {
		logger.Printf("%s:\n%v", errObtenerTab, err)
		return
	}

	if err := realizarBusqueda(driver, tab, date, logger); err != nil {
		return
	}

	tab, err = getSelectionProcessTab(driver)
	if err != nil {
		logger.Printf("%s:\n%v", errObtenerTab, err)
		return
	}

	recordsObtained, err := findTotalAmountOfRows(tab)
	if err != nil {
		logger.Printf("%s:\n%v", errEncontrarFilas, err)
		return
	}

	if err := procesarRegistros(driver, recordsObtained, logger); err != nil {
		return
	}

	logger.Println("Proceso finalizado exitosamente")
}

func inicializarTabProcedimientos(driver selenium.WebDriver, logger *log.Logger) error {
	button, err := driver.FindElement(selenium.ByXPATH, selectionProceduresTabButton)
	if err != nil {
		logger.Printf("%s:\n%v", errEncontrarTab, err)
		return err
	}

	if err := button.Click(); err != nil {
		logger.Printf("%s:\n%v", errHacerClicTab, err)
		return err
	}

	time.Sleep(initialWaitTime)
	return nil
}

func realizarBusqueda(driver selenium.WebDriver, tab selenium.WebElement, date time.Time, logger *log.Logger) error {
	if err := fillDates(tab, date); err != nil {
		logger.Printf("%s:\n%v", errRellenarFechas, err)
		return err
	}

	time.Sleep(searchWaitTime)

	if _, err := driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil); err != nil {
		logger.Printf("%s:\n%v", errDesplazarsePagina, err)
		return err
	}

	return nil
}

func takeScreenshot(driver selenium.WebDriver, filename string) error {
	// Crear directorio si no existe
	if err := os.MkdirAll(screenshotsDir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio de screenshots:\n%w", err)
	}

	// Tomar screenshot
	screenshot, err := driver.Screenshot()
	if err != nil {
		return fmt.Errorf("error al tomar screenshot:\n%w", err)
	}

	// Guardar screenshot
	path := filepath.Join(screenshotsDir, filename)
	if err := os.WriteFile(path, screenshot, 0644); err != nil {
		return fmt.Errorf("error al guardar screenshot:\n%w", err)
	}

	return nil
}

func procesarRegistros(driver selenium.WebDriver, recordsObtained int64, logger *log.Logger) error {
	if recordsObtained == 0 {
		logger.Println(errNoRegistros)
		return nil
	}

	logger.Printf("Cantidad total de filas obtenidas: %d\n", recordsObtained)

	if err := waitForTableToLoad(driver); err != nil {
		logger.Printf("%s:\n%v", errEsperarCargaPagina, err)
		return err
	}

	rowIdentifierFormat, err := extractRowIdentifierFormat(driver)
	if err != nil {
		logger.Printf("%s:\n%v", errExtraerIdentificador, err)
		return err
	}
	logger.Printf("Formato de identificador de fila extraído: %s\n", rowIdentifierFormat)

	printHeader()
	for i := 0; i < int(recordsObtained); i++ {
		logger.Printf("Procesando registro %d de %d\n", i+1, recordsObtained)

		tab, err := getSelectionProcessTab(driver)
		if err != nil {
			logger.Printf("%s:\n%v", errObtenerTab, err)
			return err
		}

		if err := selectElement(driver, tab, i, rowIdentifierFormat, logger); err != nil {
			logger.Printf("%s %d:\n%v", errProcesarRegistro, i+1, err)

			// Tomar screenshot del error
			screenshotName := fmt.Sprintf("error_registro_%d_%s.png", i+1, time.Now().Format("20060102_150405"))
			if screenshotErr := takeScreenshot(driver, screenshotName); screenshotErr != nil {
				logger.Printf("Error al tomar screenshot del error: %v", screenshotErr)
			} else {
				logger.Printf("Screenshot guardado como: %s", screenshotName)
			}

			break
		}

		logger.Printf("Registro %d procesado correctamente\n", i+1)
	}

	return nil
}

func waitForTableToLoad(driver selenium.WebDriver) error {
	return driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		_, err := wd.FindElement(selenium.ByID, tableDataSelector)
		return err == nil, nil
	}, pageLoadTimeout)
}

func getSelectionProcessTab(driver selenium.WebDriver) (selenium.WebElement, error) {
	// Esperar a que el tab esté presente y sea interactivo
	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		// Intentar encontrar el tab
		tab, err := wd.FindElement(selenium.ByID, selectionProceduresTabID)
		if err != nil {
			return false, nil // No es un error, solo que aún no está disponible
		}

		// Verificar si el elemento es visible y habilitado
		visible, err := tab.IsDisplayed()
		if err != nil || !visible {
			return false, nil
		}

		enabled, err := tab.IsEnabled()
		if err != nil || !enabled {
			return false, nil
		}

		return true, nil
	}, elementWaitTimeout)

	if err != nil {
		return nil, fmt.Errorf("timeout esperando el tab de procedimientos de selección:\n%w", err)
	}

	// Una vez que sabemos que el elemento está disponible, lo obtenemos
	tab, err := driver.FindElement(selenium.ByID, selectionProceduresTabID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo encontrar el formulario de procedimientos de selección:\n%w", err)
	}

	return tab, nil
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

func fillDates(tab selenium.WebElement, date time.Time) error {
	advancedSearchButton, err := tab.FindElement(selenium.ByCSSSelector, advancedSearchSelector)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el botón de búsqueda avanzada:\n%w", err)
	}

	err = advancedSearchButton.Click()
	if err != nil {
		return fmt.Errorf("no se pudo hacer clic en el botón de búsqueda avanzada:\n%w", err)
	}

	time.Sleep(2 * time.Second)
	formattedDate := date.Format("02/01/2006")
	startDateSelector, err := tab.FindElement(selenium.ByID, startDateSelector)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el selector de fecha de inicio:\n%w", err)
	}
	err = startDateSelector.SendKeys(formattedDate)
	if err != nil {
		return fmt.Errorf("no se pudo establecer el valor de la fecha de inicio:\n%w", err)
	}

	time.Sleep(2 * time.Second)
	endDateSelector, err := tab.FindElement(selenium.ByID, endDateSelector)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el selector de fecha de fin:\n%w", err)
	}
	err = endDateSelector.SendKeys(formattedDate)
	if err != nil {
		return fmt.Errorf("no se pudo establecer el valor de la fecha de fin:\n%w", err)
	}

	button, err := tab.FindElement(selenium.ByID, searchButtonSelector)
	if err != nil {
		return fmt.Errorf("no se pudo obtener el botón de búsqueda:\n%w", err)
	}
	err = button.Click()
	if err != nil {
		return fmt.Errorf("no se pudo hacer clic en el botón de búsqueda:\n%w", err)
	}

	return nil
}

func findTotalAmountOfRows(tab selenium.WebElement) (int64, error) {
	retrievedRowsData, err := tab.FindElement(selenium.ByCSSSelector, retrievedRowsDataContainerSelector)
	if err != nil {
		return 0, fmt.Errorf("no se pudo obtener el contenedor de filas recuperadas:\n%s", err)
	}

	text, err := retrievedRowsData.Text()
	if err != nil {
		return 0, fmt.Errorf("no se pudo obtener el texto del contenedor de filas recuperadas:\n%s", err)
	}

	parts := strings.Fields(text)
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

func selectElement(driver selenium.WebDriver, tab selenium.WebElement, id int, rowIdentifierFormat string, logger *log.Logger) error {
	err := goToPage(driver, tab, calculatePageNumber(id), logger)
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

func goToPage(driver selenium.WebDriver, tab selenium.WebElement, page int, logger *log.Logger) error {
	paginator, err := tab.FindElements(selenium.ByCSSSelector, ".ui-paginator-page")
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
		err = clickNextPage(tab)
		if err != nil {
			return fmt.Errorf("error al hacer clic en el botón de siguiente página:\n%s", err)
		}
		_, err = driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
		if err != nil {
			return err
		}
		return goToPage(driver, tab, page, logger)
	}

	logger.Println("retrocediendo")
	err = clickPreviousPage(tab)
	if err != nil {
		return fmt.Errorf("error al hacer clic en el botón de página anterior:\n%s", err)
	}
	return goToPage(driver, tab, page, logger)
}

func clickNextPage(tab selenium.WebElement) error {
	nextPage, err := tab.FindElement(selenium.ByCSSSelector, nextPageButton)
	if err != nil {
		return fmt.Errorf("error al obtener el botón de siguiente página:\n%s", err)
	}

	classNames, err := nextPage.GetAttribute("class")
	if err != nil {
		return fmt.Errorf("error al obtener los nombres de clase:\n%s", err)
	}

	if strings.Contains(classNames, "ui-state-disabled") {
		return fmt.Errorf("el botón de siguiente página está deshabilitado")
	}

	err = nextPage.Click()
	if err != nil {
		return fmt.Errorf("error al hacer clic en el botón de siguiente página:\n%s", err)
	}
	time.Sleep(5 * time.Second)
	return nil
}

func clickPreviousPage(tab selenium.WebElement) error {
	previousPage, err := tab.FindElement(selenium.ByCSSSelector, previousPageButton)
	if err != nil {
		return fmt.Errorf("error al obtener el botón de página anterior:\n%s", err)
	}

	classNames, err := previousPage.GetAttribute("class")
	if err != nil {
		return fmt.Errorf("error al obtener los nombres de clase:\n%s", err)
	}

	if strings.Contains(classNames, "ui-state-disabled") {
		return fmt.Errorf("el botón de página anterior está deshabilitado")
	}

	err = previousPage.Click()
	if err != nil {
		return fmt.Errorf("error al hacer clic en el botón de página anterior:\n%s", err)
	}
	time.Sleep(5 * time.Second)
	return nil
}
