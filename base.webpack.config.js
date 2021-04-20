const path = require('path');

module.exports = {
    mode: "none",
    entry: './tssrc/base.js',
    output: {
        filename: 'base.js',
        path: path.resolve(__dirname, './static/html/js'),
    },
};