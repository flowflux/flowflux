const { Transform } = require('stream');
const { EvntEndBytes } = require('./config');

class EventReader extends Transform {
  constructor(inputStream) {
    super();
    if (inputStream) {
      this._buffer = null;
      inputStream.pipe(this);
    }
  }

  _transform(chunk, enc, callback) {
    const buff = this._buffer
      ? Buffer.concat([this._buffer, chunk])
      : chunk;
    const reader = eventsReader(buff);

    const progress = () => {
      if (reader.next()) {
        let event = reader.current();
        // const termIdx = event.lastIndexOf('\x00');
        // if (termIdx > -1) {
        //   event = event.slice(0, termIdx);
        // }
        if (!this.push(event)) {
          this.once('drain', progress);
        } else {
          progress();
        }
      } else {
        this._buffer = reader.rest();
        callback();
      }
    };

    progress();
  }
}

function eventsReader(buff) {
  let lastRange = null;
  let currentRange = null;
  let next = () => {
    currentRange = binRangeOf(buff, EvntEndBytes, 0);
    next = () => {
      lastRange = currentRange;
      const fromIndex = currentRange.to + EvntEndBytes.length;
      currentRange = binRangeOf(buff, EvntEndBytes, fromIndex);
      return !!currentRange;
    };
    return !!currentRange;
  };
  const current = () => {
    return buff.slice(currentRange.from, currentRange.to);
  };
  const rest = () => {
    const fromIndex = lastRange.to + EvntEndBytes.length;
    if (fromIndex >= buff.length) return null;
    return buff.slice(fromIndex);
  };
  return { next: () => next(), current, rest };
}

function binRangeOf(searchIn, searchFor, fromIndex) {
  const findingIdx = searchIn.indexOf(searchFor, fromIndex);
  if (findingIdx === -1) return null;
  return {from: fromIndex, to: findingIdx};
}

module.exports = EventReader;