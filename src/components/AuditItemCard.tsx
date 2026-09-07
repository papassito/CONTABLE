import React, { useState } from 'react';
import { AuditItem } from '../types';
import { 
  CheckCircle2, 
  AlertTriangle, 
  AlertCircle, 
  Info, 
  ChevronDown, 
  ChevronUp, 
  FileCode, 
  Terminal, 
  Lightbulb, 
  Copy, 
  Check 
} from 'lucide-react';

export default function AuditItemCard({ item }: { item: AuditItem; key?: React.Key }) {
  const [expanded, setExpanded] = useState(true);
  const [copied, setCopied] = useState(false);

  const getSeverityBadge = () => {
    switch (item.severity) {
      case 'critical':
        return (
          <span className="flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-medium bg-rose-500/10 text-rose-400 border border-rose-500/20">
            <AlertCircle className="w-3 h-3" />
            Crítico
          </span>
        );
      case 'warning':
        return (
          <span className="flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-medium bg-amber-500/10 text-amber-400 border border-amber-500/20">
            <AlertTriangle className="w-3 h-3" />
            Atención
          </span>
        );
      case 'success':
        return (
          <span className="flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-medium bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <CheckCircle2 className="w-3 h-3" />
            Saneado
          </span>
        );
      default:
        return (
          <span className="flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-medium bg-blue-500/10 text-blue-400 border border-blue-500/20">
            <Info className="w-3 h-3" />
            Auditado
          </span>
        );
    }
  };

  const getStatusBadge = () => {
    switch (item.status) {
      case 'saneado':
        return <span className="text-[11px] font-medium text-emerald-400">Resuelto</span>;
      case 'pendiente_code_assist':
        return <span className="text-[11px] font-medium text-amber-400">Acción Requerida</span>;
      default:
        return <span className="text-[11px] font-medium text-slate-400">Normal</span>;
    }
  };

  const handleCopyAction = () => {
    if (!item.suggestedActionForCodeAssist) return;
    navigator.clipboard.writeText(item.suggestedActionForCodeAssist);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden hover:border-slate-700/80 transition">
      {/* Card Header */}
      <div 
        onClick={() => setExpanded(!expanded)}
        className="p-4 flex items-center justify-between gap-3 cursor-pointer select-none bg-slate-900/60"
      >
        <div className="flex items-center gap-3 min-w-0">
          {getSeverityBadge()}
          <h3 className="text-sm font-semibold text-white truncate">
            {item.title}
          </h3>
        </div>

        <div className="flex items-center gap-3 shrink-0">
          <div className="hidden sm:block">
            {getStatusBadge()}
          </div>
          <button className="text-slate-400 hover:text-slate-200">
            {expanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
          </button>
        </div>
      </div>

      {/* Card Body */}
      {expanded && (
        <div className="p-4 pt-1 border-t border-slate-800/60 space-y-3.5 text-xs text-slate-300">
          <p className="text-slate-300 leading-relaxed font-normal">
            {item.summary}
          </p>

          <div className="p-3 rounded-lg bg-slate-950/70 border border-slate-800 font-mono text-[11px] text-slate-400 space-y-1">
            <span className="text-slate-500 uppercase tracking-wider text-[10px] block font-sans font-semibold">
              Detalles de Auditoría:
            </span>
            <p className="text-slate-300 whitespace-pre-wrap leading-relaxed">
              {item.technicalDetails}
            </p>
          </div>

          {/* Target files */}
          {item.targetFiles.length > 0 && (
            <div className="flex flex-wrap items-center gap-1.5 pt-1">
              <span className="text-slate-500 text-[11px] font-mono">Archivos:</span>
              {item.targetFiles.map((f, i) => (
                <span 
                  key={i}
                  className="px-2 py-0.5 rounded bg-slate-800 text-slate-300 font-mono text-[11px] border border-slate-700/50"
                >
                  {f}
                </span>
              ))}
            </div>
          )}

          {/* Action for Code Assist */}
          {item.suggestedActionForCodeAssist && (
            <div className="p-3 rounded-lg bg-indigo-950/30 border border-indigo-500/20 text-indigo-200 flex items-start justify-between gap-3">
              <div className="space-y-1">
                <div className="flex items-center gap-1.5 text-indigo-400 font-semibold text-[11px]">
                  <Lightbulb className="w-3.5 h-3.5" />
                  <span>Sugerencia / Acción para Code Assist:</span>
                </div>
                <p className="text-indigo-100/90 text-xs">
                  {item.suggestedActionForCodeAssist}
                </p>
              </div>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  handleCopyAction();
                }}
                className="shrink-0 p-1.5 rounded bg-indigo-900/50 hover:bg-indigo-800 text-indigo-300 hover:text-white transition"
                title="Copiar acción"
              >
                {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
