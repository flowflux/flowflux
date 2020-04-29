const {ParseMessage, SendMessage} = require('../flow-messaging');
const send = SendMessage(process.stdout);
const Jimp = require('jimp');

const actions = {
  COMPOSE: async ({input, output}) => {
    if (input.length < 2) {
      console.error('Cannot compose less than 2 layers');
    } else {
      input.reverse();
      const images = await Promise.all(input.map(i => Jimp.read(i.filename)))
      const workingLayer = images[0];
      
      let width = workingLayer.bitmap.width;
      let height = workingLayer.bitmap.height;
      for (let i = 1; i < images.length; i++) {
        const layer = images[i];
        width = Math.max(width, layer.bitmap.width);
        height = Math.max(height, layer.bitmap.height);
      }

      workingLayer.contain(width, height);
      for (let i = 1; i < images.length; i++) {
        const layer = images[i];
        layer.contain(width, height);
        const {blendMode} = input[i];
        workingLayer.composite(
          layer, 0, 0, 
          {mode: jimpBlendMode[blendMode]},
        );
      }

      const buff = await workingLayer
        .quality(output.quality)
        .getBufferAsync(jimpMime[output.format]);
      
      send({ cmd: 'OPEN_FILE', filename: output.filename });

      const chunkSize = 2048;
      for (let start = 0; start < buff.length; start+=chunkSize) {
        const end = start + chunkSize;
        const chunk = end < buff.length
          ? buff.slice(start, end)
          : buff.slice(start);
        send({
          cmd: 'WRITE_CHUNK',
          chunk: chunk.toString('base64'),
          encoding: 'base64',
        });
      }

      send({ cmd: 'CLOSE_FILE' });
    }
  }
}

process.stdin.on('data', ParseMessage(msg => actions[msg.cmd](msg)));

const jimpBlendMode = {
  SOURCE_OVER: Jimp.BLEND_SOURCE_OVER,
  DESTINATION_OVER: Jimp.BLEND_DESTINATION_OVER,
  MULTIPLY: Jimp.BLEND_MULTIPLY,
  SCREEN: Jimp.BLEND_SCREEN,
  OVERLAY: Jimp.BLEND_OVERLAY,
  DARKEN: Jimp.BLEND_DARKEN,
  LIGHTEN: Jimp.BLEND_LIGHTEN,
  HARDLIGHT: Jimp.BLEND_HARDLIGHT,
  DIFFERENCE: Jimp.BLEND_DIFFERENCE,
  EXCLUSION: Jimp.BLEND_EXCLUSION,
};

const jimpMime = {
  PNG: Jimp.MIME_PNG,
  JPEG: Jimp.MIME_JPEG,
  BMP: Jimp.MIME_BMP,
};
