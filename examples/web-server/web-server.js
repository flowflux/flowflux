#!/usr/bin/env node

const http = require('http');
const path = require('path');
const cuid = require('cuid');
const {SendMessage, ParseMessage} = require('../flow-messaging');
const send = SendMessage(process.stdout);

const port = 3000;
const pending = {};

const server = http.createServer((request, response) => {
  const id = cuid();
  
  pending[id] = {id, request, response, needsHeaders: true};

  response.setTimeout(10000);

  send({id, cmd: 'READ_FILE', url: request.url});
  send({id, cmd: 'FIND_MEDIA_TYPE', ext: path.extname(request.url)});
});

server.listen(port, (err) => {
  if (err) {
    return console.error('something bad happened', err);
  }
});

const commands = {
  PROCESS_FILE_CHUNK: ({id, contentType, payload, encoding}) => {
    const {response, needsHeaders} = pending[id];
    if (needsHeaders) {
      response.setHeader('Content-Type', contentType);
      pending[id].needsHeaders = false;
    }
    response.write(Buffer.from(payload, encoding));
  },
  CONCLUDE_FILE: ({id}) => {
    const {response} = pending[id];
    response.end();
    delete pending[id];
  },
  FAIL_WITH_NOT_FOUND: ({id, error}) => {
    const {response} = pending[id];
    response.writeHead(404, { 'Content-Type': 'text/html' });
    response.end(`<h1>NOT FOUND</h1><p>${error}</p>`, 'utf-8')
    delete pending[id];
  },
};

process.stdin.on('data', ParseMessage(msg => commands[msg.cmd](msg)));
