import React, { useState } from 'react';
import { useEventStreamStore } from '../../stores/useEventStreamStore';

export const HumanInTheLoopModal: React.FC = () => {
  const { currentBarrier, resolveBarrier } = useEventStreamStore();
  const [captchaValue, setCaptchaValue] = useState('');

  if (!currentBarrier) return null;

  return (
    <div className="fixed inset-0 bg-slate-900/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-slate-800 border border-amber-500/50 rounded-xl max-w-lg w-full p-6 shadow-2xl">
        <div className="flex items-center gap-3 text-amber-400 mb-4">
          <span className="text-2xl">⚠️</span>
          <h3 className="text-lg font-bold">Intervención Humana Requerida (RPA Barrier)</h3>
        </div>
        
        <p className="text-slate-300 text-sm mb-4">
          El conector oficial <strong className="text-white">{currentBarrier.portal_name}</strong> requiere validación para el trámite de <strong className="text-white">{currentBarrier.expediente_code}</strong>.
        </p>

        {currentBarrier.barrier_type === 'CAPTCHA' && (
          <div className="space-y-4">
            <div className="bg-slate-950 p-4 rounded-lg flex justify-center border border-slate-700">
              <img src={currentBarrier.captcha_base64_image} alt="CAPTCHA Portal Oficial" className="h-16 object-contain" />
            </div>
            <input
              type="text"
              value={captchaValue}
              onChange={(e) => setCaptchaValue(e.target.value)}
              placeholder="Ingrese el texto de la imagen..."
              className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-amber-500"
              autoFocus
            />
          </div>
        )}

        <div className="mt-6 flex justify-end gap-3">
          <button
            onClick={() => resolveBarrier(currentBarrier.worker_id, captchaValue)}
            className="bg-amber-500 hover:bg-amber-600 text-slate-950 font-bold px-5 py-2 rounded-lg transition-colors"
          >
            Confirmar y Reanudar Robot
          </button>
        </div>
      </div>
    </div>
  );
};