import React, { useState } from 'react';
import { 
  ShieldAlert, 
  CheckSquare, 
  Square, 
  ListChecks, 
  Clock, 
  Scale
} from 'lucide-react';

interface AuditCriterion {
  id: string;
  title: string;
  category: 'critical' | 'architecture' | 'database' | 'security';
  priority: 'ALTA' | 'CRÍTICA' | 'MEDIA';
  description: string;
  deliverable: string;
}

const AUDIT_CRITERIA: AuditCriterion[] = [
  {
    id: 'crit-1',
    category: 'critical',
    priority: 'CRÍTICA',
    title: 'Invariante de Partida Doble (Débito == Crédito)',
    description: 'Bajo ninguna circunstancia un asiento contable puede quedar desbalanceado. Validar que la suma total de débitos sea exactamente igual a la suma de créditos antes de permitir el estado CONTABILIZADO.',
    deliverable: 'Implementar validación en journal_service.go y test unitario en journal_service_test.go con múltiples líneas.',
  },
  {
    id: 'crit-2',
    category: 'critical',
    priority: 'CRÍTICA',
    title: 'Inmutabilidad de Asientos Contabilizados',
    description: 'Prohibir operaciones SQL UPDATE o DELETE sobre asientos en estado CONTABILIZADO. Las correcciones contables se realizan exclusivamente mediante asientos de reversión/anulación con contrapartida.',
    deliverable: 'Método ReverseEntry() en journal_service.go con registro del motivo y referencia al asiento original.',
  },
  {
    id: 'crit-3',
    category: 'architecture',
    priority: 'ALTA',
    title: 'Transaccionalidad Atómica (ACID)',
    description: 'La operación PostEntry() debe ejecutarse dentro de una transacción única de base de datos (tx): actualizar estado de asiento, insertar en libro mayor y actualizar saldos de cuentas de forma indivisible.',
    deliverable: 'Implementar BeginTx() en la capa de persistencia y pasar contexto transaccional en repository.',
  },
  {
    id: 'crit-4',
    category: 'database',
    priority: 'ALTA',
    title: 'Manejo de Moneda y Precisión Decimal',
    description: 'En contabilidad profesional no se debe usar float64 para saldos acumulados debido a errores de redondeo de punto flotante. Usar tipos decimales de precisión fija (ej: shopspring/decimal o entero en centavos int64).',
    deliverable: 'Ajustar campos monetarios en domain a tipos decimales y columnas NUMERIC(18, 4) en PostgreSQL.',
  },
  {
    id: 'crit-5',
    category: 'architecture',
    priority: 'ALTA',
    title: 'Regla de Cuentas Auxiliares (Plan Contable)',
    description: 'Solo las cuentas de último nivel jerárquico (auxiliares) que tengan accepts_move = true pueden recibir imputaciones contables en las líneas de los asientos.',
    deliverable: 'Validación en account_service.go y journal_service.go antes de crear o contabilizar líneas.',
  },
  {
    id: 'crit-6',
    category: 'security',
    priority: 'ALTA',
    title: 'Pista de Auditoría Automática (AuditLog)',
    description: 'Toda acción que altere el estado contable (crear cuenta, contabilizar asiento, anular factura) debe dejar un registro inmutable con usuario, IP, timestamp y detalle en audit.go.',
    deliverable: 'Middleware de auditoría en handler/http y registro en base de datos.',
  },
];

