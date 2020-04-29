const fs = require('fs');
const {SendMessage, ParseMessage} = require('../flow-messaging');
const send = SendMessage(process.stdout);

// Config-parser

process.stdin.on('data', data => {
  const configStr = data.toString().trim();
  const config = JSON.parse(configStr);
  const {input, output} = config;
  send({cmd: 'COMPOSE', input, output});
});