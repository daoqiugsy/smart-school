module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 您的Go后端服务地址
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/api'
        }
      }
    }
  }
}