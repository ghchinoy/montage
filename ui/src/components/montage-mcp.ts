import { LitElement, html, css } from "lit";
import { customElement, property, state } from "lit/decorators.js";
import { ToolGroup } from "../api/types";

@customElement("montage-mcp")
export class MontageMCP extends LitElement {
  @property({ type: Array }) toolGroups: ToolGroup[] = [];
  @property({ type: String }) targetName = "app";

  @state() private activeClient = "opencode";
  @state() private copied = false;

  static styles = css`
    :host {
      display: block;
      width: 100%;
    }

    .mcp-layout {
      display: flex;
      flex-direction: column;
      gap: 2rem;
    }

    .group-card {
      background: #0f172a;
      border: 1px solid #1e293b;
      border-radius: 14px;
      overflow: hidden;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
    }

    .group-header {
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid #1e293b;
      background: #161f30;
    }

    .group-title {
      font-size: 1.1rem;
      font-weight: 700;
      color: #f8fafc;
      margin-bottom: 0.25rem;
    }

    .group-desc {
      font-size: 0.825rem;
      color: #94a3b8;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      text-align: left;
      font-size: 0.875rem;
    }

    th {
      background: #0f172a;
      color: #64748b;
      font-weight: 700;
      font-size: 0.75rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      padding: 0.75rem 1.25rem;
      border-bottom: 1px solid #1e293b;
    }

    td {
      padding: 0.85rem 1.25rem;
      border-bottom: 1px solid #162032;
      color: #cbd5e1;
    }

    .tool-name {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.85rem;
      font-weight: 700;
      color: #38bdf8;
    }

    .badge-ro {
      font-size: 0.725rem;
      font-weight: 600;
      background: rgba(16, 185, 129, 0.12);
      color: #34d399;
      border: 1px solid rgba(16, 185, 129, 0.25);
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
    }

    .badge-mut {
      font-size: 0.725rem;
      font-weight: 600;
      background: rgba(244, 63, 94, 0.12);
      color: #f87171;
      border: 1px solid rgba(244, 63, 94, 0.25);
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
    }

    .config-card {
      background: #0f172a;
      border: 1px solid #1e293b;
      border-radius: 14px;
      padding: 1.5rem;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
    }

    .config-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 1rem;
    }

    .config-title {
      font-size: 1rem;
      font-weight: 700;
      color: #f8fafc;
    }

    .client-tabs {
      display: flex;
      gap: 0.4rem;
    }

    .tab-btn {
      background: #1e293b;
      border: 1px solid #334155;
      color: #94a3b8;
      font-size: 0.8rem;
      font-weight: 600;
      padding: 0.3rem 0.7rem;
      border-radius: 6px;
      cursor: pointer;
    }

    .tab-btn.active {
      background: rgba(6, 182, 212, 0.15);
      border-color: #06b6d4;
      color: #22d3ee;
    }

    .code-container {
      position: relative;
      background: #090d16;
      border: 1px solid #1e293b;
      border-radius: 8px;
      padding: 1rem;
      overflow-x: auto;
    }

    pre {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.825rem;
      color: #cbd5e1;
      line-height: 1.5;
    }

    .copy-btn {
      position: absolute;
      top: 0.75rem;
      right: 0.75rem;
      background: #1e293b;
      border: 1px solid #334155;
      color: #cbd5e1;
      font-size: 0.75rem;
      font-weight: 600;
      padding: 0.3rem 0.6rem;
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.15s ease;
    }

    .copy-btn:hover {
      background: #273549;
      color: #f8fafc;
    }
  `;

  private getConfigSnippet(): string {
    const app = this.targetName.toLowerCase();
    switch (this.activeClient) {
      case "opencode":
        return JSON.stringify(
          {
            mcpServers: {
              [app]: {
                command: `bin/${app}`,
                args: ["mcp"],
              },
            },
          },
          null,
          2
        );
      case "claude":
        return JSON.stringify(
          {
            mcpServers: {
              [app]: {
                command: `/usr/local/bin/${app}`,
                args: ["mcp"],
              },
            },
          },
          null,
          2
        );
      case "cursor":
        return JSON.stringify(
          {
            mcpServers: {
              [app]: {
                command: `${app}`,
                args: ["mcp"],
              },
            },
          },
          null,
          2
        );
      default:
        return "";
    }
  }

  private async copyConfig() {
    await navigator.clipboard.writeText(this.getConfigSnippet());
    this.copied = true;
    setTimeout(() => {
      this.copied = false;
    }, 2000);
  }

  render() {
    return html`
      <div class="mcp-layout">
        ${this.toolGroups && this.toolGroups.length > 0
          ? this.toolGroups.map(
              (grp) => html`
                <div class="group-card">
                  <div class="group-header">
                    <div class="group-title">${grp.name}</div>
                    <div class="group-desc">${grp.description}</div>
                  </div>
                  <table>
                    <thead>
                      <tr>
                        <th>Tool Name</th>
                        <th>Access</th>
                        <th>Description</th>
                      </tr>
                    </thead>
                    <tbody>
                      ${grp.tools.map(
                        (t) => html`
                          <tr>
                            <td><span class="tool-name">${t.name}</span></td>
                            <td>
                              ${t.is_read_only
                                ? html`<span class="badge-ro">🔒 Read-Only</span>`
                                : html`<span class="badge-mut">⚡ Mutating</span>`}
                            </td>
                            <td>${t.description}</td>
                          </tr>
                        `
                      )}
                    </tbody>
                  </table>
                </div>
              `
            )
          : html`<p style="color: #64748b; text-align: center;">No clustered toolsets available.</p>`}

        <div class="config-card">
          <div class="config-header">
            <div class="config-title">Quickstart MCP Client Configuration</div>
            <div class="client-tabs">
              <button
                class="tab-btn ${this.activeClient === 'opencode' ? 'active' : ''}"
                @click="${() => (this.activeClient = 'opencode')}"
              >
                OpenCode
              </button>
              <button
                class="tab-btn ${this.activeClient === 'claude' ? 'active' : ''}"
                @click="${() => (this.activeClient = 'claude')}"
              >
                Claude Desktop
              </button>
              <button
                class="tab-btn ${this.activeClient === 'cursor' ? 'active' : ''}"
                @click="${() => (this.activeClient = 'cursor')}"
              >
                Cursor
              </button>
            </div>
          </div>

          <div class="code-container">
            <button class="copy-btn" @click="${this.copyConfig}">
              ${this.copied ? "✓ Copied" : "Copy"}
            </button>
            <pre>${this.getConfigSnippet()}</pre>
          </div>
        </div>
      </div>
    `;
  }
}
