/**
 * wailsService.ts
 * Puente de comunicación e integración tipado entre React (TS) y el FCOS Kernel (Go) de Wails.
 */

// Importación directa de los bindings generados automáticamente por el motor de Wails
import {
  CreateDraft,
  PostEntry,
  ReverseEntry,
  GetEntry,
  CreateAccount,
  GetAccount,
  GetAccountByCode,
  ListAccounts
} from '../wailsjs/go/main/App';

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

export type AccountType = 'ACTIVO' | 'PASIVO' | 'PATRIMONIO' | 'INGRESO' | 'GASTO' | 'COSTO';
export type AccountStatus = 'ACTIVA' | 'INACTIVA';

export interface Account {
  id?: string;
  code: string;
  name: string;
  type: AccountType;
  status?: AccountStatus;
  accepts_move?: boolean;
  current_bal?: number; // Representado en decimales en UI (se mapea a céntimos en el Kernel)
  parent_id?: string;
  created_at?: string;
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

// --- SERIALIZADORES Y SANITIZADORES DEL PUENTE ---

/**
 * Prepara un JournalEntry convirtiendo importes de UI (decimales) a céntimos enteros de Go.
 */
export function prepareEntryForKernel(entry: JournalEntry): any {
  return {
    ...entry,
    lines: (entry.lines || []).map(line => ({
      ...line,
      debit: toCents(line.debit),
      credit: toCents(line.credit)
    }))
  };
}

/**
 * Normaliza un JournalEntry proveniente del Kernel convirtiendo céntimos enteros a decimales legibles.
 */
export function parseEntryFromKernel(entry: any): JournalEntry {
  if (!entry) return entry;
  return {
    ...entry,
    lines: (entry.lines || []).map((line: any) => ({
      ...line,
      debit: fromCents(line.debit || 0),
      credit: fromCents(line.credit || 0)
    })),
    total_debit: entry.total_debit !== undefined ? fromCents(entry.total_debit) : undefined,
    total_credit: entry.total_credit !== undefined ? fromCents(entry.total_credit) : undefined
  } as JournalEntry;
}

/**
 * Prepara una cuenta convirtiendo el saldo decimal a céntimos enteros de Go.
 */
export function prepareAccountForKernel(account: Account): any {
  return {
    ...account,
    current_bal: toCents(account.current_bal || 0)
  };
}

/**
 * Normaliza una cuenta del Kernel convirtiendo los céntimos del saldo a formato decimal legible.
 */
export function parseAccountFromKernel(account: any): Account {
  if (!account) return account;
  return {
    ...account,
    current_bal: fromCents(account.current_bal || 0)
  };
}

// --- MÉTODOS DEL SERVICIO INTEGRADO ---

/**
 * Registra un nuevo borrador de asiento contable en el sistema.
 * @param entry Estructura del asiento contable con sus líneas.
 * @returns El asiento creado normalizado en decimales.
 */
export async function createJournalDraft(entry: JournalEntry): Promise<JournalEntry> {
  try {
    const prepared = prepareEntryForKernel(entry);
    const result = await CreateDraft(prepared);
    return parseEntryFromKernel(result);
  } catch (error) {
    throw handleKernelError(error, 'Error al crear borrador contable');
  }
}

/**
 * Consulta un asiento contable por su identificador único.
 * @param entryId ID único del asiento contable.
 * @returns Estructura completa del asiento consultado en formato legible.
 */
export async function getJournalEntry(entryId: string): Promise<JournalEntry> {
  try {
    const result = await GetEntry(entryId);
    return parseEntryFromKernel(result);
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
 * @returns El asiento de reversión creado en formato legible.
 */
export async function reverseJournalEntry(entryId: string, reason: string): Promise<JournalEntry> {
  try {
    const result = await ReverseEntry(entryId, reason);
    return parseEntryFromKernel(result);
  } catch (error) {
    throw handleKernelError(error, 'Error al revertir el asiento');
  }
}

/**
 * Registra una nueva cuenta contable en el catálogo.
 * @param account Datos estructurados de la cuenta contable.
 */
export async function createAccount(account: Account): Promise<void> {
  try {
    const prepared = prepareAccountForKernel(account);
    await CreateAccount(prepared);
  } catch (error) {
    throw handleKernelError(error, 'Error al crear cuenta contable');
  }
}

/**
 * Consulta una cuenta contable mediante su identificador único.
 */
export async function getAccount(accountId: string): Promise<Account> {
  try {
    const result = await GetAccount(accountId);
    return parseAccountFromKernel(result);
  } catch (error) {
    throw handleKernelError(error, 'Error al consultar cuenta contable');
  }
}

/**
 * Consulta una cuenta contable por su código estructurado.
 */
export async function getAccountByCode(code: string): Promise<Account> {
  try {
    const result = await GetAccountByCode(code);
    return parseAccountFromKernel(result);
  } catch (error) {
    throw handleKernelError(error, 'Error al buscar cuenta por código');
  }
}

/**
 * Obtiene el catálogo completo de cuentas contables del sistema.
 */
export async function listAccounts(): Promise<Account[]> {
  try {
    const results = await ListAccounts();
    return (results || []).map((acc: any) => parseAccountFromKernel(acc));
  } catch (error) {
    throw handleKernelError(error, 'Error al listar el catálogo de cuentas');
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
  if (rawMessage.includes('ya existe') || rawMessage.includes('already exists')) {
    return new Error('Código Duplicado: Ya existe una cuenta contable registrada con ese código.');
  }
  if (rawMessage.includes('código contable inválido') || rawMessage.includes('invalid account code')) {
    return new Error('Código Inválido: El código debe iniciar con un dígito de 1 a 6.');
  }
  if (rawMessage.includes('cuenta padre no encontrada') || rawMessage.includes('parent account not found')) {
    return new Error('Jerarquía Inválida: La cuenta padre especificada para esta clasificación no existe.');
  }
  
  return new Error(`${defaultMsg}: ${rawMessage}`);
}