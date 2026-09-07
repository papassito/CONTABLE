import { CheckSquare, Square, AlertCircle, ArrowRight, Copy, Check } from 'lucide-react';
import { useState } from 'react';

interface Task {
  id: string;
  title: string;
  command?: string;
  description: string;
  priority: 'ALTA' | 'MEDIA' | 'INFORMACIÓN';
}

const tasks: Task[] = [
  {
    id: '1',
    title: '1. Integrar el Módulo Go (go.mod)',
    command: 'git clone https://github.com/papassito/CONTABLE.git temp_contable && cp -r temp_contable/go-skeleton ./go-skeleton && rm -rf temp_contable',
    description: 'Los scripts reparar_y_probar.ps1 y run-tests.ps1 requieren la presencia de go.mod en la raíz o en go-skeleton.',
    priority: 'ALTA'
  },
  {
    id: '2',
    title: '2. Ejecutar la Reparación de Contratos de Dominio',
    command: 'powershell -ExecutionPolicy Bypass -File .\\reparar_y_probar.ps1',
    description: 'Reescribe internal/domain/document_parser.go, reubica anchor.go y purga config/uow.go.',
    priority: 'ALTA'
  },
  {
    id: '3',
    title: '3. Dinamizar la Ruta de SignTool en CI/CD',
    command: 'release.yml: Reemplazar la ruta fija 10.0.22621.0 por resolución dinámica en el PATH de Windows Kits',
    description: 'Evita fallos en GitHub Actions cuando la imagen windows-latest renueva el SDK.',
    priority: 'MEDIA'
  },
  {
    id: '4',
    title: '4. Ejecutar Scanner-Agresivo en Máquina de Desarrollo',
    command: 'powershell -ExecutionPolicy Bypass -File .\\Scanner-Agresivo.ps1',
    description: 'Modo solo lectura: valida Windows Defender, conexiones TCP activas y descartar malware en Temp/AppData.',
    priority: 'INFORMACIÓN'
  }
];

export default function ActionChecklist() {
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const copyCommand = (cmd: string, id: string) => {
    navigator.clipboard.writeText(cmd);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <CheckSquare className="w-5 h-5 text-indigo-400" />
          <h3 className="text-sm font-semibold text-white">
            Ruta de Acción Recomendada para Code Assist
          </h3>
        </div>
        <span className="text-[11px] font-mono text-slate-400 bg-slate-950 px-2.5 py-1 rounded border border-slate-800">
          4 Pasos Clave
        </span>
      </div>

      <div className="space-y-3">
        {tasks.map(task => (
          <div 
            key={task.id}
            className="p-3.5 rounded-lg bg-slate-950/60 border border-slate-800/80 hover:border-slate-700/60 transition space-y-2"
          >
            <div className="flex items-center justify-between gap-2">
              <span className="text-xs font-semibold text-white flex items-center gap-2">
                {task.title}
              </span>
              <span className={`text-[10px] font-mono px-2 py-0.5 rounded font-medium ${
                task.priority === 'ALTA' 
                  ? 'bg-rose-500/10 text-rose-400 border border-rose-500/20' 
                  : task.priority === 'MEDIA'
                  ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
                  : 'bg-slate-800 text-slate-400'
              }`}>
                {task.priority}
              </span>
            </div>

            <p className="text-xs text-slate-400">
              {task.description}
            </p>

            {task.command && (
              <div className="flex items-center justify-between gap-2 p-2 rounded bg-slate-900 border border-slate-800 font-mono text-[11px] text-slate-300">
                <span className="truncate">{task.command}</span>
                <button
                  onClick={() => copyCommand(task.command!, task.id)}
                  className="shrink-0 p-1 rounded hover:bg-slate-800 text-slate-400 hover:text-white transition"
                  title="Copiar comando"
                >
                  {copiedId === task.id ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
                </button>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
