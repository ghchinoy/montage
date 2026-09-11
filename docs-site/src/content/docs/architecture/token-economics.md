---
title: Token Economics & Persona Clustering
description: How Montage optimizes context windows through deterministic pre-filtering and clustered toolsets.
---

## Token Economics

Passing massive, unstructured outputs into an LLM context window causes:
- High inference costs.
- Model latency degradation.
- "Lost in the middle" reasoning failures.

Montage's **Skills with Code** pattern pairs agents with lightweight helper scripts that execute locally, filter out noise, and return compact structured JSON payloads, achieving **60–90% context reduction**.

## Persona Clustering

Exposing dozens of flat MCP tools pollutes the model's system prompt and degrades tool selection accuracy. Montage groups candidate tools into **operational personas** (e.g. *Strategy & Compliance*, *Visual Staging*, *Media Synthesis*), allowing hosts to load only the relevant tool cluster for the active workflow.
