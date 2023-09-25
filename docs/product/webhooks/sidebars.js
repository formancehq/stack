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
        label: 'Webhooks',
        type: 'category',
        collapsible: true,
        collapsed: true,
        items: [
          {
            type: 'doc',
            id: 'index',
            label: 'Introduction',
          },
        ],
      }
    ]
};

module.exports = sidebars;
