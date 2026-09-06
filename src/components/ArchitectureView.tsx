import React from 'react';
import { Layers, Database, ShieldCheck, Cpu, Globe, FolderGit2, BookOpen, CheckCircle2 } from 'lucide-react';

interface ArchitectureViewProps {
  onSelectFile: (path: string) => void;
}

export const ArchitectureView: React.FC<ArchitectureViewProps> = ({ onSelectFile }) => {
  return (
    <div id="architecture-view-container" className="p-6 space-y-6 max-w-5xl mx-auto overflow-y-auto h-full">
      {/* Overview Banner */}
      <div className="bg-gradient-to-r from-slate-900 to-slate-950 p-6 rounded-xl border border-slate-800 shadow-sm">
        <div className="flex items-start justify-between">
          <div className="space-y-2">
            <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-cyan-950/80 border border-cyan-800/40 text-cyan-400 text-xs font-semibold">
              <FolderGit2 className="w-3.5 h-3.5" />
              <span>Arquitectura Modular en Go</span>
            </div>
            <h2 className="text-xl font-bold text-white">Contable Fix by KLIK</h2>
            <p className="text-sm text-slate-400 max-w-2xl leading-relaxed">
              Esqueleto limpio, desacoplado y vacío listo para ser llenado con las reglas contables de la empresa.
              Sigue los principios de <strong className="text-slate-200">Clean Architecture / Puertos y Adaptadores</strong> para que el motor contable no dependa de ningún framework o base de datos específica.
            </p>
          </div>
          <div className="hidden md:flex flex-col items-end space-y-1 text-xs text-slate-400 font-mono">
            <span className="bg-slate-800/80 px-2.5 py-1 rounded border border-slate-700/60 text-slate-300">Package: github.com/klik/contable-fix</span>
            <span className="text-cyan-400">Go 1.22+ Standard Library</span>
          </div>
        </div>
      </div>

      {/* Clean Architecture Diagram Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {/* Layer 1: HTTP Handlers */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-lg p-4 space-y-3">
          <div className="flex items-center space-x-2 text-cyan-400">
            <Globe className="w-4 h-4" />
            <h3 className="text-xs font-bold uppercase tracking-wider">1. Entrega HTTP (API)</h3>
          </div>
          <p className="text-xs text-slate-400">
            Recibe peticiones HTTP, decodifica JSON y serializa respuestas.
          </p>
          <div className="space-y-1.5 pt-2">
            <button
              onClick={() => onSelectFile('internal/handler/http/router.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-cyan-300 font-mono transition-colors block truncate cursor-pointer"
            >
              router.go
            </button>
            <button
              onClick={() => onSelectFile('internal/handler/http/account_handler.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-cyan-300 font-mono transition-colors block truncate cursor-pointer"
            >
              account_handler.go
            </button>
            <button
              onClick={() => onSelectFile('internal/handler/http/journal_handler.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-cyan-300 font-mono transition-colors block truncate cursor-pointer"
            >
              journal_handler.go
            </button>
          </div>
        </div>

        {/* Layer 2: Services / Use Cases */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-lg p-4 space-y-3">
          <div className="flex items-center space-x-2 text-purple-400">
            <Cpu className="w-4 h-4" />
            <h3 className="text-xs font-bold uppercase tracking-wider">2. Casos de Uso (Service)</h3>
          </div>
          <p className="text-xs text-slate-400">
            Orquesta transacciones contables, partida doble y validación de balances.
          </p>
          <div className="space-y-1.5 pt-2">
            <button
              onClick={() => onSelectFile('internal/service/journal_service.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-purple-300 font-mono transition-colors block truncate cursor-pointer"
            >
              journal_service.go
            </button>
            <button
              onClick={() => onSelectFile('internal/service/account_service.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-purple-300 font-mono transition-colors block truncate cursor-pointer"
            >
              account_service.go
            </button>
            <button
              onClick={() => onSelectFile('internal/service/ledger_service.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-purple-300 font-mono transition-colors block truncate cursor-pointer"
            >
              ledger_service.go
            </button>
          </div>
        </div>

        {/* Layer 3: Domain Entities */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-lg p-4 space-y-3">
          <div className="flex items-center space-x-2 text-emerald-400">
            <Layers className="w-4 h-4" />
            <h3 className="text-xs font-bold uppercase tracking-wider">3. Dominio Contable</h3>
          </div>
          <p className="text-xs text-slate-400">
            Modelos de cuentas, asientos, partidas dobles, libro mayor e impuestos.
          </p>
          <div className="space-y-1.5 pt-2">
            <button
              onClick={() => onSelectFile('internal/domain/account.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-emerald-300 font-mono transition-colors block truncate cursor-pointer"
            >
              account.go
            </button>
            <button
              onClick={() => onSelectFile('internal/domain/journal.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-emerald-300 font-mono transition-colors block truncate cursor-pointer"
            >
              journal.go
            </button>
            <button
              onClick={() => onSelectFile('internal/domain/ledger.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-emerald-300 font-mono transition-colors block truncate cursor-pointer"
            >
              ledger.go
            </button>
          </div>
        </div>

        {/* Layer 4: Repositories */}
        <div className="bg-slate-900/90 border border-slate-800 rounded-lg p-4 space-y-3">
          <div className="flex items-center space-x-2 text-amber-400">
            <Database className="w-4 h-4" />
            <h3 className="text-xs font-bold uppercase tracking-wider">4. Puertos / Persistencia</h3>
          </div>
          <p className="text-xs text-slate-400">
            Interfaces puras para desacoplar PostgreSQL, MySQL, SQLite o memoria.
          </p>
          <div className="space-y-1.5 pt-2">
            <button
              onClick={() => onSelectFile('internal/repository/account_repository.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-amber-300 font-mono transition-colors block truncate cursor-pointer"
            >
              account_repository.go
            </button>
            <button
              onClick={() => onSelectFile('internal/repository/journal_repository.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-amber-300 font-mono transition-colors block truncate cursor-pointer"
            >
              journal_repository.go
            </button>
            <button
              onClick={() => onSelectFile('internal/repository/ledger_repository.go')}
              className="w-full text-left text-xs p-2 rounded bg-slate-800/60 hover:bg-slate-800 text-slate-200 hover:text-amber-300 font-mono transition-colors block truncate cursor-pointer"
            >
              ledger_repository.go
            </button>
          </div>
        </div>
      </div>

      {/* Design Principles Checklist */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 space-y-4">
        <div className="flex items-center space-x-2 text-slate-200 font-semibold text-sm">
          <ShieldCheck className="w-4 h-4 text-cyan-400" />
          <span>Características del Esqueleto Go Diseñado para KLIK</span>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
          <div className="flex items-start space-x-2.5 p-3 rounded bg-slate-950/60 border border-slate-800/60">
            <CheckCircle2 className="w-4 h-4 text-cyan-400 shrink-0 mt-0.5" />
            <div>
              <strong className="text-slate-200 block mb-1">Estructura Vacía (Clean Stubs)</strong>
              <span className="text-slate-400">Funciones y métodos con firmas y contratos listos, con marcas <code>// TODO:</code> para ser llenados.</span>
            </div>
          </div>

          <div className="flex items-start space-x-2.5 p-3 rounded bg-slate-950/60 border border-slate-800/60">
            <CheckCircle2 className="w-4 h-4 text-purple-400 shrink-0 mt-0.5" />
            <div>
              <strong className="text-slate-200 block mb-1">Partida Doble & Dominio Contable</strong>
              <span className="text-slate-400">Entidades preparadas para activo, pasivo, patrimonio, balance de comprobación y cuadre matemático.</span>
            </div>
          </div>

          <div className="flex items-start space-x-2.5 p-3 rounded bg-slate-950/60 border border-slate-800/60">
            <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
            <div>
              <strong className="text-slate-200 block mb-1">Inyección de Dependencias</strong>
              <span className="text-slate-400">Configurado en <code>cmd/api/main.go</code> sin librerías invasivas, usando Go idiomático.</span>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Start Guide */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 space-y-3">
        <div className="flex items-center space-x-2 text-slate-200 font-semibold text-sm">
          <BookOpen className="w-4 h-4 text-cyan-400" />
          <span>Cómo comenzar en tu entorno local con Go</span>
        </div>
        <div className="p-3 bg-slate-950 rounded-lg border border-slate-800/80 font-mono text-xs text-slate-300 space-y-2">
          <div className="text-slate-400"># 1. Descomprimir el esqueleto o clonar la carpeta</div>
          <div className="text-cyan-300">cd contable-fix-klik</div>
          <div className="text-slate-400"># 2. Descargar o verificar dependencias del módulo</div>
          <div className="text-cyan-300">go mod tidy</div>
          <div className="text-slate-400"># 3. Compilar o ejecutar directamente el servidor</div>
          <div className="text-cyan-300">go run cmd/api/main.go</div>
        </div>
      </div>
    </div>
  );
};
