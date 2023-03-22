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
      label: 'Introduction',
      type: 'doc',
      id: 'index'
    },
    {
      label: 'Getting Started',
      type: 'category',
      collapsible: true,
      collapsed: true,
      items: [
        {
          type: 'doc',
          id: 'stack/tutorials/installation',
          customProps: {
            // icon: '💾',
          }
        },
        {
          type: 'doc',
          id: 'help',
          label: 'Getting Help',
        }
      ],
    },
    {
      label: 'Services Reference',
      type: 'category',
      collapsible: true,
      collapsed: false,
      items: [
        {
          label: 'Ledger',
          type: 'category',
          collapsible: true,
          collapsed: true,
          items: [
            {
              type: 'doc',
              id: 'ledger/index',
              label: 'Introduction'
            },
            {
              label: 'Tutorials',
              type: 'category',
              collapsible: true,
              collapsed: true,
              link: { type: 'doc', id: 'ledger/get-started/index' },
              items: [
                {
                  type: 'doc',
                  id: 'ledger/get-started/installation',
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
                  link: { type: 'doc', id: 'ledger/get-started/hello-world/index' },
                  items: [
                    // 'ledger/get-started/hello-world/start-the-server',
                    'ledger/get-started/hello-world/introducing-money',
                    'ledger/get-started/hello-world/checking-balances',
                    'ledger/get-started/hello-world/your-first-transaction'
                  ]
                },
                {
                  label: 'Mastering Numscript',
                  type: 'category',
                  collapsible: true,
                  collapsed: true,
                  link: { type: 'doc', id: 'ledger/numscript/index' },
                  customProps: {
                    // icon: '🔢',
                    description: 'Get started by creating your first transaction.',
                  },
                  items: [
                    {
                        type:'doc',
                        id: 'ledger/numscript/prerequisites',
                        customProps: {
                            // icon: '1️⃣',
                        },
                    },
                    {
                        type:'doc',
                        id: 'ledger/numscript/multi-destination/index',
                        customProps: {
                            // icon: '➗',
                        },
                    },
                    {
                        type:'doc',
                        id: 'ledger/numscript/multi-source/index',
                        customProps: {
                          // icon: '✖️',
                        },
                    },
                    {
                        type:'doc',
                        id: 'ledger/numscript/http/index',
                        customProps: {
                          // icon: '🕸',
                        },
                    },
                    {
                        type:'doc',
                        id: 'ledger/numscript/templates/index',
                        customProps: {
                          // icon: '📝',
                        },
                    },
                    {
                        type:'doc',
                        id: 'ledger/numscript/currencies/index',
                        customProps: {
                          // icon: '💴',
                        },
                    },
                    {
                        type:'doc',
                        id: 'ledger/numscript/metadata/index',
                        customProps: {
                          // icon: '📌',
                        },
                    },
                    {
                        type:'doc',
                        id: 'ledger/numscript/kept/index',
                        customProps: {
                          // icon: '🍕',
                        },
                    },
                    {
                      type:'doc',
                      id: 'ledger/numscript/overdraft/index',
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
                  id: 'ledger/advanced/publisher',
                  label: 'Publishing to HTTP / Kafka'
                },
                {
                  type: 'doc',
                  id: 'ledger/advanced/asset-conversion',
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
                'ledger/operations/installation',
                'ledger/operations/configuration',
                'ledger/operations/env-vars',
                'ledger/operations/storages',
                'ledger/operations/upgrade',
                'ledger/operations/authentication',
                'ledger/operations/using-the-control-dashboard',
                // 'operations/running-in-production',
                'ledger/api/sdks'
              ],
            },
            {
              label: 'Reference',
              type: 'category',
              collapsible: true,
              collapsed: true,
              items: [
                'ledger/reference/ledgers',
                'ledger/reference/accounts',
                'ledger/reference/transactions',
                'ledger/reference/architecture',
                'ledger/reference/concurrency-model',
                {
                  label: 'Numscript',
                  type: 'category',
                  collapsible: true,
                  collapsed: true,
                  items: [
                    'ledger/reference/numscript/machine',
                    'ledger/reference/numscript/sources',
                    'ledger/reference/numscript/destinations',
                    'ledger/reference/numscript/variables',
                    'ledger/reference/numscript/metadata',
                    'ledger/reference/numscript/rounding',
                  ],
                },
              ],
            },
            // {
            //   label: 'Examples',
            //   type: 'category',
            //   collapsible: true,
            //   collapsed: true,
            //   items: [
            //     'ledger/examples/in-app-currency',
            //     'ledger/examples/marketplace-sales-routing',
            //   ],
            // },
          ],
        },
        {
          label: 'Payments',
          type: 'category',
          collapsible: true,
          collapsed: true,
          items: [
            {
              type: 'doc',
              id: 'payments/index',
              label: 'Introduction'
            },
            {
              type: 'category',
              label: 'Available Connectors',
              collapsible: true,
              collapsed: true,
              link: {
                type: 'doc',
                id: 'payments/connectors/index',
                // label: 'Connectors'
              },
              items: [
                {
                  type: 'doc',
                  id: 'payments/connectors/stripe',
                  label: 'Stripe',
                },
                {
                  type: 'doc',
                  id: 'payments/connectors/modulr',
                  label: 'Modulr',
                },
                {
                  type: 'doc',
                  id: 'payments/connectors/wise',
                  label: 'Wise',
                },
                {
                  type: 'doc',
                  id: 'payments/connectors/currencycloud',
                  label: 'CurrencyCloud',
                },
                {
                  type: 'doc',
                  id: 'payments/connectors/bankingcircle',
                  label: 'BankingCircle',
                }
              ],
            },
            {
              type: 'doc',
              id: 'payments/connectors/framework',
              label: 'Framework',
            },
          ]
        },
        {
          label: 'Wallets (beta)',
          type: 'category',
          collapsible: true,
          collapsed: true,
          items: [
            {
              type: 'doc',
              id: 'wallets/index',
              label: 'Introduction'
            },
            // {
            //   type: 'doc',
            //   id: 'wallets/model',
            // },
            {
              type: 'doc',
              id: 'wallets/creation',
              label: 'Creating wallets',
            },
            {
              type: 'category',
              label: 'Managing funds',
              collapsible: true,
              collapsed: true,
              items: [
                {
                  type: 'doc',
                  id: 'wallets/funds/adding-funds',
                  label: 'Adding funds',
                },
                {
                  type: 'doc',
                  id: 'wallets/funds/spending-funds',
                  label: 'Spending funds',
                },
                {
                  type: 'doc',
                  id: 'wallets/funds/holds',
                  label: 'Hold and confirm',
                },
                // {
                //   type: 'doc',
                //   id: 'wallets/funds/expirable-funds',
                //   label: 'Expirable funds',
                // },
                // {
                //   type: 'doc',
                //   id: 'wallets/funds/funds-reserve',
                //   label: 'Reserved funds',
                // }
              ],
            },
            {
              type: 'doc',
              id: 'wallets/configuration',
              label: 'Configuration',
            }
          ]
        },
      ]
    },
    {
      label: 'Stack Reference',
      type: 'category',
      collapsible: true,
      collapsed: true,
      items: [
        {
          type: 'doc',
          id: 'stack/architecture',
          label: 'Architecture',
        },
        {
          type: 'doc',
          id: 'stack/authentication/index',
          label: 'Authentication',
        },
        {
          type: 'doc',
          id: 'stack/sdk/index',
          label: 'SDKs',
        },
        {
          type: 'doc',
          id: 'stack/webhooks/index',
          label: 'Webhooks',
        },
        {
          type: 'doc',
          id: 'stack/unambiguous-monetary-notation',
          label: 'Monetary Notation (UMN)',
        },
        {
          type: 'category',
          label: 'Self-hosting',
          items: [
            {
              type: 'doc',
              id: 'stack/reference/docker',
              label: 'Docker',
            },
            {
              type: 'doc',
              id: 'stack/reference/helm',
              label: 'Kubernetes / Helm',
            },
            // {
            //   type: 'doc',
            //   id: 'stack/reference/production',
            //   label: 'Production checklist',
            // }
          ],
        },
        // {
        //   type: 'doc',
        //   id: 'stack/telemetry/index',
        //   label: 'Telemetry',
        // },
      ],
    },
    {
      label: 'Cloud',
      type: 'category',
      collapsible: true,
      collapsed: true,
      items: [
        {
          type: 'doc',
          id: 'stack/fctl',
          label: 'Formance CLI',
        },
        {
          type: 'doc',
          id: 'cloud/sandboxes',
          label: 'Sandbox environment',
        },
        {
          type: 'doc',
          id: 'cloud/production',
          label: 'Production environment',
        },
        {
          type: 'doc',
          id: 'cloud/authentication',
          label: 'Authentication',
        },
        // {
        //   type: 'category',
        //   label: 'Conventions',
        //   items: [
        //     {
        //       type: 'doc',
        //       id: 'stack/unambiguous-monetary-notation',
        //       label: 'Monetary Notation',
        //     }
        //   ],
        // },
      ],
    },
  ],
};

module.exports = sidebars;
