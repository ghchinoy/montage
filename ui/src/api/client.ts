import { AssessmentReport } from "./types";

export interface AssessRequest {
  target: string;
  use_llm?: boolean;
  scaffold?: boolean;
}

export async function assessRepository(target: string, useLLM = true): Promise<AssessmentReport> {
  const res = await fetch("/api/assess", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      target,
      use_llm: useLLM,
      scaffold: true,
    }),
  });

  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(errorText || `Failed to assess repository (HTTP ${res.status})`);
  }

  return await res.json();
}

export function getExportZipUrl(target: string): string {
  return `/api/export?target=${encodeURIComponent(target)}`;
}
