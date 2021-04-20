const path = require('path');

module.exports = {
    mode: "none",
    entry: './static/html/js/src/search.js',
    output: {
        filename: 'search.js',
        path: path.resolve(__dirname, './static/html/js/dist'),
    },
};