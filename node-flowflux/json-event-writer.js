const { Transform } = require('stream');
const EventWriter = require('./event-writer');

class JSONEventWriter extends Transform {
  constructor(outputStream) {
    super({objectMode: true});
    if (outputStream) {
      this._writer = new EventWriter(outputStream);
      this.pipe(this._writer);
    }
  }

  _transform(obj, enc, callback) {
    try {
      const str = JSON.stringify(obj);
      const chunk = Buffer.from(str);
      callback(null, chunk);
    } catch(err) {
      callback(err);
    }
  }
}

module.exports = JSONEventWriter;