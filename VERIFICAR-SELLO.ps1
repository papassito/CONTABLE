param([string]$ProjectPath = 'C:\Users\Radio 2027\Desktop\Soft Septioembre 2026\CONTABLE')
$ErrorActionPreference = 'Stop'
$manifest = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'SELLO-MD.json') -Raw -Encoding utf8 | ConvertFrom-Json
$failed = $false
foreach ($entry in $manifest.files) {
    $target = Join-Path $ProjectPath $entry.file
    if (-not (Test-Path -LiteralPath $target -PathType Leaf)) {
        Write-Output "FALTA: $($entry.file)"
        $failed = $true
        continue
    }
    $actual = (Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash
    if ($actual -ne $entry.sha256) {
        Write-Output "MODIFICADO: $($entry.file)"
        $failed = $true
    } else {
        Write-Output "OK: $($entry.file)"
    }
}
if ($failed) { throw 'El proyecto no coincide con el sello documental.' }
Write-Output "SELLO VERIFICADO: $($manifest.version)"
