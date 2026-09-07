import React, { useState, useEffect } from 'react';
import { 
  JournalEntry, 
  JournalLine,
  getJournalEntry, 
  postJournalEntry, 
  reverseJournalEntry, 
  fromCents 
} from '../services/wailsService';

interface Props {
  entryId: string;
  onStatusChange?: () => void;
}

export const JournalEntryDetail: React.FC<Props> = ({ entryId, onStatusChange }) => {
  const [entry, setEntry] = useState<JournalEntry | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [reversionReason, setReversionReason] = useState<string>('');
  const [showReversionModal, setShowReversionModal] = useState<boolean>(false);

  const fetchEntry = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getJournalEntry(entryId);
      setEntry(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (entryId) {
      fetchEntry();
    }
  }, [entryId]);

  const handlePost = async () => {
    if (!entry?.id) return;
    setError(null);
    try {
      await postJournalEntry(entry.id);
      await fetchEntry();
      if (onStatusChange) onStatusChange();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleReverse = async () => {
    if (!entry?.id || !reversionReason.trim()) return;
    setError(null);
    try {
      await reverseJournalEntry(entry.id, reversionReason);
      setShowReversionModal(false);
      setReversionReason('');
      await fetchEntry();
      if (onStatusChange) onStatusChange();
    } catch (err: any) {
      setError(err.message);
    }
  };

  if (loading) {
    return <div className="p-6 text-slate-400 text-sm">Cargando asiento contable...</div>;
  }

  if (error && !entry) {
    return (
      <div className="p-4 bg-rose-950/80 border border-rose-800 text-rose-300 rounded-lg text-xs">
        ⚠️ {error}
      </div>
    );
  }

  if (!entry) return null;

  const totalDebit = entry.lines.reduce((acc: number, l: JournalLine) => acc + l.debit, 0);
  const totalCredit = entry.lines.reduce((acc: number, l: JournalLine) => acc + l.credit, 0);

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 text-slate-200 space-y-6 max-w-4xl mx-auto shadow-lg">
      {/* Alerta de error superior */}
      {error && (
        <div className="p-3 bg-rose-950/80 border border-rose-800 text-rose-300 rounded text-xs">
          ⚠️ {error}
        </div>
      )}

      {/* Encabezado del Asiento */}
      <div className="flex justify-between items-start border-b border-slate-800 pb-4">
        <div>
          <div className="flex items-center space-x-3">
            <h2 className="text-xl font-bold text-white">Asiento: {entry.number}</h2>
            <span className={`text-[10px] px-2.5 py-0.5 rounded-full font-mono font-bold ${
              entry.status === 'CONTABILIZADO' 
                ? 'bg-emerald-950 text-emerald-300 border border-emerald-800/50' 
                : entry.status === 'ANULADO'
                ? 'bg-rose-950 text-rose-300 border border-rose-800/50'
                : 'bg-amber-950 text-amber-300 border border-amber-800/50'
            }`}>
              {entry.status || 'BORRADOR'}
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">{entry.concept}</p>
          <span className="text-[11px] font-mono text-cyan-400/90 mt-0.5 block">
            Ref: {entry.reference || 'Sin referencia'}
          </span>
        </div>

        <div className="text-right text-xs text-slate-400 font-mono">
          <div>Fecha: {new Date(entry.date).toLocaleDateString()}</div>
          <div className="text-[10px] text-slate-500">ID: {entry.id}</div>
        </div>
      </div>

      {/* Tabla de Partidas */}
      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs border-collapse">
          <thead>
            <tr className="border-b border-slate-800 text-slate-400 font-mono uppercase text-[10px]">
              <th className="py-2 px-3">Cuenta</th>
              <th className="py-2 px-3">Descripción</th>
              <th className="py-2 px-3 text-right">Débito</th>
              <th className="py-2 px-3 text-right">Crédito</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/50">
              {entry.lines.map((line: JournalLine, idx: number) => (
              <tr key={idx} className="hover:bg-slate-800/30">
                <td className="py-2.5 px-3 font-mono text-cyan-400">{line.account_code}</td>
                <td className="py-2.5 px-3 text-slate-300">{line.description}</td>
                <td className="py-2.5 px-3 text-right font-mono text-emerald-400">
                  {line.debit > 0 ? `$${fromCents(line.debit).toFixed(2)}` : '-'}
                </td>
                <td className="py-2.5 px-3 text-right font-mono text-amber-400">
                  {line.credit > 0 ? `$${fromCents(line.credit).toFixed(2)}` : '-'}
                </td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr className="border-t border-slate-700 font-bold font-mono text-xs bg-slate-950/50">
              <td colSpan={2} className="py-3 px-3 text-right uppercase text-slate-400">Totales:</td>
              <td className="py-3 px-3 text-right text-emerald-400">${fromCents(totalDebit).toFixed(2)}</td>
              <td className="py-3 px-3 text-right text-amber-400">${fromCents(totalCredit).toFixed(2)}</td>
            </tr>
          </tfoot>
        </table>
      </div>

      {/* Acciones del Asiento */}
      <div className="flex justify-end space-x-3 pt-2 border-t border-slate-800">
        {entry.status === 'BORRADOR' && (
          <button
            onClick={handlePost}
            className="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs rounded transition-all shadow-sm"
          >
            Contabilizar Definitivo
          </button>
        )}

        {entry.status === 'CONTABILIZADO' && (
          <button
            onClick={() => setShowReversionModal(true)}
            className="px-4 py-2 bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs rounded transition-all shadow-sm"
          >
            Revertir / Anular Asiento
          </button>
        )}
      </div>

      {/* Modal de Reversión */}
      {showReversionModal && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center p-4 z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 max-w-md w-full space-y-4 shadow-2xl">
            <h3 className="text-sm font-bold text-white">Motivo de Anulación Contable</h3>
            <p className="text-xs text-slate-400">
              Esta acción creará un contraasiento de reversión inmutable en el Libro Mayor.
            </p>
            <textarea
              value={reversionReason}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setReversionReason(e.target.value)}
              placeholder="Escriba el motivo de la anulación (ej. Error en imputación de cuenta)..."
              className="w-full bg-slate-950 border border-slate-700 rounded p-2 text-xs text-white h-24 focus:outline-none focus:border-cyan-500"
            />
            <div className="flex justify-end space-x-2">
              <button
                onClick={() => setShowReversionModal(false)}
                className="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs rounded"
              >
                Cancelar
              </button>
              <button
                onClick={handleReverse}
                disabled={!reversionReason.trim()}
                className="px-3 py-1.5 bg-rose-600 hover:bg-rose-500 disabled:opacity-50 text-white text-xs font-bold rounded"
              >
                Confirmar Reversión
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};