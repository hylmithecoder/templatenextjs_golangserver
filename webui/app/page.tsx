"use client";

import { useEffect, useState } from "react";

interface ApiHelloResponse {
  status: string;
  message: string;
}

export default function Home() {
  const [data, setData] = useState<ApiHelloResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    fetch("/api/v1/hello")
      .then((res) => {
        if (!res.ok) {
          throw new Error(`HTTP ${res.status}`);
        }
        return res.json();
      })
      .then((json: ApiHelloResponse) => {
        setData(json);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-zinc-50 font-sans dark:bg-zinc-950 p-6">
      <main className="flex w-full max-w-2xl flex-col gap-6 rounded-2xl border border-zinc-200 bg-white p-8 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
        <div className="flex flex-col gap-2">
          <span className="inline-flex w-fit items-center rounded-full bg-zinc-100 px-3 py-1 text-xs font-medium text-zinc-800 dark:bg-zinc-800 dark:text-zinc-200">
            Next.js + Golang (Fiber v3) Template
          </span>
          <h1 className="text-3xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">
            Full-Stack Single-Binary Starter
          </h1>
          <p className="text-sm text-zinc-600 dark:text-zinc-400">
            Static Next.js export (<code className="font-mono">webui/out</code>) embedded into a Go Fiber v3 binary with unified dev reverse proxy.
          </p>
        </div>

        <div className="rounded-xl border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-950">
          <div className="text-xs font-mono uppercase tracking-wider text-zinc-500 mb-2">
            GET /api/v1/hello
          </div>
          {loading && (
            <p className="text-sm text-zinc-500">Connecting to Go API...</p>
          )}
          {error && (
            <p className="text-sm text-red-600 dark:text-red-400">
              API unreachable ({error}) — start with <code className="font-mono">make run</code>
            </p>
          )}
          {data && (
            <div className="flex items-center justify-between gap-4">
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-100">
                {data.message}
              </span>
              <span className="rounded-full bg-emerald-100 px-2.5 py-0.5 text-xs font-medium text-emerald-800 dark:bg-emerald-950 dark:text-emerald-300">
                {data.status}
              </span>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
