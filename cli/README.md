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

The `runner.sh` requires the scrubber to be installed on `/opt/custom`. It's a wrapper for
the cron job so that it's log file and report file all get send to the same place.
