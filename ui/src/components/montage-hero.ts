import { LitElement, html, css } from "lit";
import { customElement, property, state } from "lit/decorators.js";

@customElement("montage-hero")
export class MontageHero extends LitElement {
  @property({ type: String }) target = "";
  @property({ type: Boolean }) loading = false;

  @state() private inputVal = "";

  static styles = css`
    :host {
      display: block;
      width: 100%;
    }

    .hero-container {
      max-width: 900px;
      margin: 0 auto;
      text-align: center;
      padding: 3rem 1.5rem 2rem;
    }

    .brand-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background: rgba(16, 185, 129, 0.1);
      border: 1px solid rgba(16, 185, 129, 0.25);
      color: #34d399;
      font-size: 0.8125rem;
      font-weight: 600;
      padding: 0.35rem 0.85rem;
      border-radius: 9999px;
      margin-bottom: 1.5rem;
      letter-spacing: 0.02em;
    }

    .brand-title {
      font-size: 3rem;
      font-weight: 800;
      letter-spacing: -0.03em;
      line-height: 1.15;
      background: linear-gradient(135deg, #ffffff 0%, #cbd5e1 50%, #94a3b8 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      margin-bottom: 1rem;
    }

    .brand-subtitle {
      font-size: 1.15rem;
      color: #94a3b8;
      max-width: 650px;
      margin: 0 auto 2.5rem;
      line-height: 1.6;
    }

    .input-wrapper {
      position: relative;
      display: flex;
      align-items: center;
      background: #0f172a;
      border: 1.5px solid #334155;
      border-radius: 14px;
      padding: 0.5rem 0.6rem 0.5rem 1.25rem;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.5);
      transition: all 0.2s ease;
    }

    .input-wrapper:focus-within {
      border-color: #10b981;
      box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2), 0 10px 25px -5px rgba(0, 0, 0, 0.5);
    }

    .github-icon {
      color: #64748b;
      margin-right: 0.75rem;
      flex-shrink: 0;
    }

    input {
      flex: 1;
      background: transparent;
      border: none;
      outline: none;
      color: #f8fafc;
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.95rem;
      min-width: 0;
    }

    input::placeholder {
      color: #475569;
    }

    .cta-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.6rem;
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
      color: #ffffff;
      border: none;
      font-size: 0.95rem;
      font-weight: 700;
      padding: 0.85rem 1.4rem;
      border-radius: 10px;
      cursor: pointer;
      transition: all 0.2s ease;
      white-space: nowrap;
      box-shadow: 0 4px 12px rgba(16, 185, 129, 0.25);
    }

    .cta-btn:hover:not(:disabled) {
      background: linear-gradient(135deg, #34d399 0%, #059669 100%);
      box-shadow: 0 6px 18px rgba(16, 185, 129, 0.4);
      transform: translateY(-1px);
    }

    .cta-btn:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    .sample-chips {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      justify-content: center;
      gap: 0.6rem;
      margin-top: 1.25rem;
    }

    .chips-label {
      font-size: 0.8125rem;
      color: #64748b;
      font-weight: 500;
    }

    .chip {
      background: #1e293b;
      border: 1px solid #334155;
      color: #cbd5e1;
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.775rem;
      padding: 0.25rem 0.65rem;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.15s ease;
    }

    .chip:hover {
      background: #273549;
      color: #38bdf8;
      border-color: #475569;
    }

    @media (max-width: 640px) {
      .brand-title {
        font-size: 2.25rem;
      }
      .input-wrapper {
        flex-direction: column;
        padding: 0.75rem;
        gap: 0.75rem;
      }
      .cta-btn {
        width: 100%;
        justify-content: center;
      }
    }
  `;

  connectedCallback() {
    super.connectedCallback();
    if (this.target) {
      this.inputVal = this.target;
    }
  }

  private handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Enter" && !this.loading) {
      this.submit();
    }
  }

  private submit() {
    const val = this.inputVal.trim();
    if (!val) return;
    this.dispatchEvent(
      new CustomEvent("submit-target", {
        detail: { target: val },
        bubbles: true,
        composed: true,
      })
    );
  }

  private setSample(sample: string) {
    this.inputVal = sample;
    this.submit();
  }

  render() {
    return html`
      <div class="hero-container">
        <div class="brand-badge">
          <span>🎞️</span>
          <span>Autonomous Agent Ecosystem Scaffolder</span>
        </div>

        <h1 class="brand-title">Deconstruct your codebase into the Agent Economy</h1>
        <p class="brand-subtitle">
          Enter any GitHub repository. Montage identifies interface maturity, mathematical trend crunchers,
          and domain policies to generate <strong>Pure Skills</strong>, <strong>Skills with Code</strong>, and <strong>MCP Servers</strong>.
        </p>

        <div class="input-wrapper">
          <svg class="github-icon" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
          </svg>
          <input
            type="text"
            placeholder="https://github.com/owner/repo or owner/repo"
            .value="${this.inputVal}"
            @input="${(e: any) => (this.inputVal = e.target.value)}"
            @keydown="${this.handleKeyDown}"
            ?disabled="${this.loading}"
          />
          <button class="cta-btn ${this.loading ? '' : 'pulse-glow'}" @click="${this.submit}" ?disabled="${this.loading}">
            ${this.loading
              ? html`
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="animation: spin 1s linear infinite;">
                    <style>@keyframes spin { 100% { transform: rotate(360deg); } }</style>
                    <path d="M21 12a9 9 0 11-6.219-8.56"/>
                  </svg>
                  <span>Deconstructing...</span>
                `
              : html`
                  <span>Make this part of the Agent Economy</span>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M5 12h14M12 5l7 7-7 7"/>
                  </svg>
                `}
          </button>
        </div>

        <div class="sample-chips">
          <span class="chips-label">Try sample repositories:</span>
          <button class="chip" @click="${() => this.setSample('ghchinoy/syntaxis')}">ghchinoy/syntaxis</button>
          <button class="chip" @click="${() => this.setSample('ghchinoy/kitten-haven')}">ghchinoy/kitten-haven</button>
          <button class="chip" @click="${() => this.setSample('ghchinoy/repotographer')}">ghchinoy/repotographer</button>
        </div>
      </div>
    `;
  }
}
