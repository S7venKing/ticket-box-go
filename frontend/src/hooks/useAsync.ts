import { useCallback, useState } from "react";

export function useAsync<T>(operation: () => Promise<T>) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const execute = useCallback(async () => {
    setLoading(true); setError(null);
    try { return await operation(); } catch (err) { const message = err instanceof Error ? err.message : "Unexpected error"; setError(message); throw err; }
    finally { setLoading(false); }
  }, [operation]);
  return { execute, loading, error };
}
