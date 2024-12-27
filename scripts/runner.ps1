$execution_date = $args[0]

# Check if date is empty
if ([string]::IsNullOrEmpty($execution_date)) {
  $execution_date = (Get-Date).AddDays(-7).ToString("yyyy-MM-dd")
}

# Check if "logs" directory exists, create it if not
if (!(Test-Path -Path "logs")) {
  New-Item -ItemType Directory -Path "logs"
}

# Check if "reports" directory exists, create it if not
if (!(Test-Path -Path "reports")) {
  New-Item -ItemType Directory -Path "reports"
}

$log_path = "logs\$execution_date.execution.log"
$report_path = "reports\reportes-$execution_date.csv"

.\scrubber.exe -d $execution_date | Out-File -FilePath $report_path -Encoding utf8 2>&1 | Out-File -FilePath $log_path -Encoding utf8
