// Masks a Brazilian CPF as ___.___.___-__. Returns the input unchanged when it
// isn't 11 digits (already-formatted values, partial/empty data) so nothing breaks.
export function formatCpf(value?: string | null): string | undefined {
  if (!value) return undefined;
  const digits = value.replace(/\D/g, '');
  if (digits.length !== 11) return value;
  return digits.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
}
