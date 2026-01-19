import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'ftsctl',
  description: 'Command-line tool for FTSnext',
  lang: 'en-US',
  base: '/fts-control/',

  themeConfig: {
    siteTitle: 'ftsctl',

    nav: [
      { text: 'Getting Started', link: '/getting-started/installation' },
      { text: 'GitHub', link: 'https://github.com/medizininformatik-initiative/fts-control' }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Installation', link: '/getting-started/installation' },
          { text: 'Configuration', link: '/getting-started/configuration' },
          { text: 'Commands', link: '/getting-started/commands' }
        ]
      }
    ],

    search: {
      provider: 'local'
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/medizininformatik-initiative/fts-control' }
    ]
  }
})
