import { useState, useMemo } from 'react';
import {
  FolderTree,
  FileCode,
  Download,
  Copy,
  Check,
  Search,
  Layers,
  Sparkles,
  ExternalLink,
  Code2,
  FolderGit2,
  ShieldAlert
} from 'lucide-react';
import { GO_SKELETON_FILES, GoFile } from './data/goSkeleton';
import { CodeViewer } from './components/CodeViewer';
import { ArchitectureView } from './components/ArchitectureView';
import { AuditView } from './components/AuditView';
import { downloadGoSkeletonZip } from './utils/zipExport';

type LayerFilter = 'all' | 'domain' | 'repository' | 'service' | 'handler' | 'config' | 'pkg' | 'root' | 'cmd';

export default function App() {
  const [selectedFilePath, setSelectedFilePath] = useState<string>('internal/domain/account.go');
  const [layerFilter, setLayerFilter] = useState<LayerFilter>('all');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [viewMode, setViewMode] = useState<'code' | 'architecture' | 'audit'>('code');
  const [isDownloading, setIsDownloading] = useState<boolean>(false);
  const [copiedTree, setCopiedTree] = useState<boolean>(false);

  // Find currently selected file
  const currentFile = useMemo(() => {
    return GO_SKELETON_FILES.find((f) => f.path === selectedFilePath) || GO_SKELETON_FILES[0];
  }, [selectedFilePath]);

  // Filter files based on layer and search
  const filteredFiles = useMemo(() => {
    return GO_SKELETON_FILES.filter((file) => {
      const matchesLayer = layerFilter === 'all' || file.layer === layerFilter;
      const matchesSearch =
        searchQuery.trim() === '' ||
        file.path.toLowerCase().includes(searchQuery.toLowerCase()) ||
        file.description.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesLayer && matchesSearch;
    });
  }, [layerFilter, searchQuery]);

  const handleDownloadZip = async () => {
    try {
      setIsDownloading(true);
      await downloadGoSkeletonZip();
    } catch (err) {
      console.error('Error al generar ZIP:', err);
    } finally {
      setIsDownloading(false);
    }
  };

  const handleCopyTree = () => {
    const treeText = `contable-fix/
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── domain/
│   │   ├── account.go
│   │   ├── journal.go
│   │   ├── ledger.go
│   │   ├── invoice.go
│   │   ├── audit.go
│   │   └── errors.go
│   ├── repository/
│   │   ├── account_repository.go
│   │   ├── journal_repository.go
│   │   ├── ledger_repository.go
│   │   └── invoice_repository.go
│   ├── service/
│   │   ├── account_service.go
│   │   ├── journal_service.go
│   │   └── ledger_service.go
│   └── handler/
│       └── http/
│           ├── account_handler.go
│           ├── journal_handler.go
│           └── router.go
├── pkg/
│   ├── response/
│   │   └── response.go
│   └── validator/
│       └── validator.go
├── Makefile
├── README.md
└── go.mod`;

    navigator.clipboard.writeText(treeText);
    setCopiedTree(true);
    setTimeout(() => setCopiedTree(false), 2000);
  };

  const layerBadges: { id: LayerFilter; label: string; count: number }[] = [
    { id: 'all', label: 'Todos', count: GO_SKELETON_FILES.length },
    { id: 'domain', label: 'Dominio', count: GO_SKELETON_FILES.filter((f) => f.layer === 'domain').length },
    { id: 'repository', label: 'Repositorios', count: GO_SKELETON_FILES.filter((f) => f.layer === 'repository').length },
    { id: 'service', label: 'Servicios', count: GO_SKELETON_FILES.filter((f) => f.layer === 'service').length },
    { id: 'handler', label: 'Handlers HTTP', count: GO_SKELETON_FILES.filter((f) => f.layer === 'handler').length },
    { id: 'cmd', label: 'cmd / api', count: GO_SKELETON_FILES.filter((f) => f.layer === 'cmd').length },
    { id: 'config', label: 'Config', count: GO_SKELETON_FILES.filter((f) => f.layer === 'config').length },
    { id: 'pkg', label: 'Utilidades', count: GO_SKELETON_FILES.filter((f) => f.layer === 'pkg').length },
  ];

  return (
    <div id="contable-fix-root" className="min-h-screen flex flex-col bg-slate-950 text-slate-100 font-sans antialiased">
      {/* Top Application Header */}
      <header className="border-b border-slate-800/80 bg-slate-900/90 backdrop-blur px-5 py-3 sticky top-0 z-30 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="w-9 h-9 rounded-lg bg-cyan-600/20 border border-cyan-500/40 flex items-center justify-center text-cyan-400 shadow-sm">
            <Code2 className="w-5 h-5" />
          </div>
          <div>
            <div className="flex items-center space-x-2">
              <h1 className="text-base font-bold tracking-tight text-white flex items-center gap-1.5">
                <span>Contable Fix</span>
                <span className="text-cyan-400 font-medium text-xs px-2 py-0.5 rounded bg-cyan-950/80 border border-cyan-800/50">
                  by KLIK
                </span>
              </h1>
            </div>
            <p className="text-xs text-slate-400">
              Esqueleto base limpio y modular en Go (Golang) • Vacío para llenarlo
            </p>
          </div>
        </div>

        {/* View and Export Actions */}
        <div className="flex items-center space-x-2">
          {/* Switch View Mode */}
          <div className="flex items-center bg-slate-950 rounded-lg p-1 border border-slate-800 text-xs">
            <button
              id="btn-view-code"
              onClick={() => setViewMode('code')}
              className={`px-3 py-1 rounded transition-colors flex items-center space-x-1.5 cursor-pointer ${
                viewMode === 'code' ? 'bg-cyan-600 text-white font-medium' : 'text-slate-400 hover:text-white'
              }`}
            >
              <FileCode className="w-3.5 h-3.5" />
              <span>Explorador de Código</span>
            </button>
            <button
              id="btn-view-architecture"
              onClick={() => setViewMode('architecture')}
              className={`px-3 py-1 rounded transition-colors flex items-center space-x-1.5 cursor-pointer ${
                viewMode === 'architecture' ? 'bg-cyan-600 text-white font-medium' : 'text-slate-400 hover:text-white'
              }`}
            >
              <Layers className="w-3.5 h-3.5" />
              <span>Arquitectura & Capas</span>
            </button>
            <button
              id="btn-view-audit"
              onClick={() => setViewMode('audit')}
              className={`px-3 py-1 rounded transition-colors flex items-center space-x-1.5 cursor-pointer ${
                viewMode === 'audit' ? 'bg-cyan-600 text-white font-medium' : 'text-slate-400 hover:text-white'
              }`}
            >
              <ShieldAlert className="w-3.5 h-3.5" />
              <span>Auditoría & Guía</span>
            </button>
          </div>

          {/* Copy Tree Text */}
          <button
            id="btn-copy-tree"
            onClick={handleCopyTree}
            className="flex items-center space-x-1.5 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700/80 text-slate-200 text-xs font-medium border border-slate-700/60 transition-colors cursor-pointer"
            title="Copiar estructura de directorios en texto"
          >
            {copiedTree ? (
              <>
                <Check className="w-3.5 h-3.5 text-emerald-400" />
                <span className="text-emerald-400">Árbol copiado</span>
              </>
            ) : (
              <>
                <Copy className="w-3.5 h-3.5 text-slate-300" />
                <span>Copiar Árbol</span>
              </>
            )}
          </button>

          {/* Download Full Zip */}
          <button
            id="btn-download-zip"
            onClick={handleDownloadZip}
            disabled={isDownloading}
            className="flex items-center space-x-1.5 px-3.5 py-1.5 rounded-lg bg-cyan-600 hover:bg-cyan-500 text-white text-xs font-semibold shadow-sm transition-colors cursor-pointer disabled:opacity-50"
            title="Descargar todos los archivos Go en un archivo ZIP listo para compilar"
          >
            <Download className="w-3.5 h-3.5" />
            <span>{isDownloading ? 'Generando ZIP...' : 'Descargar Esqueleto (.zip)'}</span>
          </button>
        </div>
      </header>

      {/* Main Workspace Layout */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left File Navigation Panel */}
        <aside className="w-80 shrink-0 border-r border-slate-800/80 bg-slate-900/50 flex flex-col overflow-hidden">
          {/* Search box */}
          <div className="p-3 border-b border-slate-800/80">
            <div className="relative">
              <Search className="w-3.5 h-3.5 text-slate-400 absolute left-2.5 top-2.5" />
              <input
                id="search-files-input"
                type="text"
                placeholder="Buscar archivo o struct..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-8 pr-3 py-1.5 bg-slate-950 border border-slate-800 rounded-md text-xs text-slate-200 placeholder:text-slate-400 focus:outline-none focus:border-cyan-500 transition-colors"
              />
            </div>
          </div>

          {/* Architectural Layer Badges / Filter */}
          <div className="p-3 border-b border-slate-800/80">
            <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2 flex items-center space-x-1">
              <Layers className="w-3 h-3 text-cyan-400" />
              <span>Capas Arquitectónicas</span>
            </div>
            <div className="flex flex-wrap gap-1">
              {layerBadges.map((badge) => (
                <button
                  key={badge.id}
                  onClick={() => setLayerFilter(badge.id)}
                  className={`px-2 py-1 rounded text-xs transition-colors cursor-pointer flex items-center space-x-1 ${
                    layerFilter === badge.id
                      ? 'bg-cyan-600 text-white font-medium'
                      : 'bg-slate-800/80 text-slate-400 hover:text-slate-200 hover:bg-slate-800'
                  }`}
                >
                  <span>{badge.label}</span>
                  <span className="text-[10px] opacity-70">({badge.count})</span>
                </button>
              ))}
            </div>
          </div>

          {/* File List */}
          <div className="flex-1 overflow-y-auto p-2 space-y-1">
            <div className="px-2 py-1 text-[11px] font-semibold text-slate-400 uppercase tracking-wider flex items-center justify-between">
              <span>Archivos ({filteredFiles.length})</span>
              <span className="text-slate-400 font-mono">Go 1.22</span>
            </div>

            {filteredFiles.map((file) => {
              const isSelected = selectedFilePath === file.path;
              return (
                <button
                  key={file.path}
                  onClick={() => {
                    setSelectedFilePath(file.path);
                    setViewMode('code');
                  }}
                  className={`w-full text-left p-2 rounded-md text-xs transition-all flex items-start space-x-2.5 cursor-pointer ${
                    isSelected
                      ? 'bg-slate-800 text-cyan-300 font-medium border border-cyan-500/30'
                      : 'text-slate-300 hover:bg-slate-800/60 hover:text-white'
                  }`}
                >
                  <FileCode
                    className={`w-4 h-4 shrink-0 mt-0.5 ${
                      isSelected ? 'text-cyan-400' : 'text-slate-400'
                    }`}
                  />
                  <div className="overflow-hidden flex-1">
                    <div className="flex items-center justify-between">
                      <span className="truncate font-mono">{file.name}</span>
                      <span className="text-[10px] text-slate-400 px-1 rounded bg-slate-950 font-mono shrink-0 ml-1">
                        {file.layer}
                      </span>
                    </div>
                    <span className="text-[11px] text-slate-400 truncate block font-sans">
                      {file.path}
                    </span>
                  </div>
                </button>
              );
            })}

            {filteredFiles.length === 0 && (
              <div className="p-6 text-center text-xs text-slate-400">
                No se encontraron archivos que coincidan con la búsqueda.
              </div>
            )}
          </div>

          {/* Architecture Status Widget */}
          <div className="p-3 border-t border-slate-800 bg-slate-950/60 text-xs">
            <div className="flex items-center justify-between text-slate-400">
              <span className="flex items-center space-x-1.5">
                <FolderGit2 className="w-3.5 h-3.5 text-cyan-400" />
                <span>Estructura Go</span>
              </span>
              <span className="text-emerald-400 font-medium">Modular & Limpia</span>
            </div>
          </div>
        </aside>

        {/* Main Content Area */}
        <main className="flex-1 overflow-hidden p-4 bg-slate-950">
          {viewMode === 'code' && <CodeViewer file={currentFile} />}
          {viewMode === 'architecture' && (
            <ArchitectureView
              onSelectFile={(path) => {
                setSelectedFilePath(path);
                setViewMode('code');
              }}
            />
          )}
          {viewMode === 'audit' && <AuditView />}
        </main>
      </div>

      {/* Footer Info Bar */}
      <footer className="border-t border-slate-800 px-5 py-2.5 bg-slate-900/90 text-xs text-slate-400 flex flex-col sm:flex-row items-center justify-between gap-2">
        <div className="flex items-center space-x-4">
          <span className="flex items-center space-x-1.5">
        <Sparkles className="w-3.5 h-3.5 text-cyan-400" />
        <span>FCOS Kernel v2.2.0 • Listo</span>
      </span>
    </div>
    <div className="flex items-center space-x-4 font-mono text-[10px] text-slate-500">
      <span>Klik Technologies © {new Date().getFullYear()}</span>
    </div>
  </footer>
</div>
  );
}