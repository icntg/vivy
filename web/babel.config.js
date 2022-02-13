// vue.config.js
const path = require("path");
const TerserPlugin = require("terser-webpack-plugin");
const glob = require("glob");
const {merge} = require("webpack-merge");
// const HtmlWebpackExternalsPlugin = require('html-webpack-externals-plugin');

const page = function() {
  let entryHtml = glob.sync("src/views" + "/*/*.html"); // 在vue-cli3中配置文件的路径不需要加相对路径，也不需要找绝对路径
  let obj = {};
  entryHtml.forEach(filePath => {
    let filename = filePath.substring(filePath.lastIndexOf("/") + 1, filePath.lastIndexOf("."));
    let entryName = filePath.substring(0, filePath.lastIndexOf(".")); // 我为了方便使用的html的路径，最后只是改写后缀
    let conf = {
      entry: entryName + '.ts', // 此处需要注意，如果引用ts和html不同名称与路径的文件，在此处做相应改写
      template: filePath, // 此处是html文件的访问路径
      filename: filename + ".html",
      chunks: ["chunk-vendors", "chunk-common", filename],
    };
    if (process.env.NODE_ENV === "production") {
      conf = merge(conf, {
        minify: {
          removeComments: true,
          collapseWhitespace: true,
          removeAttributeQuotes: true
        },
        chunksSortMode: "dependency"
      });
    }
    obj[filename] = conf
  });
  return obj;
};


module.exports = {
  presets: [
    '@vue/cli-plugin-babel/preset'
  ]
}
