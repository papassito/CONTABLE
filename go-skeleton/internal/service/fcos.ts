export type AuthorityDomain = 'SAT' | 'IMSS' | 'ISRTP_ESTATAL';

export type ExpedienteStage = 
  | 'CREATED' | 'PENDING' | 'PROCESSING' | 'VALIDATION' | 'READY'
  | 'REQUIRES_HUMAN' | 'EXECUTED' | 'EVIDENCE_CAPTURED' | 'RECONCILED'
  | 'COMPLETED' | 'FAILED' | 'EXTERNAL_UNAVAILABLE' | 'RETRY' | 'CANCELLED' | 'BLOCKED';

export interface TaxpayerHealthSummary {
  taxpayer_id: string;
  rfc: string;
  legal_name: string;
  sat_status: 'POSITIVE' | 'CREDIT_DETECTED' | 'OBLIGATION_PENDING';
  imss_status: 'POSITIVE' | 'CREDIT_DETECTED' | 'IN_VALIDATION';
  state_status: 'POSITIVE' | 'OBLIGATION_PENDING';
  payments_status: 'UP_TO_DATE' | 'EXPIRED';
}

export interface HumanInTheLoopBarrierEvent {
  event_id: string;
  tenant_id: string;
  expediente_id: string;
  expediente_code: string;
  worker_id: string;
  portal_name: string;
  barrier_type: 'CAPTCHA' | 'MFA' | 'LAYOUT_CHANGE';
  captcha_base64_image?: string;
  timestamp_utc: string;
}