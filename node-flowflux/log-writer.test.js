const LogWriter = require('./log-writer');
const { Writable } = require('stream');
const { EvntEnd } = require('./config');

test('errors streaming', () => {
  const errorEvents = [];
  const originals = [];

  const output = new Writable({
    write(chunk, enc, cb) {
      errorEvents.push(chunk.toString('utf8'));
      cb();
    }
  });
  const log = new LogWriter(output);
  
  for (const err of generateErrors(5000)) {
    originals.push(err);
    log.error(err);
  }

  const received = errorEvents
    .map(e => e.slice(0, e.lastIndexOf(EvntEnd)))
    .map(s => JSON.parse(s))
    .map(([t, d]) => [t, JSON.parse(d)]);

  expect(received.length).toBe(originals.length);
  
  for (let i = 0; i < originals.length; i++) {
    const { message: originalMsg } = originals[i];
    const [type, errObj] = received[i];
    expect(type).toBe('error');
    expect(originalMsg).toEqual(errObj.message);
  }
});

test('log streaming', () => {
  const receivedEvents = [];
  const originals = [];

  const output = new Writable({
    write(chunk, enc, cb) {
      receivedEvents.push(chunk.toString('utf8'));
      cb();
    }
  });
  const log = new LogWriter(output);
  
  originals.push(...generateLogMessages(1666, 'info'));
  originals.push(...generateLogMessages(1666, 'debug'));
  originals.push(...generateLogMessages(1666, 'warn'));

  for (const msg of originals) {
    const [type, ...args] = msg;
    log[type](...args);
  }

  const received = receivedEvents
    .map(e => e.slice(0, e.lastIndexOf(EvntEnd)))
    .map(s => JSON.parse(s));
  
  expect(received.length).toBe(originals.length);
  
  for (let i = 0; i < originals.length; i++) {
    const [origType, ...origArgs] = originals[i];
    const origArgsStr = origArgs.join(' ');
    const [receivedType, receivedArgsStr] = received[i];
    expect(origType).toBe(receivedType);
    expect(origArgsStr).toBe(receivedArgsStr);
  }
});

function* generateErrors(amount) {
  let index = 0;
  while (index < amount) {
    try {
      throw new Error(`Test error: ${index + 1}`);
    } catch(e) {
      yield e;
    }
    index++;
  }
}

function* generateLogMessages(amount, type) {
  let index = 0;
  while (index < amount) {
    yield [type, 'log', 'message', index + 1];
    index++;
  }
}
