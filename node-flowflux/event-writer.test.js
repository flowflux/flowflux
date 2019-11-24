const EventWriter = require('./event-writer');
const { Writable } = require('stream');
const { EvntEnd } = require('./config');
const {
  readAllExpectedEvents,
  createMapById
} = require('./test-utils');

test('event writing', done => {
  readAllExpectedEvents((err, allExpectedEvents) => {
    const originals = allExpectedEvents; // .slice(0, 10);
    const originalById = createMapById(originals);
    let count = 0;
    
    const output = new Writable({
      write(chunk, enc, cb) {
        let str = chunk.toString('utf8');
        str = str.slice(0, str.lastIndexOf(EvntEnd));
        const received = JSON.parse(str);
        const original = originalById[received.id];
        expect(received).toEqual(original);
        count++;
        cb();
        if (count === originals.length) done();
      }
    });

    const writer = new EventWriter(output);
    for (const obj of originals) {
      const str = JSON.stringify(obj);
      const chunk = Buffer.from(str);
      writer.write(chunk);
    }
  });
});
