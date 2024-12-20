#!/bin/sh

execution_date=$1

# check if date is empty
if [ -z "$execution_date" ]; then
    execution_date=$(date -d "last week" +"%Y-%m-%d")
fi

log_path="/var/log/scrubber/$execution_date.execution.log"
report_path="/var/custom/scrubber/reportes-$execution_date.csv"

/opt/custom/scrubber -d $execution_date > $report_path 2> $log_path
