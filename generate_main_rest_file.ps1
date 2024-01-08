$copyFlag = $false
$copyMarker = "//COPY-FROM-HERE"
$firstLine = "//ATTNETION: THIS FILE IS AUTO GENERATED DON'T EDIT IT HERE - EDIT IN THE mock\requests folder"

# Clear the existing main.rest file or create an empty one
$firstLine > mock\main.auto.rest

# Define a list of files to process
$files = "const.rest", "auth.rest", "users.rest", "hotels.rest", "booking.rest"

# Loop through each file in the list
foreach ($file in $files) {
  $copyFlag = $false
  Get-Content "mock\requests\$file" | ForEach-Object {
    if ($copyFlag -eq $false -and $_ -eq $copyMarker) {
      $copyFlag = $true
    }
    elseif ($copyFlag -eq $true) {
      $_
    }
  } | Out-File -Append mock\main.auto.rest
}