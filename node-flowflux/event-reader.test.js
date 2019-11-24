const EventReader = require('./event-reader');
const fs = require('fs');
const {
  readAllExpectedEvents,
  inputReaderFromEvents,
  createMapById,
  testDataFilePath
} = require('./test-utils');

test('if testdata works', done => {
  const stats = fs.statSync(testDataFilePath);
  const input = fs.createReadStream(testDataFilePath);
  let byteCount = 0;
  let buffer = Buffer.alloc(0);
  
  input.on('data', chunk => {
    byteCount += chunk.length;
    buffer = Buffer.concat([buffer, chunk]);
  });

  input.on('end', () => {
    expect(stats.size).toBe(byteCount);
    expect(stats.size).toBe(buffer.length);
    done();
  });
});

test('controlled event splitting', done => {
  readAllExpectedEvents((err, expectedEvents) => {
    expect(err).toBe(null);
    const eventForId = createMapById(expectedEvents);
    const input = inputReaderFromEvents(expectedEvents);

    const reader = new EventReader(input);
    let count = 0;

    reader.on('data', data => {
      const eventStr = data.toString('utf8');
      const event = JSON.parse(eventStr);
      const expectedEvent = eventForId[event.id];
      expect(event).toEqual(expectedEvent);
      count++;
    });

    reader.on('end', () => {
      expect(count).toBe(expectedEvents.length);
      done();
    });
  });
});

test('streaming directly from filereader', done => {
  readAllExpectedEvents((err, expectedEvents) => {
    expect(err).toBe(null);
    const eventForId = createMapById(expectedEvents);
    const input = fs.createReadStream(testDataFilePath);

    const reader = new EventReader(input);
    let count = 0;
    
    reader.on('data', data => {
      const eventStr = data.toString('utf8');
      const event = JSON.parse(eventStr);
      const expectedEvent = eventForId[event.id];
      expect(event).toEqual(expectedEvent);
      count++;
    });

    reader.on('end', () => {
      expect(count).toBe(expectedEvents.length);
      done();
    });
  });
});
