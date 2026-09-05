import { LitElement, html, css } from "lit";
import { customElement, state } from "lit/decorators.js";

interface StepConfig {
  num: number;
  label: string;
  message: string;
  progressPercent: number;
}

const STEPS: StepConfig[] = [
  {
    num: 1,
    label: "Ingest & Hygiene",
    message: "Fetching repository tree and applying hygiene filters...",
    progressPercent: 12,
  },
  {
    num: 2,
    label: "Scan APIs & Trees",
    message: "Extracting REST route decorators, Cobra command trees, and authored prompts...",
    progressPercent: 38,
  },
  {
    num: 3,
    label: "5D Scorecard",
    message: "Evaluating 5-dimension readiness scorecard and execution safety boundaries...",
    progressPercent: 65,
  },
  {
    num: 4,
    label: "Cluster Personas",
    message: "Synthesizing multi-skill blueprints and clustering MCP toolsets...",
    progressPercent: 92,
  },
];

@customElement("montage-pipeline")
export class MontagePipeline extends LitElement {
  @state() private currentStep = 1;
  @state() private statusMessage = STEPS[0].message;
  @state() private progressWidth = STEPS[0].progressPercent;

  private timers: number[] = [];

  static styles = css`
    :host {
      display: block;
      width: 100%;
      max-width: 820px;
      margin: 2rem auto;
    }

    .pipeline-card {
      background: #0f172a;
      border: 1px solid #1e293b;
      border-radius: 14px;
      padding: 1.75rem 2rem 1.5rem;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4);
    }

    .steps-container {
      display: flex;
      align-items: center;
      justify-content: space-between;
      position: relative;
      margin-bottom: 1.5rem;
    }

    .step-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.6rem;
      position: relative;
      z-index: 2;
    }

    /* Solid opaque background prevents the line behind from showing through */
    .step-circle {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.875rem;
      font-weight: 700;
      font-family: 'JetBrains Mono', monospace;
      background: #0f172a;
      border: 2px solid #334155;
      color: #64748b;
      transition: all 0.4s ease;
      position: relative;
      z-index: 2;
    }

    /* Completed Step */
    .step-completed .step-circle {
      background: #10b981;
      border-color: #10b981;
      color: #ffffff;
      box-shadow: 0 0 10px rgba(16, 185, 129, 0.3);
    }

    /* Active Step */
    .step-active .step-circle {
      background: #0f172a;
      border-color: #10b981;
      color: #34d399;
      box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.18), 0 0 16px rgba(16, 185, 129, 0.35);
      animation: pulseAura 2s infinite ease-in-out;
    }

    @keyframes pulseAura {
      0%, 100% {
        box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.18), 0 0 14px rgba(16, 185, 129, 0.3);
      }
      50% {
        box-shadow: 0 0 0 6px rgba(6, 182, 212, 0.25), 0 0 20px rgba(6, 182, 212, 0.45);
      }
    }

    .step-label {
      font-size: 0.8125rem;
      font-weight: 600;
      color: #64748b;
      text-align: center;
      white-space: nowrap;
      transition: color 0.3s ease;
    }

    .step-completed .step-label {
      color: #94a3b8;
    }

    .step-active .step-label {
      color: #38bdf8;
      font-weight: 700;
    }

    /* Connecting Track Line */
    .step-line {
      position: absolute;
      top: 20px;
      left: 20px;
      right: 20px;
      height: 3px;
      background: #1e293b;
      z-index: 1;
      border-radius: 2px;
      overflow: hidden;
    }

    .step-line-progress {
      height: 100%;
      background: linear-gradient(90deg, #10b981, #06b6d4);
      transition: width 0.6s cubic-bezier(0.4, 0, 0.2, 1);
    }

    /* Status Ticker */
    .status-ticker {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.65rem;
      padding: 0.75rem 1rem;
      background: rgba(30, 41, 59, 0.6);
      border: 1px solid #1e293b;
      border-radius: 8px;
    }

    .ticker-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: #10b981;
      box-shadow: 0 0 8px #10b981;
      animation: blink 1.2s infinite ease-in-out;
      flex-shrink: 0;
    }

    @keyframes blink {
      0%, 100% { opacity: 0.3; transform: scale(0.85); }
      50% { opacity: 1; transform: scale(1.15); }
    }

    .ticker-text {
      font-size: 0.825rem;
      font-family: 'JetBrains Mono', monospace;
      color: #cbd5e1;
      text-align: center;
    }
  `;

  connectedCallback() {
    super.connectedCallback();
    this.startProgress();
  }

  disconnectedCallback() {
    super.disconnectedCallback();
    this.clearTimers();
  }

  private clearTimers() {
    this.timers.forEach((t) => window.clearTimeout(t));
    this.timers = [];
  }

  private startProgress() {
    this.clearTimers();
    this.currentStep = 1;
    this.statusMessage = STEPS[0].message;
    this.progressWidth = STEPS[0].progressPercent;

    // Advance to Step 2 at 1.5s
    this.timers.push(
      window.setTimeout(() => {
        this.currentStep = 2;
        this.statusMessage = STEPS[1].message;
        this.progressWidth = STEPS[1].progressPercent;
      }, 1500)
    );

    // Advance to Step 3 at 4.0s
    this.timers.push(
      window.setTimeout(() => {
        this.currentStep = 3;
        this.statusMessage = STEPS[2].message;
        this.progressWidth = STEPS[2].progressPercent;
      }, 4000)
    );

    // Advance to Step 4 at 7.0s
    this.timers.push(
      window.setTimeout(() => {
        this.currentStep = 4;
        this.statusMessage = STEPS[3].message;
        this.progressWidth = STEPS[3].progressPercent;
      }, 7000)
    );
  }

  render() {
    return html`
      <div class="pipeline-card">
        <div class="steps-container">
          <div class="step-line">
            <div class="step-line-progress" style="width: ${this.progressWidth}%"></div>
          </div>

          ${STEPS.map((s) => {
            const isCompleted = s.num < this.currentStep;
            const isActive = s.num === this.currentStep;
            const stateClass = isCompleted ? "step-completed" : isActive ? "step-active" : "step-pending";

            return html`
              <div class="step-item ${stateClass}">
                <div class="step-circle">
                  ${isCompleted
                    ? html`
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <path d="M5 13l4 4L19 7" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                      `
                    : isActive
                    ? html`
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" style="animation: spin 1.2s linear infinite;">
                          <style>@keyframes spin { 100% { transform: rotate(360deg); } }</style>
                          <path d="M21 12a9 9 0 11-6.219-8.56"/>
                        </svg>
                      `
                    : s.num}
                </div>
                <span class="step-label">${s.label}</span>
              </div>
            `;
          })}
        </div>

        <div class="status-ticker">
          <div class="ticker-dot"></div>
          <div class="ticker-text">${this.statusMessage}</div>
        </div>
      </div>
    `;
  }
}
