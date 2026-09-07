/**
 * wailsService.ts
 * Puente de comunicación e integración tipado entre React (TS) y el FCOS Kernel (Go) de Wails.
 */

// Importación directa de los bindings generados automáticamente por el motor de Wails
import { CreateDraft, PostEntry, ReverseEntry, GetEntry } from '../wailsjs/go/main/App';

// --- CONTRATOS DE TIPOS (Coincidentes con Go Domain) ---

export type EntryStatus = 'BORRADOR' | 'CONTABILIZADO' | 'ANULADO';

export interface JournalLine {
  account_id: string;
  account_code: string;
  description: string;
  debit: number;  // Representado en céntimos (ej. $10.00 = 1000)
  credit: number; // Representado en céntimos (ej. $10.00 = 1000)
  third_party_id?: string;
}

export interface JournalEntry {
  id?: string;
  number: string;
  date: string; // Formato RFC3339 / ISO string
  concept: string;
  reference: string;
  status?: EntryStatus;
  lines: JournalLine[];
  total_debit?: number;
  total_credit?: number;
}

// --- HELPER DE CONVERSIÓN MONETARIA ---

/**
 * Convierte un monto decimal de UI (ej: 100.50) a céntimos enteros de Go (ej: 10050)
 */
export const toCents = (amount: number): number => {
  return Math.round(amount * 100);
};

/**
 * Convierte céntimos enteros de Go (ej: 10050) a decimal legible en UI (ej: 100.50)
 */
export const fromCents = (cents: number): number => {
  return cents / 100;
};

// --- MÉTODOS DEL SERVICIO INTEGRADO ---

/**
 * Registra un nuevo borrador de asiento contable en el sistema.
 * @param entry Estructura del asiento contable con sus líneas.
 * @returns El asiento creado con su respectivo ID autogenerado.
 */
export async function createJournalDraft(entry: JournalEntry): Promise<JournalEntry> {
  try {
    // Invocar el binding nativo de Go
    const result = await CreateDraft(entry as any);
    return result as unknown as JournalEntry;
  } catch (error) {
    throw handleKernelError(error, 'Error al crear borrador contable');
  }
}

/**
 * Consulta un asiento contable por su identificador único.
 * @param entryId ID único del asiento contable.
 * @returns Estructura completa del asiento consultado.
 */
export async function getJournalEntry(entryId: string): Promise<JournalEntry> {
  try {
    const result = await GetEntry(entryId);
    return result as unknown as JournalEntry;
  } catch (error) {
    throw handleKernelError(error, 'Error al consultar asiento contable');
  }
}

/**
 * Contabiliza de forma definitiva e irreversible un asiento en borrador.
 * @param entryId ID único del asiento contable.
 */
export async function postJournalEntry(entryId: string): Promise<void> {
  try {
    await PostEntry(entryId);
  } catch (error) {
    throw handleKernelError(error, 'Error al contabilizar asiento definitivo');
  }
}

/**
 * Genera un contraasiento de reversión/anulación con saldo inverso.
 * @param entryId ID único del asiento original.
 * @param reason Motivo justificado de la anulación para la pista de auditoría.
 */
export async function reverseJournalEntry(entryId: string, reason: string): Promise<JournalEntry> {
  try {
    const result = await ReverseEntry(entryId, reason);
    return result as unknown as JournalEntry;
  } catch (error) {
    throw handleKernelError(error, 'Error al revertir el asiento');
  }
}

// --- GESTOR DE ERRORES DEL KERNEL CONTABLE ---

function handleKernelError(error: any, defaultMsg: string): Error {
  const rawMessage = typeof error === 'string' ? error : error?.message || '';
  
  if (rawMessage.includes('unbalanced journal')) {
    return new Error('Violación de Partida Doble: Los débitos y créditos no cuadran.');
  }
  if (rawMessage.includes('closed period')) {
    return new Error('Periodo Cerrado: El periodo contable correspondiente está cerrado o ya fue declarado.');
  }
  if (rawMessage.includes('account not found')) {
    return new Error('Cuenta Inexistente: Una o más cuentas especificadas no existen.');
  }
  
  return new Error(`${defaultMsg}: ${rawMessage}`);
}