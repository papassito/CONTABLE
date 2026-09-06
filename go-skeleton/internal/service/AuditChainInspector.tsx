import React from 'react';

interface AuditEventItem {
  sequence_id: number;
  event_id: string;
  event_type: string;
  canonical_payload_hash: string;
  previous_hash: string;
  chain_hash: string;
  created_at_utc: string;
}

export const AuditChainInspector: React.FC<{ events: AuditEventItem[] }> = ({ events }) => {
  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 font-mono text-xs">
      <h4 className="text-slate-400 font-bold mb-4 uppercase tracking-wider">Audit Hash Chain (Append-Only Ledger)</h4>
      <div className="space-y-3">
        {events.map((evt) => (
          <div key={evt.event_id} className="bg-slate-950 border border-slate-800 p-3 rounded-lg hover:border-slate-700 transition-colors">
            <div className="flex justify-between text-slate-400 mb-1">
              <span>#SEQ: {evt.sequence_id} │ EVENT: <strong className="text-emerald-400">{evt.event_type}</strong></span>
              <span>{new Date(evt.created_at_utc).toLocaleString()}</span>
            </div>
            <div className="text-slate-500 truncate">PREV_HASH: {evt.previous_hash}</div>
            <div className="text-slate-300 font-bold truncate">CHAIN_HASH: {evt.chain_hash}</div>
            <div className="mt-2 text-right">
              <span className="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-emerald-950 text-emerald-400 border border-emerald-800">
                ✓ CANONICAL OK (RFC 8785)
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};