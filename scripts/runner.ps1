$execution_date = Read-Host "Ingresa la fecha a procesar (yyyy-MM-dd): "

if (!(Test-Path -Path "reportes")) {
  New-Item -ItemType Directory -Path "reportes"
}

$report_path = "reportes\reporte-$execution_date.csv"

.\scrubber.exe -d $execution_date > $report_path
