#!/usr/bin/env node

const propName = process.env.PROP_NAME;
if (!propName) {
  console.error('PROP_NAME missing');
  process.exit(-1);
}

const JSONEventReader = require('../../node-flowflux/json-event-reader');
const LogWriter = require('../../node-flowflux/log-writer');
const JSONEventWriter = require('../../node-flowflux/json-event-writer');

const input = new JSONEventReader(process.stdin);
const log = new LogWriter(process.stderr);
const output = new JSONEventWriter(process.stdout);

input.on('data', (inObj) => {
  try {
    const value = inObj[propName].toUpperCase();
    const outObj = {
      id: inObj.id,
      [propName]: value,
    };

    log.info(`Did extract "${propName}" from: "${inObj.id}"`);

    if (!output.write(outObj)) {
        input.pause();
      output.once('drain', () => input.resume());
    }
  } catch(err) {
    log.error(err);
  }
});

input.on('end', () => {
  output.end();
});

input.resume();
