const {ParseMessage} = require('../flow-messaging');
const fs = require('fs');

function Sequentialize(outputStream) {
  let p = Promise.resolve();

  return [
    buff => {
      p = p.then(() => {
        const waitForDrain = !outputStream.write(buff);
        if (waitForDrain) return new Promise(resolve => {
          outputStream.once('drain', resolve);
        });
      });
    },
    () => {
      p = p.then(() => {
        outputStream.end();
      });
    },
  ]
}

let write, end;

const actions = {
  OPEN_FILE: ({filename}) => {
    [write, end] = Sequentialize(fs.createWriteStream(filename));
  },
  WRITE_CHUNK: ({chunk, encoding}) => write(Buffer.from(chunk, encoding)),
  CLOSE_FILE: () => end(),
};

process.stdin.on('data', ParseMessage(msg => actions[msg.cmd](msg)));