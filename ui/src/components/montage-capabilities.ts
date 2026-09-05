import { LitElement, html, css } from "lit";
import { customElement, property, state } from "lit/decorators.js";
import { Capability } from "../api/types";

@customElement("montage-capabilities")
export class MontageCapabilities extends LitElement {
  @property({ type: Array }) capabilities: Capability[] = [];

  @state() private activeFilter = "all";
  @state() private searchQuery = "";

  static styles = css`
    :host {
      display: block;
      width: 100%;
    }

    .matrix-card {
      background: #0f172a;
      border: 1px solid #1e293b;
      border-radius: 14px;
      overflow: hidden;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
    }

    .matrix-header {
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid #1e293b;
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
    }

    .filter-tabs {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;
    }

    .filter-btn {
      background: #1e293b;
      border: 1px solid #334155;
      color: #94a3b8;
      font-size: 0.8125rem;
      font-weight: 600;
      padding: 0.35rem 0.75rem;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.15s ease;
    }

    .filter-btn:hover {
      background: #273549;
      color: #f8fafc;
    }

    .filter-btn.active {
      background: rgba(16, 185, 129, 0.15);
      border-color: #10b981;
      color: #34d399;
    }

    .search-box {
      background: #1e293b;
      border: 1px solid #334155;
      border-radius: 8px;
      padding: 0.35rem 0.75rem;
      color: #f8fafc;
      font-size: 0.8125rem;
      font-family: 'JetBrains Mono', monospace;
      outline: none;
      min-width: 220px;
    }

    .search-box:focus {
      border-color: #06b6d4;
    }

    .table-container {
      overflow-x: auto;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      text-align: left;
      font-size: 0.875rem;
    }

    th {
      background: #161f30;
      color: #64748b;
      font-weight: 700;
      font-size: 0.75rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      padding: 0.85rem 1.25rem;
      border-bottom: 1px solid #1e293b;
      white-space: nowrap;
    }

    td {
      padding: 1rem 1.25rem;
      border-bottom: 1px solid #162032;
      vertical-align: top;
      color: #cbd5e1;
    }

    tr:hover td {
      background: rgba(30, 41, 59, 0.4);
    }

    .cap-name {
      font-weight: 700;
      color: #f8fafc;
      display: block;
      margin-bottom: 0.2rem;
    }

    .cap-desc {
      font-size: 0.8rem;
      color: #94a3b8;
      line-height: 1.4;
    }

    .source-code {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.775rem;
      background: #1e293b;
      padding: 0.2rem 0.4rem;
      border-radius: 4px;
      color: #38bdf8;
      border: 1px solid #334155;
      display: inline-block;
      white-space: nowrap;
    }

    .badge {
      display: inline-flex;
      align-items: center;
      gap: 0.3rem;
      padding: 0.25rem 0.6rem;
      border-radius: 6px;
      font-size: 0.75rem;
      font-weight: 700;
      white-space: nowrap;
    }

    .badge-pure {
      background: rgba(147, 51, 234, 0.15);
      border: 1px solid rgba(147, 51, 234, 0.35);
      color: #c084fc;
    }

    .badge-code {
      background: rgba(245, 158, 11, 0.15);
      border: 1px solid rgba(245, 158, 11, 0.35);
      color: #fbbf24;
    }

    .badge-mcp {
      background: rgba(6, 182, 212, 0.15);
      border: 1px solid rgba(6, 182, 212, 0.35);
      color: #22d3ee;
    }

    .badge-cli {
      background: rgba(100, 116, 139, 0.2);
      border: 1px solid rgba(100, 116, 139, 0.4);
      color: #cbd5e1;
    }

    .nature-badge {
      font-size: 0.725rem;
      font-weight: 600;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
      display: inline-block;
      white-space: nowrap;
    }

    .nature-readonly {
      background: rgba(16, 185, 129, 0.12);
      color: #34d399;
      border: 1px solid rgba(16, 185, 129, 0.25);
    }

    .nature-mutating {
      background: rgba(244, 63, 94, 0.12);
      color: #f87171;
      border: 1px solid rgba(244, 63, 94, 0.25);
    }

    .domain-tag {
      font-size: 0.75rem;
      color: #94a3b8;
      font-weight: 600;
      background: #1e293b;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
      display: inline-block;
    }

    .empty-state {
      padding: 3rem 1.5rem;
      text-align: center;
      color: #64748b;
    }
  `;

