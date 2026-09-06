import React, { useState } from 'react';
import { Copy, Check, Terminal, FileCode } from 'lucide-react';
import { GoFile } from '../data/goSkeleton';

interface CodeViewerProps {
  file: GoFile;
}

export const CodeViewer: React.FC<CodeViewerProps> = ({ file }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(file.code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const lines = file.code.split('\n');

  // Simple token parser for clean, fast Go syntax highlighting
  const renderLine = (line: string) => {
    // Check comments first
    const commentIndex = line.indexOf('//');
    if (commentIndex !== -1) {
      const codePart = line.substring(0, commentIndex);
      const commentPart = line.substring(commentIndex);
      const isTodo = commentPart.includes('TODO:');
      return (
        <>
          {renderTokens(codePart)}
          <span className={isTodo ? 'text-amber-400 font-medium' : 'text-slate-400 italic'}>
            {commentPart}
          </span>
        </>
      );
    }
    return renderTokens(line);
  };

  const renderTokens = (text: string) => {
    // Regex for Go tokens: strings, backticks, keywords, types
    const tokenRegex = /(`[^`]*`|"[^"]*"|\b(?:package|import|type|struct|interface|func|return|var|const|if|else|range|for|nil|true|false)\b|\b(?:string|int|int64|float64|bool|error|time\.Time|context\.Context|map)\b|[{}()[\],;]|:=|=[^=]|!=|==)/g;

    const parts: React.ReactNode[] = [];
    let lastIndex = 0;
    let match: RegExpExecArray | null;

    while ((match = tokenRegex.exec(text)) !== null) {
      if (match.index > lastIndex) {
        parts.push(text.substring(lastIndex, match.index));
      }

      const token = match[0];
      if (token.startsWith('"') || token.startsWith('`')) {
        parts.push(<span key={match.index} className="text-emerald-400">{token}</span>);
      } else if (/^(package|import|type|struct|interface|func|return|var|const|if|else|range|for|nil|true|false)$/.test(token)) {
        parts.push(<span key={match.index} className="text-purple-400 font-medium">{token}</span>);
      } else if (/^(string|int|int64|float64|bool|error|time\.Time|context\.Context|map)$/.test(token)) {
        parts.push(<span key={match.index} className="text-cyan-400">{token}</span>);
      } else {
        parts.push(<span key={match.index} className="text-slate-200">{token}</span>);
      }

      lastIndex = tokenRegex.lastIndex;
    }

    if (lastIndex < text.length) {
      parts.push(text.substring(lastIndex));
    }

    return parts.length > 0 ? parts : text;
  };

  return (
    <div id="code-viewer-container" className="flex flex-col h-full bg-slate-900 border border-slate-800 rounded-lg overflow-hidden shadow-xl">
      {/* File Header Bar */}
      <div className="flex items-center justify-between px-4 py-3 bg-slate-950 border-b border-slate-800">
        <div className="flex items-center space-x-3 overflow-hidden">
          <FileCode className="w-5 h-5 text-cyan-400 shrink-0" />
          <div className="truncate">
            <span className="text-xs text-slate-400 font-mono block truncate">{file.path}</span>
            <span className="text-xs text-slate-400 block truncate">{file.description}</span>
          </div>
        </div>

        <div className="flex items-center space-x-2 shrink-0">
          <span className="text-xs px-2.5 py-1 rounded bg-slate-800 text-slate-300 font-mono">
            {lines.length} líneas
          </span>
          <button
            id="copy-code-btn"
            onClick={handleCopy}
            className="flex items-center space-x-1.5 px-3 py-1.5 bg-cyan-600 hover:bg-cyan-500 text-white rounded text-xs font-medium transition-colors cursor-pointer"
            title="Copiar código al portapapeles"
          >
            {copied ? (
              <>
                <Check className="w-3.5 h-3.5 text-white" />
                <span>Copiado</span>
              </>
            ) : (
              <>
                <Copy className="w-3.5 h-3.5 text-white" />
                <span>Copiar código</span>
              </>
            )}
          </button>
        </div>
      </div>

      {/* Code Lines with gutter */}
      <div className="flex-1 overflow-auto font-mono text-xs leading-relaxed p-4 bg-slate-950/80 selection:bg-cyan-900/50">
        <table className="w-full border-collapse">
          <tbody>
            {lines.map((line, idx) => (
              <tr key={idx} className="hover:bg-slate-800/40 transition-colors">
                <td className="w-12 pr-4 text-right select-none text-slate-400/80 align-top border-r border-slate-800/60 font-mono">
                  {idx + 1}
                </td>
                <td className="pl-4 whitespace-pre text-slate-200 align-top">
                  {renderLine(line)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Bottom info bar */}
      <div className="px-4 py-2 bg-slate-950 border-t border-slate-800 text-xs text-slate-400 flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <Terminal className="w-3.5 h-3.5 text-cyan-400" />
          <span>Esqueleto base en Go listo para implementar lógica y persistencia</span>
        </div>
        <span className="text-slate-400 font-mono">Go 1.22+ • UTF-8</span>
      </div>
    </div>
  );
};
