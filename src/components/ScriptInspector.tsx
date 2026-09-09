import { useState } from 'react';
import { FileCode2, Copy, Check, Terminal, Shield, Play } from 'lucide-react';

interface ScriptFile {
  id: string;
  name: string;
  type: string;
  description: string;
  code: string;
}

const scriptsData: ScriptFile[] = [
  {
    id: 'scanner',
    name: 'Scanner-Agresivo.ps1',
    type: 'PowerShell',
    description: 'Diagnóstico forense no destructivo (procesos en Temp/AppData, conexiones TCP, persistencia, HOSTS, CPU).',
    code: `<#
.SYNOPSIS
    Scanner Agresivo de Diagnóstico (Solo Lectura)
.DESCRIPTION
    Inspecciona procesos, red, persistencia y estado de seguridad sin modificar el sistema.
#>

$ErrorActionPreference = "SilentlyContinue"

Write-Host " 🛡️  INICIANDO SCANNER AGRESIVO DE DIAGNÓSTICO (MODO: SOLO LECTURA)  🛡️ "
# [1/6] VERIFICANDO ESTADO DE WINDOWS DEFENDER / SEGURIDAD...
$defender = Get-MpComputerStatus

# [2/6] CAZANDO PROCESOS SOSPECHOSOS (Paths Inusuales)...
$suspiciousProcesses = Get-Process | Where-Object { 
    $_.Path -match "\\\\AppData\\\\" -or $_.Path -match "\\\\Temp\\\\" 
}

# [3/6] ANALIZANDO CONEXIONES DE RED ACTIVAS (ESTABLISHED)...
$connections = Get-NetTCPConnection | Where-Object { 
    $_.State -eq 'Established' -and $_.RemoteAddress -notmatch "^127\.|^192\.168\.|^10\.|^172\.(1[6-9]|2[0-9]|3[0-1])\.|^::1" 
}

# [4/6] REVISANDO REGISTROS DE ARRANQUE (PERSISTENCIA)...
$runPaths = @("HKCU:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", "HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")

# [5/6] VERIFICANDO INTEGRIDAD DEL ARCHIVO HOSTS...
$hostsPath = "$env:windir\\System32\\drivers\\etc\\hosts"

# [6/6] DETECTANDO PROCESOS DE ALTO CONSUMO (CPU/RAM)...
$highHogs = Get-Process | Sort-Object CPU -Descending | Select-Object -First 5`
  },
  {
    id: 'reparar',
    name: 'reparar_y_probar.ps1',
    type: 'PowerShell',
    description: 'Localizador de go.mod, reparación de contrato document_parser.go, purga de config/uow.go y ejecución de pruebas.',
    code: `<#
.SYNOPSIS
    Script de reparación y ejecución de pruebas para el proyecto Contable Go.
#>

# 1. Determinar rutas candidatas donde puede habitar 'go.mod'
$projectDir = $basePath

# 2. Eliminar el archivo duplicado y vacío config\\uow.go
$uowFile = Join-Path -Path $projectDir -ChildPath "config\\uow.go"
if (Test-Path -Path $uowFile) { Remove-Item -Path $uowFile -Force }

# 3. Asegurar contrato internal\\domain\\document_parser.go
$domainDir = Join-Path -Path $projectDir -ChildPath "internal\\domain"
$parserFile = Join-Path -Path $domainDir -ChildPath "document_parser.go"

# 4. Reubicar 'anchor.go' a internal\\domain\\
$wrongAnchor = Join-Path -Path $projectDir -ChildPath "internal\\service\\anchor.go"
$targetAnchor = Join-Path -Path $domainDir -ChildPath "anchor.go"

# 5. Ejecución: go mod tidy && go test -v ./...`
  },
  {
    id: 'run-tests',
    name: 'run-tests.ps1',
    type: 'PowerShell',
    description: 'Automatización de go test ./... y generación de reporte de cobertura HTML (go tool cover).',
    code: `<#
.SYNOPSIS
    Script de automatización para pruebas contables en Go y cobertura HTML.
#>
function Invoke-ContableTests {
    # 1. Buscar go.mod en Contable GO o go-contable
    # 2. Limpieza de coverage.out y coverage.html previos
    # 3. go test -v ./... @TestArgs
    # 4. Control estricto de $LASTEXITCODE
    # 5. go tool cover -html=coverage.out -o coverage.html
}`
  },
  {
    id: 'release',
    name: 'release.yml',
    type: 'GitHub Actions',
    description: 'Pipeline de CI/CD para compilar con Wails en Windows, firmar digitalmente con signtool.exe y publicar.',
    code: `name: Build and Release FCOS v2.2
on:
  push:
    tags: ['v*']

jobs:
  build-windows:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - uses: pnpm/action-setup@v3
        with: { version: 9 }
      - run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - run: pnpm install
      - run: wails build -platform windows/amd64 -clean -nsis
      - name: Firma Digital EV (signtool)
        run: |
          $signtool = (Get-ChildItem -Path "C:\\Program Files (x86)\\Windows Kits" -Filter "signtool.exe" -Recurse | Sort-Object LastWriteTime -Descending | Select-Object -First 1).FullName
          & $signtool sign /f \${{ secrets.CERT_PATH }} /p \${{ secrets.CERT_PASS }} /tr http://timestamp.digicert.com /td sha256 /fd sha256 build\\bin\\ContableFix.exe`
  },
  {
    id: 'wails',
    name: 'wails.json',
    type: 'JSON',
    description: 'Configuración de Wails v2: Contable Fix FCOS v2.2.0, empaquetado con pnpm y salida NSIS.',
    code: `{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "ContableFix",
  "outputfilename": "ContableFix",
  "frontend:dir": "./",
  "frontend:install": "pnpm install",
  "frontend:build": "pnpm build",
  "frontend:dev:watcher": "pnpm dev",
  "info": {
    "companyName": "ContableFix",
    "productName": "Contable Fix FCOS",
    "productVersion": "2.2.0"
  },
  "wailsjsdir": "src/wailsjs",
  "assetdir": "dist",
  "build:dir": "build",
  "nsis": { "runAfterFinish": true }
}`
  }
];

