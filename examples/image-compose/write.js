const {ParseMessage} = require('../flow-messaging');
const fs = require('fs');

const actions = {
  WRITE: async ({filename, payload, encoding}) => {
    const data = Buffer.from(payload, encoding);
    fs.writeFile(filename, data, err => {
      if (err) throw err;
    });
  }
}

process.stdin.on('data', ParseMessage(msg => actions[msg.cmd](msg)));