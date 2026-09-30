import { useState } from 'react';
import { Copy, Check } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { copyText } from '@/lib/clipboard';

interface CopyEmailButtonProps {
  email: string;
}

export function CopyEmailButton({ email }: CopyEmailButtonProps) {
  const [copied, setCopied] = useState(false);

  if (!email) return null;

  async function handleClick(e: React.MouseEvent) {
    e.stopPropagation();
    const ok = await copyText(email);
    if (ok) {
      setCopied(true);
      toast.success('E-mail copiado');
      setTimeout(() => setCopied(false), 1200);
    } else {
      toast.error('Não foi possível copiar o e-mail');
    }
  }

  return (
    <Button
      variant="ghost"
      size="icon-xs"
      className="opacity-0 group-hover:opacity-100 focus-visible:opacity-100"
      onClick={handleClick}
      aria-label="Copiar e-mail"
    >
      {copied ? (
        <Check className="size-3.5 text-emerald-500" />
      ) : (
        <Copy className="size-3.5 text-muted-foreground" />
      )}
    </Button>
  );
}
