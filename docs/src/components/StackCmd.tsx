import React from 'react';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

export function StackCmd({ children }) {
  if (children.length == 1) {
    return <code>{children}</code>;
  }

  return (
    <Tabs>
      <TabItem value="fctl" label="fctl">
        {children[0]}
      </TabItem>
      <TabItem value="api" label="curl">
        {children[1]}
      </TabItem>
    </Tabs>
  );
}
