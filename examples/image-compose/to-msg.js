const fs = require('fs');
const {SendMessage, ParseMessage} = require('../flow-messaging');
const send = SendMessage(process.stdout);

// Config-parser

process.stdin.on('data', data => {
  const configStr = data.toString().trim();
  const tableRegex = RegExp(/^(\|[^\n]+\|\r?\n)((?:\|:?[-]+:?)+\|)(\n(?:\|[^\n]+\|\r?\n?)*)?$/, 'gm');
  const tableMatches = tableRegex.exec(configStr);

  if (tableMatches && tableMatches.length > 3) {
    const linesStr = tableMatches[3];
    const config = linesStr.split('\n')
      .map(line => {
        const cleanLineRegex = RegExp(/\|(.*?\|.*?)\|/, 'gm');
        const cleanLineMatches = cleanLineRegex.exec(line);
        if (cleanLineMatches && cleanLineMatches.length > 1) {
          const cleanLine = cleanLineMatches[1];
          return cleanLine.split('|').map(col => col.trim());
        }
        return null;
      })
      .filter(l => !!l);
    console.error('sending config message:', JSON.stringify(config));
    send({cmd: 'COMPOSE', config});
  } else {
    throw new Error(`Markdown table in config file "${configFilename}" could not be parsed, only strict markdown tables are supported`);
  }
});