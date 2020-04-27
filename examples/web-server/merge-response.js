#!/usr/bin/env node

const {ParseMessage, SendMessage} = require('../flow-messaging');
const send = SendMessage(process.stdout);

const dispatchForID = {};

process.stdin.on('data', ParseMessage(msg => {
  switch (msg.cmd) {
    case 'FAIL_WITH_NOT_FOUND':
      send(msg);
      delete dispatchForID[msg.id];
      break;
    default:
      let dispatch = dispatchForID[msg.id];
      if (!dispatch) {
        dispatch = makeBufferedDispatch(msg.id);
        dispatchForID[msg.id] = dispatch;
      }
      dispatch(msg);
      break;
  }
}));

function CONCLUDE_FILE(msg) {
  send(msg);
  delete dispatchForID[msg.id];
}

function makeBufferedDispatch(id) {
  const buffer = [];
  const commands = {
    PROCESS_MEDIA_TYPE: msg => {
      const dispatch = makeDispatch(msg.mediaType);
      buffer.forEach(dispatch);
      dispatchForID[id] = dispatch;
    },
    PROCESS_FILE_CHUNK: msg => {
      buffer.push(msg);
    },
    CONCLUDE_FILE,
  };
  return msg => commands[msg.cmd](msg);
}

function makeDispatch(contentType) {
  const commands = {
    PROCESS_FILE_CHUNK: msg => {
      send({...msg, contentType});
    },
    CONCLUDE_FILE,
  };
  return msg => commands[msg.cmd](msg);
}
