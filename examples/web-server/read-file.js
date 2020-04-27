#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const {ParseMessage, SendMessage} = require('../flow-messaging');
const send = SendMessage(process.stdout);

// STREAMING FILE READER

process.stdin.on('data', ParseMessage(({
  id,
  cmd, 
  url,
}) => {
  if (cmd === 'READ_FILE') {
    const filePath = (url === '') || (url === '/')
      ? 'index.html'
      : url;

    const fullPath = path.join(process.cwd(), 'static', filePath);
    
    try {
      fs.accessSync(fullPath, fs.constants.R_OK);

      const stream = fs.createReadStream(fullPath);

      stream.on('data', chunk => {      
        send({
          id,
          cmd: 'PROCESS_FILE_CHUNK',
          encoding: 'base64',
          payload: chunk.toString('base64'),
        });
      });
      
      stream.on('end', () => {      
        send({
          id,
          cmd: 'CONCLUDE_FILE',
        });
      });

      stream.on('error', () => {      
        send({
          id,
          cmd: 'FAIL_WITH_NOT_FOUND',
          error: err.message,
        });
      });
      
    } catch(err) {
      send({
        id,
        cmd: 'FAIL_WITH_NOT_FOUND',
        error: err.message,
      });
    }
  }
}));
