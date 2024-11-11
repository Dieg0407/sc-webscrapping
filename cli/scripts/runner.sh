#!/bin/sh

last_week=$(date -d "last week" +"%Y-%m-%d")
log_path="/var/log/scrubber/$last_week.execution.log"
report_path="/var/custom/scrubber/reportes-$last_week.csv"

/opt/custom/scrubber -d $last_week > $report_path 2> $log_path
