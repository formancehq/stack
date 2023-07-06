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
          id: 'getting-started/fctl-quick-start',
          customProps: {
            // icon: '💾',
          }
        },
        {
          type: 'doc',
          id: 'guides/newSandbox',
        }
      ],
    },
    {
      label: 'Products',
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
        {
          label: 'Flows (beta)',
          type: 'category',
          collapsible: true,
          collapsed: true,
          items: [
            {
              type: 'doc',
              id: 'flows/index',
              label: 'Introduction',
            },
            {
              type: 'doc',
              id: 'flows/definition',
              label: 'Workflows definition',
            },
            {
              type: 'doc',
              id: 'flows/execution',
              label: 'Workflows execution',
            },
            {
              type: 'category',
              label: 'Stages reference',
              collapsible: true,
              collapsed: false,
              items: [
                {
                  type: 'doc',
                  id: 'flows/stages/send',
                  label: 'Send',
                },
                {
                  type: 'doc',
                  id: 'flows/stages/wait-event',
                  label: 'Waiting for events',
                },
                {
                  type: 'doc',
                  id: 'flows/stages/wait-delay',
                  label: 'Waiting for a delay',
                }
              ],
            },
            {
              type: 'category',
              label: 'Examples',
              items: [
                {
                  type: 'doc',
                  id: 'flows/examples/ledger-to-ledger',
                  label: 'Ledger to Ledger',
                },
                {
                  type: 'doc',
                  id: 'flows/examples/payment-to-wallet',
                  label: 'Payment to Wallet',
                },
                {
                  type: 'doc',
                  id: 'flows/examples/stripe-payout',
                  label: 'Ledger to Payout',
                }
              ],
            }
          ],
        }
      ]
    },
    {
      type: 'category',
      label: 'Deployment',
      link: {
        type: 'doc',
        id: 'deployment/introduction',
      },
      items: [
        // {
        //   type: 'category',
        //   label: 'Formance Cloud',
        //   link: {
        //     type: 'doc',
        //     id: 'deployment/cloud/intro',
        //   },
        //   items: [
        //     {
        //       type: 'doc',
        //       id: 'deployment/cloud/regions',
        //       label: 'Regions',
        //     }
        //   ],
        // },
        {
          type: 'category',
          label: 'Formance CloudPrem',
          link: { type: 'doc', id: 'deployment/cloudprem/intro' },
          items: [
            {
              type: 'doc',
              id: 'deployment/cloudprem/cluster-config',
            },
            {
              type: 'doc',
              id: 'deployment/cloudprem/private-regions',
            },
            {
              type: 'doc',
              id: 'deployment/cloudprem/operator',
            },
            {
              type: 'doc',
              id: 'deployment/cloudprem/usage',
            },
          ],
        },
        // {
        //   type: 'category',
        //   label: 'Elements',
        //   link: {
        //     type: 'doc',
        //     id: 'deployment/elements/intro',
        //   },
        //   items: [
        //     {
        //       type: 'doc',
        //       id: 'deployment/elements/docker',
        //       label: 'Docker',
        //     },
        //     {
        //       type: 'doc',
        //       id: 'deployment/elements/kubernetes',
        //     }
        //   ],
        // },
      ],
    },
    {
      type: 'category',
      label: 'SDKs',
      items: [
        {
          type: 'doc',
          id: 'stack/sdk/index',
        },
      ],
    },
    {
      label: 'Resources',
      type: 'category',
      collapsible: true,
      collapsed: true,
      items: [
        {
          type: 'doc',
          id: 'stack/architecture',
        },
        {
          type: 'doc',
          id: 'stack/authentication/index',
          label: 'Authentication',
        },
        {
          type: 'doc',
          id: 'stack/unambiguous-monetary-notation',
          label: 'Monetary Notation (UMN)',
        }
//        // {
        //   type: 'doc',
        //   id: 'stack/telemetry/index',
        //   label: 'Telemetry',
        // },
      ],
    },
    {
      type: 'doc',
      id: 'help',
      label: 'Getting Help',
    },
    // {
    //   label: 'Guides',
    //   type: 'category',
    //   collapsible: true,
    //   collapsed: false,
    //   items: [

    //   ]
    // },
  ],
};

module.exports = sidebars;
