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
      label: 'Payments',
      type: 'category',
      collapsible: true,
      collapsed: true,
      items: [
        {
          type: 'doc',
          id: 'index',
          label: 'Introduction'
        },
        {
          type: 'category',
          label: 'Available Connectors',
          collapsible: true,
          collapsed: true,
          link: {
            type: 'doc',
            id: 'connectors/index',
            // label: 'Connectors'
          },
          items: [
            {
              type: 'doc',
              id: 'connectors/stripe',
              label: 'Stripe',
            },
            {
              type: 'doc',
              id: 'connectors/modulr',
              label: 'Modulr',
            },
            {
              type: 'doc',
              id: 'connectors/wise',
              label: 'Wise',
            },
            {
              type: 'doc',
              id: 'connectors/currencycloud',
              label: 'CurrencyCloud',
            },
            {
              type: 'doc',
              id: 'connectors/bankingcircle',
              label: 'BankingCircle',
            },
            {
              type: 'doc',
              id: 'connectors/mangopay',
              label: 'Mangopay',
            },
            {
              type: 'doc',
              id: 'connectors/moneycorp',
              label: 'Moneycorp',
            }
          ],
        },
        {
          type: 'doc',
          id: 'connectors/framework',
          label: 'Framework',
        },
      ]
    },
    ]
};

module.exports = sidebars;
