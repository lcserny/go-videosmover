const path = require('path');

module.exports = {
    mode: "none",
    entry: {
        base: './tssrc/base.js',
        search: {
            import: './tssrc/search.js',
            dependOn: 'vendor'
        },
        vendor: ['jquery', 'bootstrap', 'popper.js']
    },
    output: {
        filename: '[name].bundle.js',
        path: path.resolve(__dirname, './static/html/js'),
    },
};