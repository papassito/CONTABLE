import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { HumanInTheLoopModal } from './HumanInTheLoopModal';
import { useEventStreamStore } from '../../stores/useEventStreamStore';

vi.mock('../../stores/useEventStreamStore');

describe('UI/UX: HumanInTheLoopModal Interceptor', () => {
  it('debe permanecer oculto si no existen barreras RPA activas', () => {
    (useEventStreamStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      currentBarrier: null,
      resolveBarrier: vi.fn(),
    });

    const { container } = render(<HumanInTheLoopModal />);
    expect(container.firstChild).toBeNull();
  });

  it('debe desplegar el modal con CAPTCHA al recibir un evento IntervencionHumanaRequerida', () => {
    const resolveMock = vi.fn();
    (useEventStreamStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      currentBarrier: {
        event_id: 'evt-100',
        portal_name: 'Portal IMSS SIPARE',
        expediente_code: 'EXP-2026-09-001',
        barrier_type: 'CAPTCHA',
        captcha_base64_image: 'data:image/png;base64,iVBORw0KGgo...',
        worker_id: 'worker-01',
      },
      resolveBarrier: resolveMock,
    });

    render(<HumanInTheLoopModal />);

    expect(screen.getByText(/Intervención Humana Requerida/i)).toBeInTheDocument();
    expect(screen.getByText(/Portal IMSS SIPARE/i)).toBeInTheDocument();

    const input = screen.getByPlaceholderText(/Ingrese el texto de la imagen/i);
    fireEvent.change(input, { target: { value: 'AB12CD' } });

    const submitBtn = screen.getByText(/Confirmar y Reanudar Robot/i);
    fireEvent.click(submitBtn);

    expect(resolveMock).toHaveBeenCalledWith('worker-01', 'AB12CD');
  });
});