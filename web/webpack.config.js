const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');


function genHtmlWebpackPluginConfig(name) {
    return new HtmlWebpackPlugin({
        chunks: [name],
        inject: true,
        favicon: `./src/ico/${name}.ico`,
        publicPath: "/",
        template: `./src/html/${name}.html`,
        filename: `${name}.html`
    });
}

module.exports = {
    entry: {
        'index': './src/js/index.js',
        '404': './src/js/404.js'
    },
    plugins: [
        new CleanWebpackPlugin(),
        genHtmlWebpackPluginConfig('index'),
        genHtmlWebpackPluginConfig('404'),
        new MiniCssExtractPlugin({
            filename: 'assets/[name].[contenthash].css'
        })
    ],
    output: {
        filename: 'assets/[name].[contenthash].js',
        path: path.resolve(__dirname, 'dist')
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader'
                ]
            }, {
                test: /\.(sass|scss)$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader',
                    'sass-loader'
                ]
            },
            {
                test: /\.js$/,
                exclude: /(node_modules|bower_components)/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env'],
                        cacheDirectory: true
                    }
                }
            },
            {
                test: /\.(png|jpg|gif|svg|mp3)$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: "[name].[contenthash].[ext]",
                            outputPath: "assets",
                            publicPath: "./",
                            useRelativePath: true
                        }
                    }
                ]
            }
        ]
    },
    target: ['web'],
    mode: "development"
};
