const path = require("path");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
  entry: {
    bundle: "./assets/js/index.js",
  },
  output: {
    filename: "[name].js",
    path: path.resolve(__dirname, "public/js"),
  },
  devtool: "source-map",
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: [/node_modules/],
        use: [{ loader: "babel-loader" }],
      },
    ],
  },
  plugins: [
    new CopyPlugin({
      patterns: [{ from: "assets/css", to: "../css" }],
    }),
  ],
};
