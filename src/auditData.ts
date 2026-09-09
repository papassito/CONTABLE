import { AuditItem, ForensicReport } from './types';

export const forensicData: ForensicReport = {
  encabezado_auditoria: {
    fecha_ejecucion_utc: "2026-09-06T16:30:28Z",
    sistema_operativo: "Microsoft Windows 11 Pro",
    arquitectura: "64 bits",
    nombre_equipo: "DEV-STATION-01"
  },
  evidencia_archivos: [
    {
      ruta: "C:\\Windows\\System32\\drivers\\etc\\hosts",
      tamaño_bytes: 907,
      sha256: "B78800996089EAB544F838A14405A8504D9466005A177730331AF6D209318FA4",
      creacion_utc: "2022-05-07T05:24:58Z",
      ultima_modificacion_utc: "2026-08-27T19:22:37Z",
      ultimo_acceso_utc: "2026-09-06T15:57:03Z"
    },
    {
      ruta: "C:\\Windows\\System32\\drivers\\etc\\hosts.ics",
      tamaño_bytes: 439,
      sha256: "EC80311A1ED7FB7CB5AF9A602AD3D7C43ECAE2B3151FB008920BF7AF982CD9E5",
      creacion_utc: "2026-05-28T06:23:06Z",
      ultima_modificacion_utc: "2026-05-29T06:48:05Z",
      ultimo_acceso_utc: "2026-09-06T07:35:08Z"
    },
    {
      ruta: "C:\\Windows\\System32\\drivers\\etc\\lmhosts.sam",
      tamaño_bytes: 3683,
      sha256: "000FE9C924B4D155477CAD15B4CFD30616C37523B4B848D6ECBD003507A55EDF",
      creacion_utc: "2024-04-01T07:26:16Z",
      ultima_modificacion_utc: "2024-04-01T07:24:05Z",
      ultimo_acceso_utc: "2024-04-01T07:24:05Z"
    },
    {
      ruta: "C:\\Windows\\System32\\drivers\\etc\\networks",
      tamaño_bytes: 407,
      sha256: "21DE93ED8293DBB9C53B59A5C1AF04B1FD997CF7DFBD0BA5F21CB315D845B7A8",
      creacion_utc: "2022-05-07T05:24:58Z",
      ultima_modificacion_utc: "2022-05-07T05:22:33Z",
      ultimo_acceso_utc: "2026-09-06T09:31:21Z"
    },
    {
      ruta: "C:\\Windows\\System32\\drivers\\etc\\protocol",
      tamaño_bytes: 1358,
      sha256: "52CF86496F3859D0F3E58776ECCFF1D6589792004D2291D3B0CE8D7635BE7278",
      creacion_utc: "2022-05-07T05:24:58Z",
      ultima_modificacion_utc: "2022-05-07T05:22:33Z",
      ultimo_acceso_utc: "2026-09-06T09:31:21Z"
    },
    {
      ruta: "C:\\Windows\\System32\\drivers\\etc\\services",
      tamaño_bytes: 17635,
      sha256: "C1723F7F29B224C42F26452C3EFA8F80F6EE8500EE78513E0C0732BA55399F7D",
      creacion_utc: "2022-05-07T05:24:58Z",
      ultima_modificacion_utc: "2022-05-07T05:22:33Z",
      ultimo_acceso_utc: "2026-09-06T16:01:13Z"
    }
  ]
};

