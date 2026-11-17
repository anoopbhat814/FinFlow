export interface Page<T> {
  items: T[];
  nextCursor: string | null;
}

export interface ApiErrorBody {
  code: string;
  message: string;
  requestId?: string;
  fields?: Record<string, string>;
}

export class ApiError extends Error {
  constructor(public readonly status: number, public readonly body: ApiErrorBody) {
    super(body.message);
    this.name = 'ApiError';
  }
}