  render() {
    const counts = {
      all: this.capabilities.length,
      pure_skill: this.capabilities.filter((c) => c.recommended_target === "pure_skill").length,
      skill_with_code: this.capabilities.filter((c) => c.recommended_target === "skill_with_code").length,
      mcp_tool: this.capabilities.filter((c) => c.recommended_target === "mcp_tool").length,
      cli_command: this.capabilities.filter((c) => c.recommended_target === "cli_command").length,
    };

    let filtered = this.capabilities;
    if (this.activeFilter !== "all") {
      filtered = filtered.filter((c) => c.recommended_target === this.activeFilter);
    }
    if (this.searchQuery.trim()) {
      const q = this.searchQuery.toLowerCase();
      filtered = filtered.filter(
        (c) =>
          c.name.toLowerCase().includes(q) ||
          c.description.toLowerCase().includes(q) ||
          c.source_ref.toLowerCase().includes(q) ||
          c.target_grouping.toLowerCase().includes(q)
      );
    }

    return html`
      <div class="matrix-card">
        <div class="matrix-header">
          <div class="filter-tabs">
            <button
              class="filter-btn ${this.activeFilter === 'all' ? 'active' : ''}"
              @click="${() => (this.activeFilter = 'all')}"
            >
              All (${counts.all})
            </button>
            <button
              class="filter-btn ${this.activeFilter === 'pure_skill' ? 'active' : ''}"
              @click="${() => (this.activeFilter = 'pure_skill')}"
            >
              📝 Pure Skills (${counts.pure_skill})
            </button>
            <button
              class="filter-btn ${this.activeFilter === 'skill_with_code' ? 'active' : ''}"
              @click="${() => (this.activeFilter = 'skill_with_code')}"
            >
              ⚡ Skills with Code (${counts.skill_with_code})
            </button>
            <button
              class="filter-btn ${this.activeFilter === 'mcp_tool' ? 'active' : ''}"
              @click="${() => (this.activeFilter = 'mcp_tool')}"
            >
              🔌 MCP Tools (${counts.mcp_tool})
            </button>
            <button
              class="filter-btn ${this.activeFilter === 'cli_command' ? 'active' : ''}"
              @click="${() => (this.activeFilter = 'cli_command')}"
            >
              💻 CLI Tools (${counts.cli_command})
            </button>
          </div>

          <input
            type="text"
            class="search-box"
            placeholder="Search capabilities..."
            .value="${this.searchQuery}"
            @input="${(e: any) => (this.searchQuery = e.target.value)}"
          />
        </div>

        <div class="table-container">
          ${filtered.length === 0
            ? html`<div class="empty-state">No capabilities match your active filter.</div>`
            : html`
                <table>
                  <thead>
                    <tr>
                      <th>Capability</th>
                      <th>Source Reference</th>
                      <th>Nature</th>
                      <th>Target Ecosystem Artifact</th>
                      <th>Domain Grouping</th>
                      <th>Architectural Rationale</th>
                    </tr>
                  </thead>
                  <tbody>
                    ${filtered.map((cap) => {
                      let targetBadge = html`<span class="badge badge-mcp">🔌 MCP Tool</span>`;
                      if (cap.recommended_target === "pure_skill") {
                        targetBadge = html`<span class="badge badge-pure">📝 Pure Skill</span>`;
                      } else if (cap.recommended_target === "skill_with_code") {
                        targetBadge = html`<span class="badge badge-code">⚡ Skill with Code</span>`;
                      } else if (cap.recommended_target === "cli_command") {
                        targetBadge = html`<span class="badge badge-cli">💻 CLI Tool</span>`;
                      }

                      const natureBadge = cap.is_mutating
                        ? html`<span class="nature-badge nature-mutating">⚡ Mutating</span>`
                        : html`<span class="nature-badge nature-readonly">🔒 Read-Only</span>`;

                      return html`
                        <tr>
                          <td>
                            <span class="cap-name">${cap.name}</span>
                            <span class="cap-desc">${cap.description}</span>
                          </td>
                          <td>
                            <span class="source-code">${cap.source_ref}</span>
                          </td>
                          <td>${natureBadge}</td>
                          <td>${targetBadge}</td>
                          <td>
                            <span class="domain-tag">${cap.target_grouping}</span>
                          </td>
                          <td style="font-size: 0.8125rem; color: #94a3b8; max-width: 280px;">
                            ${cap.rationale}
                          </td>
                        </tr>
                      `;
                    })}
                  </tbody>
                </table>
              `}
        </div>
      </div>
    `;
  }
}
