'use strict';

const PSD = require('psd');
const psd = PSD.fromFile("./testdata/test.psd")
psd.parse();
const document = psd.tree().export();
print(document);

function print(document) {
  document.children.forEach((data) => {
    if (data.children) {
      print(data);
    }
    if (data.text) {
      console.log(data.text);
    }
  });
}
