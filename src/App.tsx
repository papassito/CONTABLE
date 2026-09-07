import { useState } from 'react';
import Header from './components/Header';
import AuditSummaryCards from './components/AuditSummaryCards';
import AuditItemCard from './components/AuditItemCard';
import ForensicTable from './components/ForensicTable';
import ScriptInspector from './components/ScriptInspector';
import ActionChecklist from './components/ActionChecklist';
import { auditItems } from './auditData';
import { ShieldCheck, Filter, AlertCircle, FileSpreadsheet, Layers, Sparkles } from 'lucide-react';

export default function App() {
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [selectedSeverity, setSelectedSeverity] = useState<string>('all');

  const filteredItems = auditItems.filter(item => {
    if (selectedCategory !== 'all' && item.category !== selectedCategory) return false;
    if (selectedSeverity !== 'all' && item.severity !== selectedSeverity) return false;
    return true;
  });

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col selection:bg-indigo-500 selection:text-white font-sans antialiased">
      <Header />

      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
        
        {/* Top Notification Banner: Sanation Complete */}
        <div className="p-4 rounded-xl bg-gradient-to-r from-emerald-950/40 via-slate-900 to-indigo-950/30 border border-emerald-500/20 flex flex-col sm:flex-row sm:items-center justify-between gap-4 shadow-sm">
          <div className="flex items-start sm:items-center gap-3">
            <div className="h-9 w-9 rounded-lg bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 shrink-0 mt-0.5 sm:mt-0">
              <ShieldCheck className="w-5 h-5" />
            </div>
            <div>
              <h2 className="text-sm font-semibold text-white">
                Saneamiento de Duplicados Ejecutado con Éxito
              </h2>
              <p className="text-xs text-slate-400 mt-0.5">
                Se detectaron y eliminaron 8 archivos redundantes con prefijos y sufijos erróneos. El espacio de trabajo está listo para la auditoría y correcciones de Code Assist.
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2 self-end sm:self-center shrink-0 font-mono text-[11px] bg-slate-950/80 px-3 py-1.5 rounded-lg border border-slate-800 text-emerald-400">
            <span>Árbol Canónico OK</span>
          </div>
        </div>

        {/* Summary Metrics */}
        <AuditSummaryCards onSelectCategory={(cat) => setSelectedCategory(cat)} />

        {/* Action Checklist & Forensic Evidence Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
          <div className="lg:col-span-7 space-y-6">
            <ForensicTable />
            <ScriptInspector />
          </div>

          <div className="lg:col-span-5 space-y-6">
            <ActionChecklist />

            {/* Audit Filter & List */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Filter className="w-4 h-4 text-slate-400" />
                  <h3 className="text-sm font-semibold text-white">
                    Hallazgos Milimétricos ({filteredItems.length})
                  </h3>
                </div>

                <div className="flex items-center gap-2">
                  <select
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                    className="bg-slate-900 border border-slate-800 text-xs rounded-lg px-2.5 py-1 text-slate-300 focus:outline-none focus:border-indigo-500 font-mono"
                  >
                    <option value="all">Todas las Categorías</option>
                    <option value="workspace">Workspace</option>
                    <option value="go-backend">Backend Go</option>
                    <option value="forensic">Forense</option>
                    <option value="wails">Wails</option>
                    <option value="cicd">CI/CD</option>
                    <option value="powershell">PowerShell</option>
                  </select>
                </div>
              </div>

              <div className="space-y-3">
                {filteredItems.map(item => (
                  <AuditItemCard key={item.id} item={item} />
                ))}
              </div>
            </div>
          </div>
        </div>

      </main>

      {/* Footer */}
      <footer className="border-t border-slate-800/80 bg-slate-950 py-6 text-center text-xs text-slate-500 font-mono">
        <p>Contable Fix by KLIK • FCOS v2.2 • Auditoría Técnica y Diagnóstico Forense</p>
      </footer>
    </div>
  );
}
