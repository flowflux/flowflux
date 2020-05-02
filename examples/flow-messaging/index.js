
const prefixSize = 20;

function ParseMessage(callback) {
  let cache = '';
  let endIndex = 0;
  const parseEndIndex = () => {
    if (cache.length >= prefixSize) {
      const lengthStr = cache.slice(0, prefixSize);
      const length = Number.parseInt(lengthStr);
      return prefixSize + length;
    }
    return 0;
  };
  return data => {
    cache += data.toString();
    endIndex = endIndex || parseEndIndex();
    while (endIndex && (cache.length >= endIndex)) {
      const msgStr = cache.slice(prefixSize, endIndex);
      cache = cache.slice(endIndex);
      endIndex = parseEndIndex();
      const msg = JSON.parse(msgStr);
      callback(msg);
    }
  };
}

function Envelope(msg) {
  const payload = Buffer.from(JSON.stringify(msg));
  const prefix = `${payload.length}`.padStart(prefixSize, '0');
  const length = Buffer.from(prefix);
  return Buffer.concat([length, payload]);
}

function SendMessage(outputStream) {
  let sequentializer = Promise.resolve();
  return msg => {
    sequentializer = sequentializer.then(() => {
      const waitForDrain = !outputStream.write(Envelope(msg));
      if (waitForDrain) return new Promise(resolve => {
        outputStream.once('drain', resolve);
      });
    });
  };
}

module.exports = {
  ParseMessage,
  Envelope,
  SendMessage,
};
