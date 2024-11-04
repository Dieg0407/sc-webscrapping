# About
This is a small project meant to create a scrapper with GOLANG that uses
chromedriver and selenium as a driver.

## Usage

This is a cli tool, following conventions defined [here](https://clig.dev/#the-basics), 
the debug output will be sent to `stderr` while the useful data will be sent to
`stdout`.

To `stdout`, the program will send the winners of the evaluated processes separated by the
character `;` with the first row being always the header.

```
Identificador;Entidad;Nomenclarura;Objecto;Descripción;Valor;Moneda;Ganador;Es MYPE;Es Selva
```

You can then run the program with the command sending the date in the format `YYYY-MM-DD`

```bash
go run cmd/cli.go -d "2024-11-01" > reportes-2024-11-01.csv
```

## Scripts

The `install.sh` will build the program and add it to the `/usr/local/bin/` folder.

The `runner.sh` requires install to be ran first, it will run the program taking as the 
input date the week previous to the execution date