export const auditItems: AuditItem[] = [
  {
    id: 'sanation-duplicates',
    category: 'workspace',
    title: 'Saneamiento y Purga de Archivos Duplicados',
    severity: 'success',
    status: 'saneado',
    summary: 'Se eliminaron exitosamente los 8 archivos duplicados y desalineados introducidos por la subida.',
    technicalDetails: 'Archivos eliminados: `.env-1.example`, `-1.gitignore`, `bun-1.lock`, `index-1.html`, `metadata-1.json`, `package-1.json`, `tsconfig-1.json` y `vite.config-1.ts`. El árbol de archivos ahora preserva la estructura canónica intacta.',
    targetFiles: ['.env-1.example', '-1.gitignore', 'bun-1.lock', 'index-1.html', 'package-1.json', 'tsconfig-1.json', 'vite.config-1.ts'],
    suggestedActionForCodeAssist: 'No requiere acción adicional. La integridad estructural del workspace quedó normalizada.'
  },
  {
    id: 'go-backend-missing',
    category: 'go-backend',
    title: 'Integridad del Módulo Go Verificada',
    severity: 'success',
    status: 'saneado',
    summary: 'El módulo Go está presente en la Raíz, completo y pasa todas las pruebas y análisis estáticos (`go vet`).',
    technicalDetails: 'Se ha verificado la presencia de `go.mod` y la estructura de paquetes de dominio. Los scripts de prueba (`run-tests.ps1`) y auditoría (`audit_fcos_anatomy.ps1`) se ejecutan correctamente, confirmando la salud del backend.',
    targetFiles: ['go.mod', 'run-tests.ps1', 'audit_fcos_anatomy.ps1'],
    suggestedActionForCodeAssist: 'No requiere acción. El backend de Go está en un estado saludable y validado.'
  },
  {
    id: 'go-contracts-repair',
    category: 'go-backend',
    title: 'Saneamiento de Contratos de Dominio Duplicados',
    severity: 'success',
    status: 'saneado',
    summary: 'Se eliminaron los archivos de dominio conflictivos que requerían la ejecución de scripts de reparación.',
    technicalDetails: 'Se eliminaron `internal/domain/vault.go` e `internal/domain/secret_envelope.go`. Estos archivos eran placeholders que causaban errores de compilación (EOF) y creaban inconsistencias con las definiciones correctas en `internal/secrets/domain`.',
    targetFiles: ['internal/domain/vault.go', 'internal/domain/secret_envelope.go'],
    suggestedActionForCodeAssist: 'No requiere acción. La arquitectura de dominio ha sido consolidada. Los scripts de reparación ya no son necesarios para la compilación.'
  },
  {
    id: 'forensic-hashes',
    category: 'forensic',
    title: 'Auditoría Forense de Integridad del Sistema',
    severity: 'info',
    status: 'verificado',
    summary: 'Evaluación del reporte forense emitido en la máquina DEV-STATION-01 (Windows 11 Pro).',
    technicalDetails: 'Se auditaron los hashes SHA-256 de 6 archivos de red de Windows (hosts, hosts.ics, lmhosts.sam, networks, protocol, services). La modificación reciente en hosts (2026-08-27) y hosts.ics (2026-05-29) es congruente con entornos de desarrollo locales con resolución DNS local o Docker/Wsl.',
    targetFiles: ['informe_forense.json', 'Scanner-Agresivo.ps1'],
    suggestedActionForCodeAssist: 'Revisar si existen líneas no comentadas en C:\\Windows\\System32\\drivers\\etc\\hosts durante la ejecución de Scanner-Agresivo.ps1 paso 5/6.'
  },
  {
    id: 'scanner-agresivo-analysis',
    category: 'powershell',
    title: 'Análisis Estático del Scanner-Agresivo.ps1',
    severity: 'info',
    status: 'verificado',
    summary: 'Script de solo lectura con 6 fases de escrutinio profundo del sistema.',
    technicalDetails: 'Fases: 1) Windows Defender & RTP. 2) Procesos en \\AppData\\ y \\Temp\\. 3) Conexiones TCP ESTABLISHED no locales. 4) Llaves de persistencia en Run (HKCU y HKLM). 5) Redirecciones activas en HOSTS. 6) Top 5 procesos de alto consumo de CPU.',
    targetFiles: ['Scanner-Agresivo.ps1'],
    suggestedActionForCodeAssist: 'El script es puramente de lectura (no destructivo). Puede ser ejecutado en la máquina host con privilegios de Administrador para confirmar diagnóstico.'
  },
  {
    id: 'wails-config-audit',
    category: 'wails',
    title: 'Configuración de Wails v2 (Contable Fix FCOS v2.2)',
    severity: 'info',
    status: 'verificado',
    summary: 'wails.json alineado con el empaquetador pnpm y salida NSIS.',
    technicalDetails: 'Configuración: name: "ContableFix", version: "2.2.0", build: pnpm build, watcher: pnpm dev, wailsjsdir: "src/wailsjs", assetdir: "dist". La carpeta "src/wailsjs" aún no ha sido generada por el CLI de Wails.',
    targetFiles: ['wails.json', 'package.json'],
    suggestedActionForCodeAssist: 'Ejecutar `wails dev` o `wails generate module` en el entorno con Wails instalado para vincular los bindings TS de Go en `src/wailsjs`.'
  },
  {
    id: 'cicd-signtool-path',
    category: 'cicd',
    title: 'Resolución Dinámica de SignTool en CI/CD',
    severity: 'success',
    status: 'saneado',
    summary: 'El pipeline de GitHub Actions ya no utiliza una ruta fija para signtool.exe, evitando fallos futuros.',
    technicalDetails: 'El script de `release.yml` ahora busca dinámicamente la última versión de `signtool.exe` dentro de `C:\\Program Files (x86)\\Windows Kits`. Esto garantiza que la firma digital funcione incluso si GitHub actualiza la versión del Windows SDK en sus runners `windows-latest`.',
    targetFiles: ['release.yml'],
    suggestedActionForCodeAssist: 'No requiere acción. La robustez del pipeline de release ha sido mejorada.'
  },
  {
    id: 'pkg-pnpm-consistency',
    category: 'workspace',
    title: 'Coherencia de Gestores (pnpm vs npm vs bun)',
    severity: 'info',
    status: 'verificado',
    summary: 'Archivos de configuración pnpm-workspace.yaml y .npmrc auditados.',
    technicalDetails: '.npmrc define `only-built-dependencies[]` para @google/genai, esbuild y protobufjs. pnpm-workspace.yaml restringe esbuild y protobufjs. Se recomienda mantener consistencia en CI usando siempre pnpm v9.',
    targetFiles: ['.npmrc', 'pnpm-workspace.yaml', 'pnpm-lock.yaml', 'package.json'],
    suggestedActionForCodeAssist: 'Asegurar que los desarrolladores y el CI usen `pnpm install --frozen-lockfile`.'
  }
];

