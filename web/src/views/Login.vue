<template>
  <div class="login-container">
    <div class="login-box">
      <h2>智慧校园系统</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">用户名</label>
          <input 
            type="text" 
            id="username" 
            v-model="loginForm.username" 
            required
          />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input 
            type="password" 
            id="password" 
            v-model="loginForm.password" 
            required
          />
        </div>
        <div class="form-group">
          <button type="submit" class="login-btn">登录</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LoginView',
  data() {
    return {
      loginForm: {
        username: '',
        password: ''
      }
    }
  },
  methods: {
    async handleLogin() {
      try {
        console.log('Attempting login with:', this.loginForm.username);
        // 修改这里：接收完整的 response 对象，或者直接使用 response.data
        const response = await this.$axios.post('/auth/login', {
          username: this.loginForm.username,
          password: this.loginForm.password
        });
        
        // response.data 应该对应后端返回的整个对象 { code: ..., data: { token: ..., user: ... }, msg: ... }
        const responseData = response; 
        console.log('Full login response from backend (inside response.data):', responseData);

        // 现在使用 responseData 来访问 code, msg, 和内层的 data
        if (responseData.code == 200) { 
          console.log('Login successful, code is 200.');
          // 保存token - 注意路径变为 responseData.data.token
          localStorage.setItem('token', responseData.data.token);
          console.log('Token saved to localStorage:', localStorage.getItem('token'));
          // 保存用户信息 - 注意路径变为 responseData.data.user
          localStorage.setItem('userInfo', JSON.stringify(responseData.data.user));
          console.log('UserInfo saved to localStorage:', localStorage.getItem('userInfo'));
          
          console.log('Vue Router instance:', this.$router);
          console.log('Attempting to redirect to /');
          try {
            await this.$router.push('/');
            console.log('Redirect to / initiated.');
          } catch (routerError) {
            console.error('Error during router.push:', routerError);
            alert('登录成功，但页面跳转失败，请查看控制台获取更多信息。');
          }
        } else {
          // 使用 responseData.code 和 responseData.msg
          console.warn('Login failed with code:', responseData.code, 'Message:', responseData.msg);
          alert(responseData.msg || '登录失败');
        }
      } catch (error) {
        console.error('Login error:', error);
        // 检查 error.response.data 是否有更详细的错误信息
        if (error.response && error.response.data && error.response.data.msg) {
          alert(error.response.data.msg);
        } else {
          alert('登录失败，请稍后重试');
        }
      }
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.login-box {
  width: 400px;
  padding: 30px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.login-btn {
  width: 100%;
  padding: 12px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.login-btn:hover {
  background-color: #45a049;
}
</style>