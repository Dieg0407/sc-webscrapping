$execution_date = $args[0]

# Comprobar si la fecha está vacía
if ([string]::IsNullOrEmpty($execution_date)) {
  $execution_date = (Get-Date).AddDays(-7).ToString("yyyy-MM-dd")
}

$log_path = "logs\$execution_date.execution.log"
$report_path = "reportes\reportes-$execution_date.csv"

scrubber.exe -d $execution_date | Out-File -FilePath $report_path -Encoding utf8 2>&1 | Out-File -FilePath $log_path -Encoding utf8
