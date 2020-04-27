#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const {ParseMessage, SendMessage} = require('../flow-messaging');
const send = SendMessage(process.stdout);

loadMediaTypes((err, mediaTypeForExt) => {
  if (err) {
    console.error('Cannot load media types file:', err);
  } else {
    listen(mediaTypeForExt);
  }
});

function listen(mediaTypeForExt) {
  process.stdin.on('data', ParseMessage(({
    id,
    cmd, 
    ext,
  }) => {
    if (cmd === 'FIND_MEDIA_TYPE') {
      const extension = ext === '' ? '.html' : ext;
      const mediaType = mediaTypeForExt[extension] || null;
      send({
        id,
        cmd: 'PROCESS_MEDIA_TYPE',
        mediaType,
      });
    }
  }));
}

function loadMediaTypes(callback) {
  const csvFilePath = path.join(__dirname, 'media-types.csv');
  fs.readFile(csvFilePath, 'utf8', (err, data) => {
    if (err) {
      callback(err);
      return;
    }
    const acc = {};
    const lines = data.split('\n');
    for (let l of lines) {
      const cell = l.split(',');
      acc[cell[0]] = cell[1];
    }
    callback(undefined, acc);
  });
}