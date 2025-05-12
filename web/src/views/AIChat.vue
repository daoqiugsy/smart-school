<template>
  <div class="chat-container">
    <div class="chat-header">
      <h2>AI 智能助手</h2>
      <div class="debug-toggle">
        <label>
          <input type="checkbox" v-model="debugMode" />
          调试模式
        </label>
        <button @click="testStreamingMode" :disabled="isLoading" class="test-btn">
          测试流式传输
        </button>
      </div>
    </div>

    <div v-if="debugMode" class="debug-panel">
      <div class="debug-info">
        <div><strong>总接收数据块:</strong> {{ debugStats.chunks }}</div>
        <div><strong>总接收字符数:</strong> {{ debugStats.totalChars }}</div>
        <div><strong>平均块延迟:</strong> {{ debugStats.avgChunkDelay.toFixed(2) }}ms</div>
        <div><strong>响应总时间:</strong> {{ debugStats.totalTime }}ms</div>
      </div>
      <div class="debug-chart">
        <div v-for="(chunk, i) in debugStats.chunkSizes" :key="i" 
             class="chunk-bar" 
             :style="{height: Math.min(chunk / 10, 100) + 'px'}">
          <span class="chunk-tooltip">
            块 #{{ i+1 }}: {{ chunk }}字符<br>
            延迟: {{ debugStats.chunkDelays[i] }}ms
          </span>
        </div>
      </div>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div
        v-for="(message, index) in messages"
        :key="index"
        :class="['message', message.type]">
        <div class="message-content" v-if="message.type === 'system' || message.type === 'user'">
          {{ message.content }}
        </div>
        <div v-else-if="message.type === 'ai' && message.isStreaming" class="message-content streaming-content">
          <div v-html="renderedMarkdown(message.content)"></div>
          <span class="cursor"></span>
        </div>
        <div v-else class="message-content markdown-content" v-html="renderedMarkdown(message.content)">
        </div>
      </div>
    </div>

    <div class="chat-input">
      <input
        type="text"
        v-model="inputMessage"
        @keyup.enter="sendMessage"
        placeholder="输入您的问题..."
        :disabled="isLoading"
      />
      <button @click="sendMessage" :disabled="isLoading">
        <span v-if="isLoading">处理中...</span>
        <span v-else>发送</span>
      </button>
    </div>
  </div>
</template>

<script>
import { marked } from 'marked';

