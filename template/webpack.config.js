const BrowserSyncPlugin = require('browser-sync-webpack-plugin')
const CleanPlugin = require('clean-webpack-plugin')
const HtmlPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const OptimizeCSSAssetsPlugin = require("optimize-css-assets-webpack-plugin");
const ScriptExtHtmlPlugin = require('script-ext-html-webpack-plugin')
const TerserPlugin = require('terser-webpack-plugin')
const path = require('path')

module.exports = {
  mode: 'development',
  optimization: {
    minimizer: [
      new TerserPlugin(),
      new OptimizeCSSAssetsPlugin({}),
    ]
  },
  entry: {
    main: './js/main.js',
  },
  output: {
    publicPath: '/assets/',
    path: path.resolve(__dirname, 'public'),
    filename: '[name].[contenthash].js',
  },
  plugins: [
	new CleanPlugin([
	  path.resolve(__dirname, 'public'),
	  path.resolve(__dirname, 'templates/layout'),
    ]),
    new MiniCssExtractPlugin({
      filename: "[name].[contenthash].css",
    }),
    new HtmlPlugin({
      template: './main.html',
      filename: '../templates/layout/main.html',
      chunks: ['main'],
      inject: 'head',
    }),
    new ScriptExtHtmlPlugin({
      defaultAttribute: 'defer',
    }),
    new BrowserSyncPlugin({
      proxy: 'http://localhost:4040',
      port: 4050,
      open: false,
      files: [
        'templates/**/*.html',
        '{{Name}}',
      ],
      snippetOptions: {
        rule: {
          match: /<\/head>/i,
          fn: function (snippet, match) {
            return snippet + match
          },
        },
      },
    }),
  ],
  module: {
    rules: [{
      test: /\.js$/,
      exclude: /turbolinks/,
      use: [{
        loader: 'babel-loader',
        options: {
          presets: ['@babel/preset-env'],
          plugins: ['@babel/plugin-syntax-dynamic-import'],
        },
      }],
    }, {
      test: /\.css$/,
      use: [{
        loader: MiniCssExtractPlugin.loader,
      }, {
        loader: 'css-loader',
      }, {
		loader: 'postcss-loader',
		options: {
		  plugins: [
            require('postcss-import'),
            require('tailwindcss')('./css/tailwind.js'),
			require('postcss-preset-env'),
		  ],
        },
      }],
    }],
  }
}