export default function ScriptInspector() {
  const [activeTab, setActiveTab] = useState<string>('scanner');
  const [copied, setCopied] = useState(false);

  const currentScript = scriptsData.find(s => s.id === activeTab) || scriptsData[0];

  const handleCopy = () => {
    navigator.clipboard.writeText(currentScript.code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
      <div className="border-b border-slate-800 p-3 sm:p-4 bg-slate-900/80 flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Terminal className="w-4 h-4 text-indigo-400" />
          <span className="text-sm font-semibold text-white">
            Explorador de Scripts y Configuración de Auditoría
          </span>
        </div>

        <button
          onClick={handleCopy}
          className="flex items-center gap-1.5 px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition text-xs font-mono"
        >
          {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
          <span>{copied ? 'Copiado' : 'Copiar Código'}</span>
        </button>
      </div>

      {/* Tabs */}
      <div className="flex overflow-x-auto border-b border-slate-800 bg-slate-950/40 text-xs scrollbar-thin">
        {scriptsData.map(script => (
          <button
            key={script.id}
            onClick={() => setActiveTab(script.id)}
            className={`px-4 py-2.5 font-mono whitespace-nowrap transition border-b-2 ${activeTab === script.id
              ? 'border-indigo-500 text-white bg-slate-800/40 font-medium'
              : 'border-transparent text-slate-400 hover:text-slate-200 hover:bg-slate-900/40'
              }`}
          >
            {script.name}
          </button>
        ))}
      </div>

      {/* Description */}
      <div className="p-3 bg-slate-950/70 border-b border-slate-800/60 text-xs text-slate-400 flex items-center gap-2">
        <span className="px-2 py-0.5 rounded bg-indigo-500/10 text-indigo-400 font-mono text-[10px]">
          {currentScript.type}
        </span>
        <span>{currentScript.description}</span>
      </div>

      {/* Code Display */}
      <div className="p-4 bg-slate-950 overflow-x-auto max-h-96 text-xs font-mono text-slate-300 leading-relaxed">
        <pre>{currentScript.code}</pre>
      </div>
    </div>
  );
}
