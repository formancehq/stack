/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

const { Collapse } = require('@material-ui/core');

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  main: [
        {
          label: 'Ledger',
          type: 'category',
          collapsible: false,
          collapsed: false,
          items: [
            {
              type: 'doc',
              id: 'index',
              label: 'Introduction',
            },
            {
              label: 'Tutorials',
              type: 'category',
              collapsible: true,
              collapsed: true,
              link: { type: 'doc', id: 'get-started/index' },
              items: [
                {
                  type: 'doc',
                  id: 'get-started/installation',
                  customProps: {
                    // icon: '💾',
                  }
                },
                {
                  label: 'Hello World',
                  type: 'category',
                  collapsible: true,
                  collapsed: true,
                  customProps: {
                    // icon: '👋🏾',
                    description: 'Get started by creating your first transaction.',
                  },
                  link: { type: 'doc', id: 'get-started/hello-world/index' },
                  items: [
                    // 'get-started/hello-world/start-the-server',
                    'get-started/hello-world/introducing-money',
                    'get-started/hello-world/checking-balances',
                    'get-started/hello-world/your-first-transaction'
                  ]
                },
                {
                  label: 'Mastering Numscript',
                  type: 'category',
                  collapsible: true,
                  collapsed: true,
                  link: { type: 'doc', id: 'numscript/index' },
                  customProps: {
                    // icon: '🔢',
                    description: 'Get started by creating your first transaction.',
                  },
                  items: [
                    {
                        type:'doc',
                        id: 'numscript/prerequisites',
                        customProps: {
                            // icon: '1️⃣',
                        },
                    },
                    {
                        type:'doc',
                        id: 'numscript/multi-destination/index',
                        customProps: {
                            // icon: '➗',
                        },
                    },
                    {
                        type:'doc',
                        id: 'numscript/multi-source/index',
                        customProps: {
                          // icon: '✖️',
                        },
                    },
                    {
                        type:'doc',
                        id: 'numscript/http/index',
                        customProps: {
                          // icon: '🕸',
                        },
                    },
                    {
                        type:'doc',
                        id: 'numscript/templates/index',
                        customProps: {
                          // icon: '📝',
                        },
                    },
                    {
                        type:'doc',
                        id: 'numscript/currencies/index',
                        customProps: {
                          // icon: '💴',
                        },
                    },
                    {
                        type:'doc',
                        id: 'numscript/metadata/index',
                        customProps: {
                          // icon: '📌',
                        },
                    },
                    {
                        type:'doc',
                        id: 'numscript/kept/index',
                        customProps: {
                          // icon: '🍕',
                        },
                    },
                    {
                      type:'doc',
                      id: 'numscript/overdraft/index',
                      customProps: {
                        // icon: '📌',
                      },
                  },
                  ],
                }
              ],
            },
            {
              label: 'Guides',
              type: 'category',
              collapsible: true,
              collapsed: true,
              items: [
                {
                  type: 'doc',
                  id: 'advanced/publisher',
                  label: 'Publishing to HTTP / Kafka'
                },
                {
                  type: 'doc',
                  id: 'advanced/asset-conversion',
                  label: 'Currency conversion',
                },
              ],
            },
            {
              label: 'Deployment',
              type: 'category',
              collapsible: true,
              collapsed: true,
              items: [
                'operations/installation',
                'operations/configuration',
                'operations/env-vars',
                'operations/storages',
                'operations/upgrade',
                'operations/authentication',
                'operations/using-the-control-dashboard',
                'api/sdks'
              ],
            },
            {
              label: 'Reference',
              type: 'category',
              collapsible: true,
              collapsed: true,
              items: [
                'reference/ledgers',
                'reference/accounts',
                'reference/transactions',
                'reference/architecture',
                'reference/concurrency-model',
                {
                  label: 'Numscript',
                  type: 'category',
                  collapsible: true,
                  collapsed: true,
                  items: [
                    'reference/numscript/machine',
                    'reference/numscript/postings',
                    'reference/numscript/sources',
                    'reference/numscript/destinations',
                    'reference/numscript/variables',
                    'reference/numscript/metadata',
                    'reference/numscript/rounding',
                  ],
                },
              ],
            },
          ],
      }
    ]
};

module.exports = sidebars;
