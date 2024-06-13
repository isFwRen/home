module.exports = {
  lintOnSave: true,
  devServer: {
    // proxy: {
    //   '/const': {
    //     target: 'http://127.0.0.1:30080',
    //     changeOrigin: true,
    //     pathRewrite: {
    //       '^/const': ''
    //     }
    //   }
    // }
  },
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].title = '理赔2.0-录入系统'
        return args
      })
  }
}