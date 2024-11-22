import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "idpzero",
  description: "Developer first Identity Provider",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Getting Started', link: '/guide/what-is-it' }
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'What is idpzero?', link: '/guide/what-is-it' },
          { text: 'Getting Started', link: '/guide/getting-started' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/idpzero/idpzero' }
    ]
  }
})
