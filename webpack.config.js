const path = require('path');

module.exports = {
    mode: "development",
    entry: {
        base: './tssrc/base.ts',
        search: {
            import: './tssrc/search.ts',
            dependOn: 'vendor'
        },
        vendor: ['jquery', 'bootstrap', 'popper.js']
    },
    devtool: 'inline-source-map',
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: [/node_modules/, /tstest/],
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    output: {
        filename: '[name].bundle.js',
        path: path.resolve(__dirname, './static/html/js'),
    },
};