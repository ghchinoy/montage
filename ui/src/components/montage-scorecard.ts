import { LitElement, html, css } from "lit";
import { customElement, property } from "lit/decorators.js";
import { AssessmentReport } from "../api/types";

@customElement("montage-scorecard")
export class MontageScorecard extends LitElement {
  @property({ type: Object }) report!: AssessmentReport;

  static styles = css`
    :host {
      display: block;
      width: 100%;
    }

    .scorecard-grid {
      display: grid;
      grid-template-columns: 320px 1fr;
      gap: 1.5rem;
      margin-bottom: 2rem;
    }

    .score-card {
      background: #0f172a;
      border: 1px solid #1e293b;
      border-radius: 14px;
      padding: 2rem 1.5rem;
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
    }

    .circle-container {
      position: relative;
      width: 150px;
      height: 150px;
      margin-bottom: 1.25rem;
    }

    .circle-bg {
      fill: none;
      stroke: #1e293b;
      stroke-width: 10;
    }

    .circle-progress {
      fill: none;
      stroke-width: 10;
      stroke-linecap: round;
      transform: rotate(-90deg);
      transform-origin: 50% 50%;
      transition: stroke-dashoffset 1s ease;
    }

    .score-text {
      position: absolute;
      inset: 0;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
    }

    .score-val {
      font-size: 2.75rem;
      font-weight: 800;
      line-height: 1;
      font-family: 'JetBrains Mono', monospace;
    }

    .score-max {
      font-size: 0.875rem;
      color: #64748b;
      font-weight: 600;
      margin-top: 0.25rem;
    }

    .tier-badge {
      display: inline-block;
      padding: 0.4rem 1rem;
      border-radius: 8px;
      font-size: 0.925rem;
      font-weight: 700;
      margin-bottom: 0.75rem;
    }

    .tier-native {
      background: rgba(16, 185, 129, 0.15);
      border: 1px solid rgba(16, 185, 129, 0.35);
      color: #34d399;
    }

    .tier-capable {
      background: rgba(6, 182, 212, 0.15);
      border: 1px solid rgba(6, 182, 212, 0.35);
      color: #22d3ee;
    }

    .tier-scriptable {
      background: rgba(245, 158, 11, 0.15);
      border: 1px solid rgba(245, 158, 11, 0.35);
      color: #fbbf24;
    }

    .tier-other {
      background: rgba(244, 63, 94, 0.15);
      border: 1px solid rgba(244, 63, 94, 0.35);
      color: #f87171;
    }

    .repo-submeta {
      font-size: 0.825rem;
      color: #94a3b8;
      line-height: 1.4;
    }

    .dimensions-card {
      background: #0f172a;
      border: 1px solid #1e293b;
      border-radius: 14px;
      padding: 1.75rem;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
      display: flex;
      flex-direction: column;
      gap: 1.25rem;
    }

    .dimension-row {
      display: flex;
      flex-direction: column;
      gap: 0.35rem;
    }

    .dim-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .dim-name {
      font-size: 0.95rem;
      font-weight: 700;
      color: #f8fafc;
    }

    .dim-weight {
      font-size: 0.75rem;
      color: #64748b;
      font-weight: 600;
      margin-left: 0.5rem;
    }

    .dim-score {
      font-family: 'JetBrains Mono', monospace;
      font-size: 0.95rem;
      font-weight: 700;
    }

    .dim-bar-bg {
      height: 8px;
      background: #1e293b;
      border-radius: 4px;
      overflow: hidden;
    }

    .dim-bar-fill {
      height: 100%;
      border-radius: 4px;
      transition: width 1s ease;
    }

    .dim-key-obs {
      font-size: 0.8125rem;
      color: #94a3b8;
    }

    .guidance-box {
      grid-column: 1 / -1;
      background: rgba(99, 102, 241, 0.08);
      border: 1px solid rgba(99, 102, 241, 0.25);
      border-radius: 12px;
      padding: 1.25rem 1.5rem;
      display: flex;
      align-items: flex-start;
      gap: 1rem;
    }

    .spark-icon {
      color: #818cf8;
      font-size: 1.4rem;
      line-height: 1;
      margin-top: 0.15rem;
    }

    .guidance-content h4 {
      font-size: 0.925rem;
      font-weight: 700;
      color: #c7d2fe;
      margin-bottom: 0.25rem;
    }

    .guidance-content p {
      font-size: 0.875rem;
      color: #e2e8f0;
      line-height: 1.55;
    }

    @media (max-width: 840px) {
      .scorecard-grid {
        grid-template-columns: 1fr;
      }
    }
  `;

  render() {
    if (!this.report) return html``;

    const score = this.report.overall_score;
    const tier = this.report.tier;

    // Radius = 60, circumference = 2 * PI * 60 ≈ 377
    const circumference = 377;
    const offset = circumference - (score / 100) * circumference;

    let strokeColor = "#10b981";
    let tierClass = "tier-native";
    if (score < 55) {
      strokeColor = "#f43f5e";
      tierClass = "tier-other";
    } else if (score < 75) {
      strokeColor = "#f59e0b";
      tierClass = "tier-scriptable";
    } else if (score < 90) {
      strokeColor = "#06b6d4";
      tierClass = "tier-capable";
    }

    return html`
      <div class="scorecard-grid">
        <div class="score-card">
          <div class="circle-container">
            <svg width="150" height="150" viewBox="0 0 150 150">
              <circle class="circle-bg" cx="75" cy="75" r="60" />
              <circle
                class="circle-progress"
                cx="75"
                cy="75"
                r="60"
                stroke="${strokeColor}"
                stroke-dasharray="${circumference}"
                stroke-dashoffset="${offset}"
              />
            </svg>
            <div class="score-text">
              <span class="score-val" style="color: ${strokeColor};">${score}</span>
              <span class="score-max">/ 100</span>
            </div>
          </div>

          <div class="tier-badge ${tierClass}">${tier}</div>
          <div class="repo-submeta">
            ${this.report.target.primary_language ? html`<strong>${this.report.target.primary_language}</strong> · ` : ""}
            ${this.report.target.description || "No description provided."}
          </div>
        </div>

        <div class="dimensions-card">
          ${this.report.dimensions.map((dim) => {
            const dimScore = dim.score;
            let barColor = "#10b981";
            if (dimScore < 50) barColor = "#f43f5e";
            else if (dimScore < 75) barColor = "#f59e0b";
            else if (dimScore < 90) barColor = "#06b6d4";

            const keyObs = dim.strengths?.[0] || dim.gaps?.[0] || dim.explanation;

            return html`
              <div class="dimension-row">
                <div class="dim-header">
                  <div>
                    <span class="dim-name">${dim.name}</span>
                    <span class="dim-weight">Weight: ${(dim.weight * 100).toFixed(0)}%</span>
                  </div>
                  <span class="dim-score" style="color: ${barColor}">${dimScore}/100</span>
                </div>
                <div class="dim-bar-bg">
                  <div class="dim-bar-fill" style="width: ${dimScore}%; background: ${barColor};"></div>
                </div>
                <div class="dim-key-obs">${keyObs}</div>
              </div>
            `;
          })}
        </div>

        ${this.report.tier_summary
          ? html`
              <div class="guidance-box">
                <div class="spark-icon">✨</div>
                <div class="guidance-content">
                  <h4>Agent Architecture Guidance</h4>
                  <p>${this.report.tier_summary}</p>
                </div>
              </div>
            `
          : ""}
      </div>
    `;
  }
}
