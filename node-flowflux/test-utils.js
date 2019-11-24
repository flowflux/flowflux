const fs = require('fs');
const { EvntEnd } = require('./config');
const { Readable } = require('stream');
const testDataFilePath = 'testdata.txt';

const readAllExpectedEvents = cb => {
  fs.readFile(testDataFilePath, 'utf8', (err, data) => {
    if (err) cb(err);
    else {
      const comps = data.split(EvntEnd).filter(c => !!c.trim().length);
      const events = comps.map(cs => JSON.parse(cs));
      cb(null, events);
    }
  });
};

const inputReaderFromEvents = events => {
  const comps = events.map(e => JSON.stringify(e));
  const dataStr = `${comps.join(EvntEnd)}${EvntEnd}`;
  const dataBuff = Buffer.from(dataStr);
  const reader = new Readable();
  reader._read = () => {};
  reader.push(dataBuff);
  reader.push(null);
  return reader;
};

const createMapById = items => items.reduce((acc, itm) => {
  acc[itm.id] = itm;
  return acc;
}, {});

module.exports = {
  readAllExpectedEvents,
  inputReaderFromEvents,
  createMapById,
  testDataFilePath
};