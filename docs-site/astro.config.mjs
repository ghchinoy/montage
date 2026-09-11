import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import catppuccin from '@catppuccin/starlight';

// https://astro.build/config
export default defineConfig({
  site: 'https://ghchinoy.github.io',
  base: '/montage',
  integrations: [
    starlight({
      title: 'Montage',
      description: 'Agent Ecosystem Readiness Assessor & Tooling Scaffolder (CLI + MCP Server)',
      logo: {
        src: './src/assets/logo.svg',
        replacesTitle: false,
      },
      social: [
        {
          icon: 'github',
          label: 'Montage on GitHub',
          href: 'https://github.com/ghchinoy/montage',
        },
      ],
      plugins: [
        catppuccin({
          dark: { flavor: 'mocha', accent: 'lavender' },
          light: { flavor: 'latte', accent: 'lavender' },
        }),
      ],
      defaultLocale: 'root',
      sidebar: [
        {
          label: 'Overview',
          items: [
            { label: 'Introduction', slug: 'index' },
            { label: 'Installation', slug: 'getting-started/installation' },
            { label: 'Quick Start', slug: 'getting-started/quickstart' },
            { label: 'Web Dashboard', slug: 'getting-started/web-dashboard' },
          ],
        },
        {
          label: 'User Guides',
          items: [
            { label: 'Reviewing Applications', slug: 'guides/reviewing-applications' },
            { label: 'Scaffolding Tooling', slug: 'guides/scaffolding-tooling' },
            { label: 'CLI Reference', slug: 'guides/cli-reference' },
            { label: 'MCP Server', slug: 'guides/mcp-server' },
          ],
        },
        {
          label: 'Architecture & Theory',
          items: [
            { label: '5 Readiness Dimensions', slug: 'architecture/readiness-dimensions' },
            { label: 'Discovery Engine & AST', slug: 'architecture/discovery-engine' },
            { label: 'Agent Ecosystem Spectrum', slug: 'architecture/agent-ecosystem-spectrum' },
            { label: 'Token Economics & Personas', slug: 'architecture/token-economics' },
          ],
        },
        {
          label: 'Reference',
          items: [
            { label: 'Agent Plugins Spec', slug: 'reference/agent-plugins-spec' },
            { label: 'Actionable Gaps & Remediation', slug: 'reference/remediation-guide' },
          ],
        },
      ],
    }),
  ],
});
