const JSONEventReader = require('./json-event-reader');
const fs = require('fs');
const {
  readAllExpectedEvents,
  createMapById,
  testDataFilePath
} = require('./test-utils');

test('JSON event reading', done => {
  readAllExpectedEvents((err, expectedEvents) => {
    expect(err).toBe(null);
    const eventForId = createMapById(expectedEvents);
    const input = fs.createReadStream(testDataFilePath);

    const reader = new JSONEventReader(input);
    let count = 0;
    
    reader.on('data', event => {
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
