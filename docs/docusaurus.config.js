// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/palenight');
const math = require('remark-math');
const katex = require('rehype-katex');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Formance Developer Docs',
  tagline: 'The open source foundation you need to build and scale money-movements within your app',
  url: 'https://docs.formance.com/',
  baseUrl: '/',
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/f-shape.ico',
  organizationName: 'formancehq', // Usually your GitHub org/user name.
  projectName: 'docs', // Usually your repo name.

  stylesheets: [
    'https://fonts.googleapis.com/icon?family=Material+Icons',
    'https://fonts.googleapis.com/css2?family=Roboto+Mono:wght@200;300;400;500;600&display=swap',
    'https://fonts.googleapis.com/css2?family=Inter:wght@200;300;400;500;600&display=swap',
    {
      href: 'https://cdn.jsdelivr.net/npm/katex@0.13.24/dist/katex.min.css',
      type: 'text/css',
      integrity:
        'sha384-odtC+0UGzzFL/6PNoE8rX/SPcQDXBJ+uRepguP4QkPCm2LBxH3FA3y+fKSiJ+AmM',
      crossorigin: 'anonymous',
    },
  ],

  plugins: [
    // 'posthog-docusaurus',
  ],

  presets: [
    [
      '@docusaurus/preset-classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        sitemap: {
          filename: 'sitemap.xml',
        },
        docs: {
          routeBasePath: '/',
          sidebarPath: require.resolve('./sidebars.js'),
          remarkPlugins: [math],
          rehypePlugins: [katex],
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        blog: false,
        pages: false,
      }),
    ],
    [
      'redocusaurus',
      {
        debug: Boolean(process.env.DEBUG || process.env.CI),
        specs: [
          {
            spec: './openapi/v1.json',
//            spec: './../openapi/build/generate.json',
            route: '/api/stack/v1.0',
            id: 'stack',
          }
      ],
      }
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      docs: {
        sidebar: {
          autoCollapseCategories: true,
        }
      },
      colorMode: {
        defaultMode: 'light',
        disableSwitch: true,
      },
      navbar: {
        // style: 'light',
        logo: {
          alt: 'Formance Logo',
          src: 'img/logo.svg',
          href: '/',
        },
        items: [
          {
            label: '⚡️ API Reference',
            position: 'right',
            items: [
              {
                label: 'v1.0',
                to: '/api/stack/v1.0',
              }
            ],
          },
          {
            label: 'Use-cases Library',
            position: 'right',
            href: 'https://www.formance.com/use-cases',
          },
          {
            href: 'https://github.com/formancehq/stack',
            label: 'GitHub',
            position: 'right',
          },
          {
            label: 'Go to Website',
            position: 'right',
            href: 'https://www.formance.com/',
          },
        ],
      },
      footer: {
        style: 'light',
        links: [
          {
            title: 'Documentation',
            items: [
              {
                label: 'Ledger',
                to: '/ledger',
              },
              {
                label: 'Payments',
                to: '/payments',
              },
              {
                label: 'Wallets',
                to: '/wallets',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Slack',
                href: 'https://bit.ly/formance-slack',
              },
              {
                label: 'Twitter',
                href: 'https://twitter.com/formancehq',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/formancehq/stack',
              },
              {
                label: 'Cloud Status',
                href: 'https://status.formance.com',
              },
            ],
          },
        ],
        copyright: `Copyright © 2021-2023 Numaire, Inc`,
      },
      prism: {
        theme: darkCodeTheme,
      },
      posthog: {
        apiKey: 'phc_hRDv01yOHJNUM7l5SmXPUtSQUuNw4r5am9FtV83Z9om',
        appUrl: 'https://app.posthog.com',  // optional
        enableInDevelopment: false,  // optional
      },
      algolia: {
        appId: 'IHGRMFJIIG',
        apiKey: '7864304f16ea5f9d27b7a553c83ad17a',
        indexName: 'numary',

        // Optional: see doc section below
        contextualSearch: true,

        // Optional: Specify domains where the navigation should occur through window.location instead on history.push. Useful when our Algolia config crawls multiple documentation sites and we want to navigate with window.location.href to them.
        externalUrlRegex: 'docs\\.formance\\.com',

        // Optional: Algolia search parameters
        searchParameters: {},

        // Optional: path for search page that enabled by default (`false` to disable it)
        searchPagePath: 'search',

        //... other Algolia params
      },
    })
};

module.exports = config;
