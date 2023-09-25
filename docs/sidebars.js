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
        },
        {
          type: 'doc',
          id: 'getting-started/invite',
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
          type: 'link',
          label: 'Ledger',
          href: "/ledger",
        },
        {
          type: 'link',
          label: 'Payments',
          href: "/payments",
        },
        {
          type: 'link',
          label: 'Wallets (beta)',
          href: "/wallets",
        },
        {
          type: 'link',
          label: 'Flows (beta)',
          href: "/flows",
        },
        {
          type: 'link',
          label: 'Operator',
          href: "/operator",
        },
      ],
    },
    {
      type: 'category',
      label: 'Deployment',
      link: {
        type: 'doc',
        id: 'deployment/introduction',
      },
      items: [
        {
          type: 'category',
          label: 'CloudPrem',
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
        {
          type: 'category',
          label: 'Elements',
          link: {
            type: 'doc',
            id: 'deployment/elements/intro',
          },
          items: [
            {
              type: 'doc',
              id: 'deployment/elements/kubernetes',
            }
          ],
        },
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
      ],
    },
    {
      type: 'doc',
      id: 'help',
      label: 'Getting Help',
    },
  ],
};

module.exports = sidebars;
