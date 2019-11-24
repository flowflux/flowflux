const { Transform } = require('stream');
const { EvntEndBytes } = require('./config');

class EventWriter extends Transform {
  constructor(outputStream) {
    super();
    if (outputStream) {
      this.pipe(outputStream);
    }
  }

  _transform(chunk, enc, callback) {
    const event = Buffer.concat([chunk, EvntEndBytes]);
    callback(null, event);
  }
}

module.exports = EventWriter;