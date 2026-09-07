import { Shield, FileText, CheckCircle2, Copy, Check } from 'lucide-react';
import { useState } from 'react';
import { forensicData } from '../auditData';

export default function ForensicTable() {
  const [copiedHash, setCopiedHash] = useState<string | null>(null);

  const copyHash = (hash: string) => {
    navigator.clipboard.writeText(hash);
    setCopiedHash(hash);
    setTimeout(() => setCopiedHash(null), 2000);
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
      <div className="p-4 sm:p-5 border-b border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-slate-900/90">
        <div>
          <div className="flex items-center gap-2">
            <Shield className="w-4 h-4 text-cyan-400" />
            <h3 className="text-sm font-semibold text-white">
              Evidencia Forense: Archivos Críticos de Red de Windows
            </h3>
          </div>
          <p className="text-xs text-slate-400 mt-0.5">
            Origen: <code className="font-mono text-cyan-300">informe_forense.json</code> | Máquina: <span className="text-slate-200">{forensicData.encabezado_auditoria.nombre_equipo}</span> ({forensicData.encabezado_auditoria.sistema_operativo})
          </p>
        </div>

        <div className="text-xs text-slate-400 font-mono bg-slate-950 px-2.5 py-1 rounded border border-slate-800">
          Ejecutado: {new Date(forensicData.encabezado_auditoria.fecha_ejecucion_utc).toLocaleString()}
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs">
          <thead className="bg-slate-950/60 text-slate-400 border-b border-slate-800 uppercase tracking-wider font-mono text-[10px]">
            <tr>
              <th className="px-4 py-3">Archivo</th>
              <th className="px-4 py-3">Tamaño</th>
              <th className="px-4 py-3">Hash SHA-256</th>
              <th className="px-4 py-3">Última Modificación</th>
              <th className="px-4 py-3">Acción</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/60 text-slate-300">
            {forensicData.evidencia_archivos.map((file, idx) => {
              const fileName = file.ruta.split('\\').pop() || file.ruta;
              return (
                <tr key={idx} className="hover:bg-slate-800/40 transition">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <FileText className="w-3.5 h-3.5 text-slate-400 shrink-0" />
                      <div>
                        <span className="font-semibold text-white font-mono">{fileName}</span>
                        <div className="text-[11px] text-slate-500 font-mono truncate max-w-xs sm:max-w-sm">
                          {file.ruta}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-3 font-mono text-slate-300">
                    {file.tamaño_bytes.toLocaleString()} B
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-1.5 font-mono text-[11px] text-slate-400">
                      <span className="truncate max-w-[140px] sm:max-w-[180px]" title={file.sha256}>
                        {file.sha256.slice(0, 16)}...{file.sha256.slice(-8)}
                      </span>
                    </div>
                  </td>
                  <td className="px-4 py-3 font-mono text-slate-400 text-[11px]">
                    {file.ultima_modificacion_utc ? new Date(file.ultima_modificacion_utc).toLocaleDateString() : 'N/A'}
                  </td>
                  <td className="px-4 py-3">
                    <button
                      onClick={() => copyHash(file.sha256)}
                      className="flex items-center gap-1 px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition text-[11px]"
                    >
                      {copiedHash === file.sha256 ? (
                        <>
                          <Check className="w-3 h-3 text-emerald-400" />
                          <span className="text-emerald-400">Copiado</span>
                        </>
                      ) : (
                        <>
                          <Copy className="w-3 h-3" />
                          <span>Copiar SHA</span>
                        </>
                      )}
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
