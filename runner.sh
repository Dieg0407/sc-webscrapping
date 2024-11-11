#!/bin/bash

last_week=$(date -d "last week" +"%Y-%m-%d")
log_path="/var/log/scrubber/$last_week.log"
report_path="/var/custom/scrubber/reportes-$last_week.csv"

/opt/custom/scrubber -d $last_week > reportes-$last_week.csv 2> $last_week.execution.log
