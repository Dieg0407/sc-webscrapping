# About
This is a small project meant to create a scrapper with GOLANG that uses
chromedriver and selenium as a driver.

**chromedriver** must be on the PATH.

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
go build -v -o scrubber cmd/cli.go
./scrubber -d "2024-11-01" > reportes-2024-11-01.csv
```


## Scripts

There are 2 main scripts for this, one for windows and one for linux. Both require the `scrubber`
to be on the same folder

The linux `runner.sh` is meant to be used as a cron, by default it will take the current 
date and substract a week. If you want to override this behaviour then send the date
in the format yyyy-MM-dd

```bash
./runner.sh 2024-11-01 &
tail -f logs/2024-11-01.execution.log
```

The windows one `runner.ps1` is meant to be used as a cli tool to extract data.
It will require to pass the date as an input and will not default to the previous week as date.

```powershell
.\runner.ps1
```
