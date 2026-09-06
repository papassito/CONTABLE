import JSZip from 'jszip';
import { GO_SKELETON_FILES } from '../data/goSkeleton';

export async function downloadGoSkeletonZip(): Promise<void> {
  const zip = new JSZip();

  // Root folder in the zip
  const root = zip.folder('contable-fix-klik');
  if (!root) return;

  for (const file of GO_SKELETON_FILES) {
    root.file(file.path, file.code);
  }

  const blob = await zip.generateAsync({ type: 'blob' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'contable-fix-klik-go-skeleton.zip';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}
