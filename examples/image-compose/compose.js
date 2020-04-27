const {ParseMessage, SendMessage} = require('../flow-messaging');
const fs = require('fs');
const Jimp = require('jimp');

const actions = {
  COMPOSE: async ({config}) => {
    if (config.length < 2) {
      console.error('Cannot compose less than 2 layers');
    } else {
      const workingLayer = await Jimp.read(config[0][0]);
      for (let i = 1; i < config.length; i++) {
        const layer = await Jimp.read(config[i][0]);
        workingLayer.composite(layer, 0, 0);
      }
      const buff = await workingLayer.quality(75).getBufferAsync(Jimp.MIME_JPEG);
      process.stdout.write(buff);
    }
  }
}

process.stdin.on('data', ParseMessage(msg => actions[msg.cmd](msg)));