#!/bin/bash

last_week=$(date -d "last week" +"%Y-%m-%d")
sc-scrubber -d $last_week > reportes-$last_week.csv 2> $last_week.execution.log
