import React from 'react';
import { NumscriptBlock } from 'react-numscript-codeblock';

export default function Numscript({ script }) {
  return (
    <NumscriptBlock script={script} callback></NumscriptBlock>
  );
}
