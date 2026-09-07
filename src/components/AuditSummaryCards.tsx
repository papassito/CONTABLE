import { CheckCircle2, AlertTriangle, FileCode2, Shield, Flame, Terminal } from 'lucide-react';
import { auditItems } from '../auditData';

export default function AuditSummaryCards({ onSelectCategory }: { onSelectCategory: (cat: string) => void }) {
  const saneados = auditItems.filter(i => i.status === 'saneado').length;
  const pendientes = auditItems.filter(i => i.status === 'pendiente_code_assist').length;
  const verificados = auditItems.filter(i => i.status === 'verificado').length;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {/* Sanation Card */}
      <div 
        onClick={() => onSelectCategory('workspace')}
        className="p-4 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition cursor-pointer group"
      >
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs font-medium text-slate-400">Saneamiento de Archivos</span>
          <div className="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 group-hover:bg-emerald-500/20 transition">
            <CheckCircle2 className="w-4 h-4" />
          </div>
        </div>
        <div className="text-2xl font-bold text-white mb-1">
          8 Duplicados
        </div>
        <p className="text-xs text-emerald-400 flex items-center gap-1">
          <span>✓ Eliminados y árbol restaurado</span>
        </p>
      </div>

      {/* Backend Go Card */}
      <div 
        onClick={() => onSelectCategory('go-backend')}
        className="p-4 rounded-xl bg-slate-900 border border-amber-900/40 hover:border-amber-700/60 transition cursor-pointer group"
      >
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs font-medium text-slate-400">Módulo Backend Go</span>
          <div className="p-1.5 rounded-lg bg-amber-500/10 text-amber-400 group-hover:bg-amber-500/20 transition">
            <AlertTriangle className="w-4 h-4" />
          </div>
        </div>
        <div className="text-2xl font-bold text-amber-300 mb-1">
          Falta go.mod
        </div>
        <p className="text-xs text-amber-400/90">
          Requiere código Go para ejecutar tests
        </p>
      </div>

      {/* Forensics Card */}
      <div 
        onClick={() => onSelectCategory('forensic')}
        className="p-4 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition cursor-pointer group"
      >
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs font-medium text-slate-400">Auditoría Forense</span>
          <div className="p-1.5 rounded-lg bg-cyan-500/10 text-cyan-400 group-hover:bg-cyan-500/20 transition">
            <Shield className="w-4 h-4" />
          </div>
        </div>
        <div className="text-2xl font-bold text-white mb-1">
          6 Archivos Core
        </div>
        <p className="text-xs text-cyan-400">
          Hashes SHA-256 en informe_forense.json
        </p>
      </div>

      {/* CI/CD & Wails Card */}
      <div 
        onClick={() => onSelectCategory('cicd')}
        className="p-4 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition cursor-pointer group"
      >
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs font-medium text-slate-400">Pipeline & Wails</span>
          <div className="p-1.5 rounded-lg bg-indigo-500/10 text-indigo-400 group-hover:bg-indigo-500/20 transition">
            <Terminal className="w-4 h-4" />
          </div>
        </div>
        <div className="text-2xl font-bold text-white mb-1">
          Wails v2 + NSIS
        </div>
        <p className="text-xs text-indigo-400">
          Revisión de signtool.exe y bindings
        </p>
      </div>
    </div>
  );
}
