var path = require('path');
var webpack = require('webpack');

module.exports = {
    devServer: true,
    debug: true,
    devtool: '#sourcemap',
    entry: [
        'webpack-dev-server/client?http://localhost:8090',
        'webpack/hot/only-dev-server',
        './ui/app.js'
    ],
    output: {
        path: path.join(__dirname, 'static', 'js'),
        filename: 'app.js',
        publicPath: 'http://localhost:8090/js/'
    },
    plugins: [
        new webpack.HotModuleReplacementPlugin(),
        new webpack.NoErrorsPlugin()
    ],
    resolve: {
        extensions: ['', '.js', '.jsx']
    },
    module: {
        loaders: [
            {
                test: /\.(js|jsx)$/,
                loaders: ['react-hot', 'babel?stage=0&optional[]=runtime'],
                include: path.join(__dirname, 'ui')
            }
        ]
    }
};
