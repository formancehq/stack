import { Box } from '@mui/material';
import React from 'react';
import { NumscriptBlock } from 'react-numscript-codeblock';

export function Numscript({ script }) {
  return (
    <Box sx={{
      mb: 2,
      border: 'solid 1px rgba(0, 0, 0, 0.12)',
      borderRadius: 1,
    }}>
      <NumscriptBlock script={script}></NumscriptBlock>
    </Box>
  );
}
