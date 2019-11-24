#!/usr/bin/env node

const EventReader = require('../../node-flowflux/event-reader');
const LogWriter = require('../../node-flowflux/log-writer');
const EventWriter = require('../../node-flowflux/event-writer');

const input = new EventReader(process.stdin);
const log = new LogWriter(process.stderr);
const output = new EventWriter(process.stdout);

input.on('data', (chunk) => {
  const str = chunk.toString('utf8').replace(/\x00/g, '');
  
  try {
    const obj = JSON.parse(str);
    obj.first_name = obj.first_name.toUpperCase();
    obj.last_name = obj.last_name.toUpperCase();
    log.info('Did process:', obj.id);
    const lineStr = JSON.stringify(obj);
    const line = Buffer.from(lineStr);

    if (!output.write(line)) {
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


