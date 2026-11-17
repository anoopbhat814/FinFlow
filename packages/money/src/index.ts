import type { Currency } from '@finflow/types';

const EXPONENTS: Record<Currency, number> = { USD: 2, EUR: 2, GBP: 2, CAD: 2, JPY: 0 };

export class Money {
  private constructor(readonly minor: number, readonly currency: Currency) {
    if (!Number.isSafeInteger(minor)) throw new RangeError('minor units must be a safe integer');
  }

  static fromMinor(minor: number, currency: Currency): Money {
    return new Money(minor, currency);
  }

  static parse(input: string, currency: Currency): Money {
    const exp = EXPONENTS[currency];
    if (!/^-?\d+(\.\d+)?$/.test(input)) throw new SyntaxError(`invalid amount: ${input}`);
    const [whole = '0', frac = ''] = input.split('.');
    if (frac.length > exp) throw new RangeError(`too many decimals for ${currency}`);
    const sign = whole.startsWith('-') ? -1 : 1;
    const minor = Math.abs(Number(whole)) * 10 ** exp + Number(frac.padEnd(exp, '0') || 0);
    return new Money(sign * minor, currency);
  }

  add(other: Money): Money {
    this.assertSame(other);
    return new Money(this.minor + other.minor, this.currency);
  }

  subtract(other: Money): Money {
    this.assertSame(other);
    return new Money(this.minor - other.minor, this.currency);
  }

  /** Splits into n parts, distributing the remainder one unit at a time. */
  allocate(n: number): Money[] {
    const base = Math.trunc(this.minor / n);
    let rem = this.minor - base * n;
    return Array.from({ length: n }, () => {
      const extra = rem > 0 ? 1 : rem < 0 ? -1 : 0;
      rem -= extra;
      return new Money(base + extra, this.currency);
    });
  }

  format(locale = 'en-US'): string {
    const exp = EXPONENTS[this.currency];
    return new Intl.NumberFormat(locale, { style: 'currency', currency: this.currency }).format(
      this.minor / 10 ** exp,
    );
  }

  private assertSame(other: Money) {
    if (other.currency !== this.currency) throw new TypeError('currency mismatch');
  }
}
