export interface ForensicEvidence {
  ruta: string;
  tamaño_bytes: number;
  sha256: string;
  creacion_utc: string;
  ultima_modificacion_utc: string;
  ultimo_acceso_utc: string;
}

export interface ForensicHeader {
  fecha_ejecucion_utc: string;
  sistema_operativo: string;
  arquitectura: string;
  nombre_equipo: string;
}

export interface ForensicReport {
  encabezado_auditoria: ForensicHeader;
  evidencia_archivos: ForensicEvidence[];
}

export interface AuditItem {
  id: string;
  category: 'workspace' | 'go-backend' | 'forensic' | 'wails' | 'cicd' | 'powershell';
  title: string;
  severity: 'critical' | 'warning' | 'info' | 'success';
  status: 'saneado' | 'pendiente_code_assist' | 'verificado' | 'atencion_requerida';
  summary: string;
  technicalDetails: string;
  targetFiles: string[];
  suggestedActionForCodeAssist?: string;
}
