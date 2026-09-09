import { ShieldAlert, CheckCircle2, Terminal, Copy, Check } from 'lucide-react';
import { useState } from 'react';
import { generateCodeAssistReportMarkdown } from '../auditData';

export default function Header() {
  const [copied, setCopied] = useState(false);

  const handleCopyReport = () => {
    navigator.clipboard.writeText(generateCodeAssistReportMarkdown());
    setCopied(true);
    setTimeout(() => setCopied(false), 2500);
  };

  return (
    <header className="border-b border-slate-800 bg-slate-900/80 backdrop-blur sticky top-0 z-20">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-indigo-500/20 to-cyan-500/20 border border-indigo-500/30 flex items-center justify-center text-indigo-400 shadow-inner">
            <ShieldAlert className="w-5 h-5" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-base sm:text-lg font-semibold text-white tracking-tight">
                Contable Fix FCOS v2.2
              </h1>
              <span className="px-2 py-0.5 text-[11px] font-mono rounded bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
                Auditoría Milimétrica
              </span>
            </div>
            <p className="text-xs text-slate-400 hidden sm:block">
              Informe Técnico & Forense para Code Assist y Desarrollo
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2 sm:gap-3">
          <div className="flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-medium">
            <CheckCircle2 className="w-3.5 h-3.5" />
            <span className="hidden sm:inline">Workspace Saneado</span>
            <span className="sm:hidden">Limpio</span>
          </div>

          <button
            id="copy-report-btn"
            onClick={handleCopyReport}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-500 active:bg-indigo-700 text-white text-xs font-medium transition shadow-sm"
          >
            {copied ? <Check className="w-3.5 h-3.5 text-emerald-300" /> : <Copy className="w-3.5 h-3.5" />}
            <span>{copied ? '¡Copiado!' : 'Copiar para Code Assist'}</span>
          </button>
        </div>
      </div>
    </header>
  );
}
