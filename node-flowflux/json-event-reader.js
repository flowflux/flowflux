const { Transform } = require('stream');
const EventReader = require('./event-reader');

class JSONEventReader extends Transform {
  constructor(inputStream) {
    super({objectMode: true});
    if (inputStream) {
      this._reader = new EventReader(inputStream);
      this._reader.pipe(this);
    }
  }

  _transform(chunk, enc, callback) {
    try {
      const str = chunk.toString('utf8').replace(/\x00/g, '');
      const obj = JSON.parse(str);
      callback(null, obj);
    } catch(err) {
      callback(err);
    }
  }
}

module.exports = JSONEventReader;