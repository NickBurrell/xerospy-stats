var path = require('path');

module.exports = {
    context: __dirname,

    mode: "production",

    entry: {
        index: './src/ts/index.tsx'
    },

    output: {
        path: path.resolve('./assets/js/'),
        filename: '[name].js'
    },
    devtool: "source-map",

    module: {
        rules: [
            {
                test: /\.tsx?$/,

                loader: 'ts-loader',
                exclude: /node_modules/
            }
        ]
    },
    resolve: {
        extensions: ['.ts', '.tsx', '.js']
    },
    externals: {
        "react": "React",
        "react-dom": "ReactDOM"
    }
}
