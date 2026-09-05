import { LitElement, html, css } from "lit";
import { customElement, state } from "lit/decorators.js";
import { assessRepository, getExportZipUrl } from "./api/client";
import { AssessmentReport } from "./api/types";

import "./components/montage-hero";
import "./components/montage-pipeline";
import "./components/montage-scorecard";
import "./components/montage-capabilities";
import "./components/montage-skills";
import "./components/montage-mcp";

@customElement("montage-app")
export class MontageApp extends LitElement {
  @state() private report: AssessmentReport | null = null;
  @state() private loading = false;
  @state() private errorMessage = "";
  @state() private activeTab = "overview";
  @state() private currentTarget = "";

  static styles = css`
    :host {
      display: flex;
      flex-direction: column;
      min-height: 100vh;
    }

    header {
      border-bottom: 1px solid #1e293b;
      background: rgba(15, 23, 42, 0.75);
      backdrop-filter: blur(12px);
      position: sticky;
      top: 0;
      z-index: 50;
    }

    .header-inner {
      max-width: 1200px;
      margin: 0 auto;
      padding: 0.85rem 1.5rem;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .brand {
      display: flex;
      align-items: center;
      gap: 0.6rem;
      text-decoration: none;
      color: #f8fafc;
      font-weight: 800;
      font-size: 1.2rem;
      letter-spacing: -0.02em;
    }

    .brand-logo {
      font-size: 1.3rem;
    }

    .version-tag {
      font-size: 0.7rem;
      font-weight: 700;
      background: #1e293b;
      color: #94a3b8;
      padding: 0.15rem 0.45rem;
      border-radius: 4px;
      border: 1px solid #334155;
    }

    .header-links {
      display: flex;
      align-items: center;
      gap: 1.25rem;
    }

    .header-link {
      color: #94a3b8;
      text-decoration: none;
      font-size: 0.875rem;
      font-weight: 600;
      transition: color 0.15s ease;
    }

    .header-link:hover {
      color: #f8fafc;
    }

    main {
      flex: 1;
      max-width: 1200px;
      width: 100%;
      margin: 0 auto;
      padding: 0 1.5rem 4rem;
    }

    .error-banner {
      background: rgba(244, 63, 94, 0.12);
      border: 1px solid rgba(244, 63, 94, 0.35);
      color: #fca5a5;
      padding: 1rem 1.25rem;
      border-radius: 10px;
      margin: 1.5rem auto;
      max-width: 900px;
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-size: 0.9rem;
    }

    .tabs-nav {
      display: flex;
      align-items: center;
      justify-content: space-between;
      border-bottom: 1px solid #1e293b;
      margin-bottom: 2rem;
      padding-bottom: 0.5rem;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .tabs-list {
      display: flex;
      gap: 0.5rem;
    }

    .tab-btn {
      background: transparent;
      border: none;
      color: #94a3b8;
      font-size: 0.9rem;
      font-weight: 600;
      padding: 0.5rem 1rem;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.15s ease;
    }

    .tab-btn:hover {
      color: #f8fafc;
      background: #1e293b;
    }

    .tab-btn.active {
      color: #10b981;
      background: rgba(16, 185, 129, 0.12);
      border: 1px solid rgba(16, 185, 129, 0.25);
    }

    .download-bundle-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background: #1e293b;
      border: 1px solid #334155;
      color: #38bdf8;
      padding: 0.5rem 1rem;
      border-radius: 8px;
      font-size: 0.85rem;
      font-weight: 700;
      text-decoration: none;
      transition: all 0.2s ease;
    }

    .download-bundle-btn:hover {
      background: #273549;
      border-color: #0284c7;
      color: #7dd3fc;
    }

    footer {
      border-top: 1px solid #1e293b;
      padding: 1.75rem 1.5rem;
      text-align: center;
      font-size: 0.8125rem;
      color: #64748b;
    }
  `;

  private async handleAssessment(e: CustomEvent<{ target: string }>) {
    const target = e.detail.target;
    this.currentTarget = target;
    this.loading = true;
    this.errorMessage = "";

    try {
      const data = await assessRepository(target, true);
      this.report = data;
      this.activeTab = "overview";
    } catch (err: any) {
      this.errorMessage = err.message || "Failed to assess repository";
    } finally {
      this.loading = false;
    }
  }

  render() {
    const hasReport = !!this.report;

    return html`
      <header>
        <div class="header-inner">
          <a href="/" class="brand">
            <span class="brand-logo">🎞️</span>
            <span>Montage</span>
            <span class="version-tag">v0.2.0</span>
          </a>

          <div class="header-links">
            <a href="https://github.com/ghchinoy/montage" target="_blank" rel="noreferrer" class="header-link">GitHub</a>
            <a href="https://github.com/agentplugins/agent-plugins-spec" target="_blank" rel="noreferrer" class="header-link">Plugins Spec</a>
          </div>
        </div>
      </header>

      <main>
        <montage-hero
          .target="${this.currentTarget}"
          .loading="${this.loading}"
          @submit-target="${this.handleAssessment}"
        ></montage-hero>

        ${this.loading ? html`<montage-pipeline></montage-pipeline>` : ""}

        ${this.errorMessage
          ? html`
              <div class="error-banner">
                <span>⚠️</span>
                <span>${this.errorMessage}</span>
              </div>
            `
          : ""}

        ${hasReport && !this.loading
          ? html`
              <div class="tabs-nav">
                <div class="tabs-list">
                  <button
                    class="tab-btn ${this.activeTab === 'overview' ? 'active' : ''}"
                    @click="${() => (this.activeTab = 'overview')}"
                  >
                    Scorecard & Overview
                  </button>
                  <button
                    class="tab-btn ${this.activeTab === 'capabilities' ? 'active' : ''}"
                    @click="${() => (this.activeTab = 'capabilities')}"
                  >
                    Capabilities Matrix (${this.report!.capabilities?.length || 0})
                  </button>
                  <button
                    class="tab-btn ${this.activeTab === 'skills' ? 'active' : ''}"
                    @click="${() => (this.activeTab = 'skills')}"
                  >
                    Specialized Skills (${this.report!.skill_blueprints?.length || 0})
                  </button>
                  <button
                    class="tab-btn ${this.activeTab === 'mcp' ? 'active' : ''}"
                    @click="${() => (this.activeTab = 'mcp')}"
                  >
                    MCP Toolsets (${this.report!.tool_groups?.length || 0})
                  </button>
                </div>

                <a
                  href="${getExportZipUrl(this.currentTarget)}"
                  class="download-bundle-btn"
                  download
                >
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3"/>
                  </svg>
                  <span>Download Agent Bundle (.zip)</span>
                </a>
              </div>

              ${this.activeTab === "overview"
                ? html`<montage-scorecard .report="${this.report}"></montage-scorecard>`
                : ""}

              ${this.activeTab === "capabilities"
                ? html`<montage-capabilities .capabilities="${this.report!.capabilities || []}"></montage-capabilities>`
                : ""}

              ${this.activeTab === "skills"
                ? html`<montage-skills .skills="${this.report!.skill_blueprints || []}"></montage-skills>`
                : ""}

              ${this.activeTab === "mcp"
                ? html`<montage-mcp .toolGroups="${this.report!.tool_groups || []}" .targetName="${this.report!.target.name}"></montage-mcp>`
                : ""}
            `
          : ""}
      </main>

      <footer>
        Montage &copy; 2026 &middot; Autonomous Agent Ecosystem Tooling Assessor &middot; Built with Lit & Go
      </footer>
    `;
  }
}
