const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin')

const IS_PROD = ['production', 'prod'].includes(process.env.NODE_ENV)
const IS_DEV = ["development", 'dev'].includes(process.env.NODE_ENV);
const webpack = require("webpack");
// const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin')
const smp = new SpeedMeasurePlugin()



module.exports = {
  lintOnSave: true,
  productionSourceMap: !IS_PROD, // 生产环境的source map
  // configureWebpack: smp.wrap({
  //   plugins: [
  //     new BundleAnalyzerPlugin()
  //   ]
  // }),
  chainWebpack: config => {
    config
    // 只打包moment.js需要的语言包
    .plugin("ignore")
    .use(
      new webpack.ContextReplacementPlugin(/moment[/\\]locale$/, /zh-cn$/)
    );

    config
      .plugin('html')
      .tap(args => {
        args[0].title = '理赔2.0-管理系统'
        return args
      })

     config.plugin('monaco').use(new MonacoWebpackPlugin())

  }
}