// DEPRECATED: This file is no longer the source of truth. The UI now fetches file structure dynamically from the /api/v1/architecture/files endpoint.
export interface GoFile {
  path: string;
  name: string;
  layer: 'cmd' | 'domain' | 'repository' | 'service' | 'handler' | 'config' | 'pkg' | 'root';
  layerLabel: string;
  description: string;
  code: string;
}

export const GO_SKELETON_FILES: GoFile[] = [];