export default {
  name: 'AIChatView',
  data() {
    return {
      messages: [
        { type: 'system', content: '您好！我是智慧校园AI助手，有什么可以帮您？' }
      ],
      inputMessage: '',
      isLoading: false,
      fullResponseText: '',
      abortController: null,
      debugMode: true, // 默认开启调试模式
      debugStats: {
        chunks: 0,
        totalChars: 0,
        totalTime: 0,
        avgChunkDelay: 0,
        chunkSizes: [],
        chunkDelays: []
      }
    }
  },
  methods: {
    // 渲染Markdown内容
    renderedMarkdown(content) {
      try {
        if (!content) return '';
        // 设置marked选项
        marked.setOptions({
          breaks: true,      // 将换行符转换为<br>
          gfm: true,         // GitHub风格Markdown
          headerIds: false,  // 不生成header IDs，避免安全问题
          mangle: false,     // 不编码HTML实体，避免显示问题
          sanitize: false    // 不需要sanitize因为我们已经使用Vue的v-html
        });
        return marked.parse(content);
      } catch (error) {
        console.error('Markdown解析错误', error);
        return content; // 解析失败时返回原文本
      }
    },

    // 重置调试统计信息
    resetDebugStats() {
      this.debugStats = {
        chunks: 0,
        totalChars: 0,
        totalTime: 0,
        avgChunkDelay: 0,
        chunkSizes: [],
        chunkDelays: []
      };
    },
    
    // 测试流式传输功能
    testStreamingMode() {
      this.inputMessage = "我有哪些课";
      this.sendMessage();
    },
    
    // 使用简单的fetch方法，直接处理文本响应
    async sendMessage() {
      if (!this.inputMessage.trim() || this.isLoading) return;

      // 重置调试统计信息
      this.resetDebugStats();

      // 添加用户消息
      this.messages.push({ type: 'user', content: this.inputMessage });
      const query = this.inputMessage;
      this.inputMessage = '';
      this.isLoading = true;

      // 添加 AI 占位消息，标记为流式输出中
      const aiIndex = this.messages.length;
      this.messages.push({ 
        type: 'ai', 
        content: '等待响应...',
        isStreaming: true 
      });
      
      const token = localStorage.getItem('token');
      
      // 重新使用fetch API但用更简单的处理方式
      try {
        // 创建AbortController用于取消请求
        if (this.abortController) {
          this.abortController.abort();
        }
        this.abortController = new AbortController();
        
        const response = await fetch('/api/schedule/ai/chat', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
            'Cache-Control': 'no-cache'
          },
          body: JSON.stringify({ query }),
          signal: this.abortController.signal
        });
        
        if (!response.ok) {
          throw new Error(`HTTP 错误! 状态: ${response.status}`);
        }
        
        let responseStartTime = Date.now();
        let receivedChunks = 0;
        let totalCharsReceived = 0;
        let lastChunkTime = responseStartTime;
        let totalDelayTime = 0;
        let currentResult = '';
        let isFirstMessage = true;
        
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        
        // 更简单直接的流处理逻辑
        const processStream = async () => {
          try {
            // eslint-disable-next-line no-constant-condition
            while (true) {
              const { done, value } = await reader.read();
              
              if (done) {
                break;
              }
              
              // 计算统计信息
              const now = Date.now();
              receivedChunks++;
              const chunkDelay = now - lastChunkTime;
              lastChunkTime = now;
              totalDelayTime += chunkDelay;
              
              // 解码数据
              const text = decoder.decode(value);
              totalCharsReceived += text.length;
              
              console.log(`[#${receivedChunks}] 收到数据块: ${text.length}字节, 延迟: ${chunkDelay}ms, 内容: "${text.replace(/\n/g, '\\n')}"`);
              
              // 更新调试统计
              this.debugStats.chunks = receivedChunks;
              this.debugStats.totalChars = totalCharsReceived;
              this.debugStats.avgChunkDelay = totalDelayTime / receivedChunks;
              this.debugStats.totalTime = now - responseStartTime;
              this.debugStats.chunkSizes.push(text.length);
              this.debugStats.chunkDelays.push(chunkDelay);
              
              // 提取和处理SSE消息
              const lines = text.split('\n');
              for (const line of lines) {
                if (line.startsWith('data:')) {
                  try {
                    const jsonStr = line.substring(5).trim();
                    const data = JSON.parse(jsonStr);
                    
                    if (data.code === 200) {
                      if (data.msg === 'success' && data.data) {
                        // 第一条消息，清除等待提示
                        if (isFirstMessage) {
                          isFirstMessage = false;
                          this.messages[aiIndex].content = '';
                        }
                        
                        currentResult += data.data;
                        this.messages[aiIndex].content = currentResult;
                        this.$nextTick(this.scrollToBottom);
                      } else if (data.msg === 'connected') {
                        console.log('已连接服务器');
                      } else if (data.msg === 'completed') {
                        console.log('响应完成');
                      }
                    } else {
                      console.error('服务器错误:', data);
                      if (!currentResult) {
                        this.messages[aiIndex].content = `错误: ${data.msg || '未知错误'}`;
                      }
                    }
                  } catch (err) {
                    // 忽略解析错误，可能是不完整的JSON
                    console.warn('无法解析JSON:', line, err);
                  }
                }
              }
            }
            
            console.log(`流处理完成，总时间: ${Date.now() - responseStartTime}ms, 数据块: ${receivedChunks}`);
            this.debugStats.totalTime = Date.now() - responseStartTime;
          } catch (err) {
            if (err.name !== 'AbortError') {
              console.error('流处理错误:', err);
              if (!currentResult) {
                this.messages[aiIndex].content = `处理错误: ${err.message}`;
              }
            }
          } finally {
            this.messages[aiIndex].isStreaming = false;
            this.isLoading = false;
          }
        };
        
        // 开始处理流
        processStream();
        
      } catch (error) {
        console.error('请求错误:', error);
        this.messages[aiIndex].content = `请求失败: ${error.message}`;
        this.messages[aiIndex].isStreaming = false;
        this.isLoading = false;
      }
    },
    
    scrollToBottom() {
      const c = this.$refs.messagesContainer;
      if (c) c.scrollTop = c.scrollHeight;
    },
    
    // 当组件销毁时取消未完成的请求
    beforeDestroy() {
      if (this.abortController) {
        this.abortController.abort();
      }
    }
  }
}
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 800px;
  margin: 0 auto;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
}

.chat-header {
  padding: 15px;
  background-color: #4CAF50;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-header h2 {
  margin: 0;
  font-size: 18px;
}

.debug-toggle {
  display: flex;
  align-items: center;
}

.debug-toggle label {
  margin-right: 10px;
  font-size: 14px;
  display: flex;
  align-items: center;
}

