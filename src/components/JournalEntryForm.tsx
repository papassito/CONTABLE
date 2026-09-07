import React, { useState } from 'react';
import { createJournalDraft, toCents, JournalLine } from '../services/wailsService';

export const JournalEntryForm: React.FC = () => {
  const [concept, setConcept] = useState('');
  const [reference, setReference] = useState('');
  const [lines, setLines] = useState<JournalLine[]>([
    { account_id: '1', account_code: '110505', description: 'Caja General', debit: 0, credit: 0 },
    { account_id: '2', account_code: '111005', description: 'Bancos Nacionales', debit: 0, credit: 0 },
  ]);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  // Manejar cambio en los montos de las líneas (se ingresan como decimales y se convierte a céntimos)
  const handleLineChange = (index: number, field: 'debit' | 'credit', value: number) => {
    const updated = [...lines];
    updated[index][field] = toCents(value);
    setLines(updated);
  };

  // Calcular totales en tiempo real
  const totalDebitCents = lines.reduce((sum: number, l: JournalLine) => sum + l.debit, 0);
  const totalCreditCents = lines.reduce((sum: number, l: JournalLine) => sum + l.credit, 0);
  const isBalanced = totalDebitCents === totalCreditCents && totalDebitCents > 0;

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setErrorMessage(null);
    setSuccessMessage(null);

    if (!isBalanced) {
      setErrorMessage('Partida doble no cuadrada. Verifique que Débito y Crédito sean iguales.');
      return;
    }

    try {
      const createdEntry = await createJournalDraft({
        number: 'AST-' + Date.now(),
        date: new Date().toISOString(),
        concept,
        reference,
        lines,
      });

      setSuccessMessage(`Borrador guardado exitosamente con ID: ${createdEntry.id || 'Generado'}`);
      setConcept('');
      setReference('');
    } catch (err: any) {
      setErrorMessage(err.message);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-6 bg-slate-900 text-white rounded-xl border border-slate-800 space-y-4 max-w-2xl mx-auto">
      <h2 className="text-lg font-bold border-b border-slate-800 pb-2">Captura de Asiento Contable</h2>

      {errorMessage && (
        <div className="p-3 bg-rose-950/80 border border-rose-800 text-rose-300 rounded text-xs">
          ❌ {errorMessage}
        </div>
      )}

      {successMessage && (
        <div className="p-3 bg-emerald-950/80 border border-emerald-800 text-emerald-300 rounded text-xs">
          ✅ {successMessage}
        </div>
      )}

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-xs font-medium text-slate-400 mb-1">Concepto / Glosa</label>
          <input
            type="text"
            value={concept}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setConcept(e.target.value)}
            required
            className="w-full bg-slate-950 border border-slate-700 rounded px-3 py-1.5 text-sm text-white"
            placeholder="Ej: Pago de servicios públicos"
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-slate-400 mb-1">Referencia</label>
          <input
            type="text"
            value={reference}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setReference(e.target.value)}
            className="w-full bg-slate-950 border border-slate-700 rounded px-3 py-1.5 text-sm text-white"
            placeholder="Ej: Factura 1234"
          />
        </div>
      </div>

      <div className="space-y-2 pt-2">
        <h3 className="text-xs font-bold text-slate-400 uppercase">Partidas / Movimientos</h3>
        {lines.map((line: JournalLine, idx: number) => (
          <div key={idx} className="flex items-center space-x-2 bg-slate-950 p-2 rounded border border-slate-800">
            <span className="text-xs font-mono text-cyan-400 w-16">{line.account_code}</span>
            <span className="text-xs text-slate-300 flex-1">{line.description}</span>
            <input
              type="number"
              step="0.01"
              placeholder="Débito"
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleLineChange(idx, 'debit', parseFloat(e.target.value) || 0)}
              className="w-24 bg-slate-900 border border-slate-700 rounded px-2 py-1 text-xs text-right text-emerald-400"
            />
            <input
              type="number"
              step="0.01"
              placeholder="Crédito"
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleLineChange(idx, 'credit', parseFloat(e.target.value) || 0)}
              className="w-24 bg-slate-900 border border-slate-700 rounded px-2 py-1 text-xs text-right text-amber-400"
            />
          </div>
        ))}
      </div>

      <div className="flex justify-between items-center pt-2 border-t border-slate-800 text-xs font-mono">
        <div>
          Débitos: <span className="text-emerald-400">${(totalDebitCents / 100).toFixed(2)}</span> | 
          Créditos: <span className="text-amber-400">${(totalCreditCents / 100).toFixed(2)}</span>
        </div>
        <button
          type="submit"
          disabled={!isBalanced}
          className={`px-4 py-2 rounded font-bold transition-all ${
            isBalanced 
              ? 'bg-cyan-600 hover:bg-cyan-500 text-white cursor-pointer' 
              : 'bg-slate-800 text-slate-500 cursor-not-allowed'
          }`}
        >
          Guardar Borrador
        </button>
      </div>
    </form>
  );
};