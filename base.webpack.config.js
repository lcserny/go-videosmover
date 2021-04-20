const path = require('path');

module.exports = {
    mode: "none",
    entry: './static/html/js/src/base.js',
    output: {
        filename: 'base.js',
        path: path.resolve(__dirname, './static/html/js/dist'),
    },
};