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
      label: 'Wallets (beta)',
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
          type: 'doc',
          id: 'creation',
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
              id: 'funds/adding-funds',
              label: 'Adding funds',
            },
            {
              type: 'doc',
              id: 'funds/spending-funds',
              label: 'Spending funds',
            },
            {
              type: 'doc',
              id: 'funds/holds',
              label: 'Hold and confirm',
            },
            // {
            //   type: 'doc',
            //   id: 'funds/expirable-funds',
            //   label: 'Expirable funds',
            // },
            // {
            //   type: 'doc',
            //   id: 'funds/funds-reserve',
            //   label: 'Reserved funds',
            // }
          ],
        },
      ]
    }
    ]
};

module.exports = sidebars;
