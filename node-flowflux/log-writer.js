const { Transform } = require('stream');
const { EvntEnd } = require('./config');

class LogWriter extends Transform {
  constructor(outputStream) {
    super({objectMode: true});
    if (outputStream) {
      this.pipe(outputStream);
    }
  }

  _transform(logMsg, enc, callback) {
    try {
      const str = `${JSON.stringify(logMsg)}${EvntEnd}`;
      const chunk = Buffer.from(str);
      callback(null, chunk);
    } catch(err) {
      callback(err);
    }
  }

  info(...args) {
    this.write(['info', argsStr(args)]);
  }

  debug(...args) {
    this.write(['debug', argsStr(args)]);
  }

  warn(...args) {
    this.write(['warn', argsStr(args)]);
  }

  error(errorInstance) {
    const ei = errorInstance;
    const obj = {
      columnNumber: ei.columnNumber,
      fileName: ei.fileName, 
      lineNumber: ei.lineNumber, 
      message: ei.message,
      name: ei.name,
      stack: ei.stack 
    }
    try {
      const str = JSON.stringify(obj);
      this.write(['error', str]);
    } catch(err) {
      const comps = [];
      if (err.name) comps.push(err.name);
      if (err.message) comps.push(err.message);
      if (err.stack) comps.push(err.stack);
      this.write(['error-stringify-error', comps.join('\n')]);
    }
  }
}

const argsStr = args => args.map(a => `${a}`).join(' ');

module.exports = LogWriter;