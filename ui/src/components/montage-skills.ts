import { LitElement, html, css } from "lit";
import { customElement, property } from "lit/decorators.js";
import { SkillBlueprint } from "../api/types";

@customElement("montage-skills")
export class MontageSkills extends LitElement {
  @property({ type: Array }) skills: SkillBlueprint[] = [];

  static styles = css`
    :host {
      display: block;
      width: 100%;
    }

    .skills-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
      gap: 1.5rem;
    }

    .skill-card {
      background: #0f172a;
      border: 1px solid #1e293b;
      border-radius: 14px;
      padding: 1.5rem;
      display: flex;
      flex-direction: column;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
      transition: all 0.2s ease;
    }

    .skill-card:hover {
      border-color: #334155;
      transform: translateY(-2px);
    }

    .skill-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 0.75rem;
      margin-bottom: 0.75rem;
    }

    .skill-name {
      font-size: 1.05rem;
      font-weight: 700;
      color: #f8fafc;
      font-family: 'JetBrains Mono', monospace;
    }

    .kind-badge {
      font-size: 0.725rem;
      font-weight: 700;
      padding: 0.2rem 0.55rem;
      border-radius: 6px;
      white-space: nowrap;
    }

    .kind-pure {
      background: rgba(147, 51, 234, 0.15);
      border: 1px solid rgba(147, 51, 234, 0.35);
      color: #c084fc;
    }

    .kind-code {
      background: rgba(245, 158, 11, 0.15);
      border: 1px solid rgba(245, 158, 11, 0.35);
      color: #fbbf24;
    }

    .skill-domain {
      font-size: 0.8rem;
      color: #38bdf8;
      font-weight: 600;
      margin-bottom: 0.6rem;
    }

    .skill-desc {
      font-size: 0.85rem;
      color: #94a3b8;
      line-height: 1.5;
      margin-bottom: 1.25rem;
      flex: 1;
    }

    .section-title {
      font-size: 0.75rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
      margin-bottom: 0.5rem;
    }

    .cap-list {
      display: flex;
      flex-wrap: wrap;
      gap: 0.4rem;
      margin-bottom: 1.25rem;
    }

    .cap-tag {
      background: #1e293b;
      color: #e2e8f0;
      font-size: 0.75rem;
      font-weight: 600;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
      border: 1px solid #334155;
    }

    .script-box {
      background: rgba(245, 158, 11, 0.08);
      border: 1px solid rgba(245, 158, 11, 0.25);
      border-radius: 8px;
      padding: 0.75rem;
      margin-bottom: 1.25rem;
    }

    .script-path {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.775rem;
      color: #fbbf24;
      display: block;
      margin-bottom: 0.25rem;
    }

    .script-savings {
      font-size: 0.75rem;
      color: #cbd5e1;
      font-weight: 600;
    }

    .instructions-list {
      list-style: none;
      display: flex;
      flex-direction: column;
      gap: 0.35rem;
    }

    .instructions-list li {
      font-size: 0.8rem;
      color: #cbd5e1;
      display: flex;
      align-items: flex-start;
      gap: 0.4rem;
    }

    .instructions-list li::before {
      content: "•";
      color: #10b981;
      font-weight: bold;
    }
  `;

  render() {
    if (!this.skills || this.skills.length === 0) {
      return html`<p style="color: #64748b; text-align: center; padding: 2rem;">No specialized skills generated.</p>`;
    }

    return html`
      <div class="skills-grid">
        ${this.skills.map((s) => {
          const isCode = s.kind === "skill_with_code";
          return html`
            <div class="skill-card">
              <div class="skill-header">
                <span class="skill-name">${s.name}</span>
                <span class="kind-badge ${isCode ? 'kind-code' : 'kind-pure'}">
                  ${isCode ? '⚡ Skill with Code' : '📝 Pure Skill'}
                </span>
              </div>

              <div class="skill-domain">Domain: ${s.domain}</div>
              <div class="skill-desc">${s.description}</div>

              ${s.helper_script
                ? html`
                    <div class="script-box">
                      <span class="script-path">📄 ${s.helper_script}</span>
                      <span class="script-savings">💡 ${s.token_savings_est}</span>
                    </div>
                  `
                : ""}

              <div class="section-title">Capabilities Included</div>
              <div class="cap-list">
                ${s.capabilities.map((c) => html`<span class="cap-tag">${c}</span>`)}
              </div>

              ${s.key_instructions && s.key_instructions.length > 0
                ? html`
                    <div class="section-title">Key Workflow Rules</div>
                    <ul class="instructions-list">
                      ${s.key_instructions.map((inst) => html`<li><span>${inst}</span></li>`)}
                    </ul>
                  `
                : ""}
            </div>
          `;
        })}
      </div>
    `;
  }
}
