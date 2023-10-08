$copyFlag = $false
$copyMarker = "//COPY-FROM-HERE"

# Clear the existing main.rest file or create an empty one
$null > mock\main_auto_generated.rest

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
  } | Out-File -Append mock\main_auto_generated.rest
}