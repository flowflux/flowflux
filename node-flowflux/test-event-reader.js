#!/usr/bin/env node

const EventReader = require('../../node-flowflux/event-reader');
const ErrorWriter = require('../../node-flowflux/error-writer');
const { EvntEnd } = require('../../node-flowflux/config');

const stdout = process.stdout;

const input = new EventReader(process.stdin);
const error = new ErrorWriter(process.stderr);

input.on('data', (chunk) => {
  const str = chunk.toString('utf8'); // .replace(/\x00/g, '');
  
  try {
    const obj = JSON.parse(str);
    obj.first_name = obj.first_name.toUpperCase();
    obj.last_name = obj.last_name.toUpperCase();
    const lineStr = `${JSON.stringify(obj)}${EvntEnd}\n`;
    const line = Buffer.from(lineStr);

    if (!stdout.write(line)) {
      input.pause();
      stdout.once('drain', () => input.resume());
    }
  } catch(err) {
    error.write(err);
    // stderr.write();
    // console.log(`${err.toString()}`);
    // console.log(`Couldn't parse: >>${str}<<`);
  }
});

input.on('end', () => {
  stdout.end();
});

input.resume();


