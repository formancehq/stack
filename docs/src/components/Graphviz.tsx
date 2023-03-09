import { Box } from '@mui/material';
import Graphviz from 'graphviz-react';
import React from 'react';

const PostingsGraph = ({ postings, caption }) => {
  const graph : string[] = [];

  postings.map((posting) => {
    const splitSource = posting.source.split(':');
    const splitDest = posting.destination.split(':');
    let source = posting.source;
    let destination = posting.destination;

    // Dot language is escaping semicolon from label.
    // To use semicolon as label we need to use dot special char <> encoding
    if (splitSource.length > 0) {
      source = `<${posting.source}>`;
    }
    if (splitDest.length > 0) {
      destination = `<${posting.destination}>`;
    }

    const color = posting.color || 'black';
    const style = posting.style || 'solid';

    graph.push(`${source} -> ${destination} [
      label="${posting.asset} ${posting.amount}",
      weight="${posting.amount}",
      color="${color}",
      style="${style}",
    ];`);
  });

  const dot = `digraph {\nrankdir=LR\n${graph.join('\n')}\n}`;

  return (
    <Box sx={{
      textAlign: 'center',
      border: '1px solid',
      borderRadius: 2,
      borderColor: 'grey.200',
      mt: 2,
      mb: 2,
    }}>
      <Graphviz
        className="Graph"
        options={{
          width: 700,
          height: 300,
          fit: true,
          useWorker: false,
          useSharedWorker: false,
        }}
        dot={dot}/>
      <Box sx={{
        textAlign: 'center',
        fontSize: 12,
        fontStyle: 'italic',
        color: 'grey.500',
        mt: 1,
        mb: 1,
        pl: 2,
      }}>
        {caption}
      </Box>
    </Box>
  );
};

export {
  PostingsGraph,
};
