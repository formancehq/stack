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
        label: 'Flows (beta)',
        type: 'category',
        collapsible: true,
        collapsed: true,
        items: [
          {
            type: 'doc',
            id: 'index',
            label: 'Introduction',
          },
          {
            type: 'doc',
            id: 'definition',
            label: 'Workflows definition',
          },
          {
            type: 'doc',
            id: 'execution',
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
                id: 'stages/send',
                label: 'Send',
              },
              {
                type: 'doc',
                id: 'stages/wait-event',
                label: 'Waiting for events',
              },
              {
                type: 'doc',
                id: 'stages/wait-delay',
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
                id: 'examples/ledger-to-ledger',
                label: 'Ledger to Ledger',
              },
              {
                type: 'doc',
                id: 'examples/payment-to-wallet',
                label: 'Payment to Wallet',
              },
              {
                type: 'doc',
                id: 'examples/stripe-payout',
                label: 'Ledger to Payout',
              }
            ],
          }
        ],
      },
    ]
};

module.exports = sidebars;
