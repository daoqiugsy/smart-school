const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    proxy: {
      '/api': { // 所有以 /api 开头的请求都会被代理
        target: 'http://localhost:8080', // 目标后端服务器地址
        changeOrigin: true, // 是否改变源地址，设置为 true，服务器收到的请求头中的 host 为目标target
        // 可选：如果您的后端API路径没有 /api 前缀，但您希望前端请求时带上 /api，可以在这里重写路径
        // pathRewrite: { '^/api': '' } 
      }
    }
  }
})
