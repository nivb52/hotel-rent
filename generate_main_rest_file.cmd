@echo off
setlocal enabledelayedexpansion

set "copyFlag="
set "copyMarker=//COPY-FROM-HERE"
set "firstLine=//ATTNETION: THIS FILE IS AUTO GENERATED DON'T EDIT IT HERE - EDIT IN THE mock\requests folder"
:: Define a list of files to process
set "files=const auth users hotels booking"
(
  for %%f in (%files%) do (
    echo %firstLine%
    for /f "delims=" %%a in (mock\requests\%%f.rest) do (
      set "line=%%a"
      if "!line!"=="%copyMarker%" (
        set "copyFlag=1"
        echo  //
        echo  //
        echo //////  %%f  ////// 
      ) else if defined copyFlag (
        echo !line!
      )
    )
    set "copyFlag="
  )
) > mock\main_auto_generated.rest

::   echo hotels

::   type mock\requests\hotels.rest