export const generateCodeAssistReportMarkdown = (): string => {
  return `# DICTAMEN DE AUDITORÍA AGRESIVA Y MILIMÉTRICA
**Proyecto:** Contable Fix FCOS v2.2  
**Fecha:** ${new Date().toISOString().split('T')[0]}  
**Destinatario:** Code Assist / Equipo de Desarrollo  

---

### 1. SANEAMIENTO DE DUPLICADOS (RESUELTO ✅)
- Se eliminaron con éxito los 8 archivos espurios y duplicados:
  - \`.env-1.example\`
  - \`-1.gitignore\`
  - \`bun-1.lock\`
  - \`index-1.html\`
  - \`metadata-1.json\`
  - \`package-1.json\`
  - \`tsconfig-1.json\`
  - \`vite.config-1.ts\`
- El árbol de configuración en la raíz se encuentra limpio y validado.

---

### 2. AUDITORÍA DEL BACKEND GO (RESUELTO ✅)
- **Estado:** El módulo Go reside en la Raíz del proyecto, completo y validado de forma nativa con Wails.
- **Verificación:** Pasa con éxito las pruebas (\`run-tests.ps1\`) y los análisis estáticos (\`go vet\`, \`audit_fcos_anatomy.ps1\`).
- **Arquitectura:** La estructura de dominio ha sido consolidada, eliminando archivos duplicados y resolviendo inconsistencias. El backend está en un estado saludable.

---

### 3. AUDITORÍA CI/CD Y WAILS (VERIFICADO ✅)
- **Archivo:** \`release.yml\`
- **Estado:** El pipeline de release ya resuelve dinámicamente la ruta de \`signtool.exe\`, lo que lo hace robusto ante actualizaciones del runner \`windows-latest\`.
- **Wails:** \`wails.json\` está configurado correctamente con \`pnpm\`. El siguiente paso es ejecutar \`wails dev\` para generar los bindings en \`src/wailsjs\`.

---

### 4. INFORME FORENSE DEL SISTEMA (DEV-STATION-01)
- **Equipo:** Windows 11 Pro 64 bits (\`DEV-STATION-01\`)
- **Archivos verificados:** 6 archivos críticos en \`C:\\Windows\\System32\\drivers\\etc\\\` con sus hashes SHA-256 registrados.
- **Scanner-Agresivo.ps1:** 6 rutinas de inspección (Defender, procesos anómalos en Temp/AppData, conexiones TCP salientes, claves Run de persistencia, integridad de HOSTS y procesos demandantes).
`;
};