export const AuditView: React.FC = () => {
  const [completedItems, setCompletedItems] = useState<Record<string, boolean>>({});

  const toggleCheck = (id: string) => {
    setCompletedItems(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  const completedCount = Object.values(completedItems).filter(Boolean).length;
  const progressPercent = Math.round((completedCount / AUDIT_CRITERIA.length) * 100);

  return (
    <div id="audit-view-container" className="p-6 space-y-6 max-w-5xl mx-auto overflow-y-auto h-full text-slate-200">
      {/* Auditor Banner */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 shadow-sm">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="space-y-1.5">
            <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-cyan-950/80 border border-cyan-800/50 text-cyan-400 text-xs font-semibold">
              <ShieldAlert className="w-3.5 h-3.5" />
              <span>Rol: Auditor Técnico & Arquitecto Líder</span>
            </div>
            <h2 className="text-xl font-bold text-white">Directrices de Auditoría y Roadmap</h2>
            <p className="text-xs text-slate-400 max-w-2xl leading-relaxed">
              Instrucciones técnicas formales para el desarrollador que completará el motor de <strong>Contable Fix by KLIK</strong>.
              Cumplimiento estricto de principios contables GAAP/NIIF y estándares de código en Go.
            </p>
          </div>

          <div className="bg-slate-950 px-4 py-3 rounded-lg border border-slate-800 shrink-0">
            <div className="text-[11px] text-slate-400 uppercase tracking-wider mb-1">Checklist de Auditoría</div>
            <div className="flex items-center space-x-3">
              <span className="text-xl font-bold text-cyan-400">{completedCount} / {AUDIT_CRITERIA.length}</span>
              <span className="text-xs text-slate-400">({progressPercent}%)</span>
            </div>
            <div className="w-32 bg-slate-800 h-1.5 rounded-full mt-2 overflow-hidden">
              <div 
                className="bg-cyan-500 h-full transition-all duration-300"
                style={{ width: `${progressPercent}%` }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* 4 Phases Roadmap */}
      <div className="space-y-3">
        <h3 className="text-sm font-bold text-white uppercase tracking-wider flex items-center space-x-2">
          <Clock className="w-4 h-4 text-cyan-400" />
          <span>Fases de Ejecución Recomendadas para el Desarrollador</span>
        </h3>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-3">
          <div className="bg-slate-900 border border-slate-800 p-4 rounded-lg space-y-2">
            <div className="text-xs font-bold text-cyan-400">FASE 1: Persistencia</div>
            <div className="text-xs text-slate-300 font-medium">PostgreSQL & Repositorios</div>
            <p className="text-[11px] text-slate-400">
              Crear migraciones SQL (tablas <code>accounts</code>, <code>journal_entries</code>, <code>journal_lines</code>). Implementar interfaces con soporte para transacciones <code>*sql.Tx</code>.
            </p>
          </div>

          <div className="bg-slate-900 border border-slate-800 p-4 rounded-lg space-y-2">
            <div className="text-xs font-bold text-purple-400">FASE 2: Plan Contable</div>
            <div className="text-xs text-slate-300 font-medium">Jerarquía y Cuentas</div>
            <p className="text-[11px] text-slate-400">
              Validar unicidad de códigos, niveles jerárquicos (1 al 5) y prevenir imputaciones directas en cuentas de mayor (padres).
            </p>
          </div>

          <div className="bg-slate-900 border border-slate-800 p-4 rounded-lg space-y-2">
            <div className="text-xs font-bold text-emerald-400">FASE 3: Motor Contable</div>
            <div className="text-xs text-slate-300 font-medium">Asientos & Libro Mayor</div>
            <p className="text-[11px] text-slate-400">
              Validar partida doble, posteo atómico a libro mayor, cálculo de saldos y método de reversión automática.
            </p>
          </div>

          <div className="bg-slate-900 border border-slate-800 p-4 rounded-lg space-y-2">
            <div className="text-xs font-bold text-amber-400">FASE 4: Seguridad & Test</div>
            <div className="text-xs text-slate-300 font-medium">AuditLog & Pruebas</div>
            <p className="text-[11px] text-slate-400">
              Middleware de auditoría, manejo centralizado de errores contables y suite de tests unitarios con casos de borde.
            </p>
          </div>
        </div>
      </div>

      {/* Audit Criteria Checklist */}
      <div className="space-y-3">
        <h3 className="text-sm font-bold text-white uppercase tracking-wider flex items-center space-x-2">
          <ListChecks className="w-4 h-4 text-cyan-400" />
          <span>Criterios Obligatorios de Aceptación (Checklist del Auditor)</span>
        </h3>

        <div className="space-y-2.5">
          {AUDIT_CRITERIA.map((crit) => {
            const isChecked = !!completedItems[crit.id];
            return (
              <div 
                key={crit.id}
                onClick={() => toggleCheck(crit.id)}
                className={`p-4 rounded-lg border transition-all cursor-pointer flex items-start space-x-3.5 ${
                  isChecked 
                    ? 'bg-slate-900/40 border-slate-800 text-slate-400' 
                    : 'bg-slate-900 border-slate-800/90 hover:border-slate-700 text-slate-200'
                }`}
              >
                <div className="mt-0.5 text-cyan-400 shrink-0">
                  {isChecked ? (
                    <CheckSquare className="w-5 h-5 text-emerald-400" />
                  ) : (
                    <Square className="w-5 h-5 text-slate-400" />
                  )}
                </div>

                <div className="flex-1 space-y-1">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <span className={`text-xs font-bold ${isChecked ? 'line-through text-slate-400' : 'text-slate-100'}`}>
                        {crit.title}
                      </span>
                      <span className={`text-[10px] px-2 py-0.5 rounded font-mono font-bold ${
                        crit.priority === 'CRÍTICA' 
                          ? 'bg-rose-950 text-rose-300 border border-rose-800/40' 
                          : 'bg-amber-950 text-amber-300 border border-amber-800/40'
                      }`}>
                        {crit.priority}
                      </span>
                    </div>
                    <span className="text-[10px] text-slate-400 uppercase font-mono">{crit.category}</span>
                  </div>

                  <p className="text-xs text-slate-400 leading-relaxed">
                    {crit.description}
                  </p>

                  <div className="text-[11px] text-cyan-400/90 font-mono pt-1">
                    ↳ <strong>Entregable requerido:</strong> {crit.deliverable}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Non-Negotiable Auditor Rules */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 space-y-3">
        <div className="flex items-center space-x-2 text-white font-bold text-sm">
          <Scale className="w-4 h-4 text-cyan-400" />
          <span>Reglas Contables No Negociables del Auditor</span>
        </div>

        <ul className="text-xs text-slate-400 space-y-2 list-disc list-inside">
          <li><strong>Prohibido modificar saldos directamente:</strong> El campo <code>current_balance</code> en una cuenta nunca se actualiza manualmente; solo se recalcula o acumula como resultado del posteo de un asiento con líneas debitadas o acreditadas.</li>
          <li><strong>Mínimo dos partidas por comprobante:</strong> Todo asiento contable debe tener al menos una partida de débito y una partida de crédito.</li>
          <li><strong>Fechas de periodo cerrado:</strong> Si una empresa cierra el mes a día 31, no se pueden registrar asientos con fechas de ese mes o anteriores salvo reapertura autorizada.</li>
        </ul>
      </div>
    </div>
  );
};
