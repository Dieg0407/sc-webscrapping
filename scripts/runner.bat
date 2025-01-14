@echo off

set /p execution_date=Ingresa la fecha a procesar (yyyy-MM-dd): 

if not exist reportes mkdir reportes

scrapper.exe -d %execution_date% > reportes\reporte-%execution_date%.csv