.debug-toggle input[type="checkbox"] {
  margin-right: 5px;
}

.test-btn {
  padding: 5px 10px;
  background-color: #388E3C;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.test-btn:disabled {
  background-color: #8BC34A;
  cursor: not-allowed;
}

.debug-panel {
  background-color: #f5f5f5;
  border-bottom: 1px solid #e0e0e0;
  padding: 10px 15px;
  font-size: 14px;
}

.debug-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.debug-chart {
  display: flex;
  align-items: flex-end;
  height: 100px;
  overflow-x: auto;
  padding-bottom: 5px;
  border-top: 1px solid #ddd;
  padding-top: 5px;
}

.chunk-bar {
  width: 10px;
  background-color: #4CAF50;
  margin-right: 2px;
  min-height: 1px;
  position: relative;
}

.chunk-bar:hover .chunk-tooltip {
  display: block;
}

.chunk-tooltip {
  display: none;
  position: absolute;
  background-color: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 5px;
  border-radius: 4px;
  bottom: 105%;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
  font-size: 12px;
  z-index: 10;
}

.chat-messages {
  flex: 1;
  padding: 15px;
  overflow-y: auto;
  background-color: #f9f9f9;
}

.message {
  margin-bottom: 15px;
  max-width: 80%;
}

.message.system {
  margin-left: auto;
  margin-right: auto;
  text-align: center;
  background-color: #f0f0f0;
  padding: 8px 12px;
  border-radius: 15px;
  color: #666;
}

.message.user {
  margin-left: auto;
  background-color: #DCF8C6;
  padding: 10px 15px;
  border-radius: 18px 0 18px 18px;
}

.message.ai {
  margin-right: auto;
  background-color: white;
  padding: 10px 15px;
  border-radius: 0 18px 18px 18px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.streaming-content .cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background-color: #333;
  margin-left: 2px;
  animation: blink 1s step-end infinite;
  vertical-align: middle;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.chat-input {
  display: flex;
  padding: 15px;
  background-color: white;
  border-top: 1px solid #e0e0e0;
}

.chat-input input {
  flex: 1;
  padding: 10px 15px;
  border: 1px solid #ddd;
  border-radius: 20px;
  margin-right: 10px;
  font-size: 16px;
}

.chat-input input:disabled {
  background-color: #f9f9f9;
  cursor: not-allowed;
}

.chat-input button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 16px;
  min-width: 80px;
}

.chat-input button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
}

.markdown-content :deep(h1) {
  font-size: 1.6em;
}

.markdown-content :deep(h2) {
  font-size: 1.4em;
}

.markdown-content :deep(h3) {
  font-size: 1.2em;
}

.markdown-content :deep(p) {
  margin: 0.5em 0;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  padding-left: 1.5em;
  margin: 0.5em 0;
}

.markdown-content :deep(li) {
  margin: 0.2em 0;
}

.markdown-content :deep(code) {
  font-family: monospace;
  background-color: #f0f0f0;
  padding: 2px 4px;
  border-radius: 3px;
  font-size: 0.9em;
}

.markdown-content :deep(pre) {
  background-color: #f0f0f0;
  padding: 8px;
  border-radius: 5px;
  overflow-x: auto;
  margin: 0.5em 0;
}

.markdown-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.markdown-content :deep(blockquote) {
  border-left: 4px solid #ddd;
  padding-left: 10px;
  margin: 0.5em 0;
  color: #666;
}

.markdown-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 0.5em 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid #ddd;
  padding: 6px;
  text-align: left;
}

.markdown-content :deep(th) {
  background-color: #f5f5f5;
}

.markdown-content :deep(a) {
  color: #4CAF50;
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(img) {
  max-width: 100%;
  border-radius: 5px;
}

.markdown-content :deep(hr) {
  border: none;
  border-top: 1px solid #eee;
  margin: 1em 0;
}

.streaming-content :deep(h1),
.streaming-content :deep(h2),
.streaming-content :deep(h3),
.streaming-content :deep(h4),
.streaming-content :deep(h5),
.streaming-content :deep(h6),
.streaming-content :deep(p),
.streaming-content :deep(ul),
.streaming-content :deep(ol),
.streaming-content :deep(li),
.streaming-content :deep(code),
.streaming-content :deep(pre),
.streaming-content :deep(blockquote),
.streaming-content :deep(table),
.streaming-content :deep(th),
.streaming-content :deep(td),
.streaming-content :deep(a),
.streaming-content :deep(img),
.streaming-content :deep(hr) {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
}
</style>