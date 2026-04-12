import React, { useState, useEffect } from 'react';
import { Upload, FileText, Send, Loader2, User, LogOut, CheckCircle } from 'lucide-react';
import api from './api';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(!!localStorage.getItem('token'));
  const [isLoginMode, setIsLoginMode] = useState(true);
  const [formData, setFormData] = useState({ email: '', password: '', username: '' });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  // Resume & JD State
  const [resume, setResume] = useState(null);
  const [jdText, setJdText] = useState('');
  const [greeting, setGreeting] = useState('');
  const [uploading, setUploading] = useState(false);

  // Handle Login/Register
  const handleAuth = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      if (isLoginMode) {
        const { data } = await api.post('/../../auth/login', { email: formData.email, password: formData.password });
        localStorage.setItem('token', data.token);
        setIsLoggedIn(true);
      } else {
        await api.post('/../../auth/register', formData);
        setIsLoginMode(true);
        alert('Registration successful! Please login.');
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Auth failed');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
  };

  // Handle Resume Upload
  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    setUploading(true);
    const uploadData = new FormData();
    uploadData.append('resume', file);

    try {
      const { data } = await api.post('/resumes/upload', uploadData);
      setResume(data.resume);
      alert('Resume uploaded and parsed!');
    } catch (err) {
      alert('Upload failed: ' + (err.response?.data?.error || err.message));
    } finally {
      setUploading(false);
    }
  };

  // Handle Greeting Generation
  const handleGenerateGreeting = async () => {
    if (!jdText) return alert('Please enter JD text');
    setLoading(true);
    try {
      const { data } = await api.post('/greetings/generate', { jd_text: jdText });
      setGreeting(data.greeting);
    } catch (err) {
      alert('Generation failed: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  if (!isLoggedIn) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
        <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8 space-y-6">
          <div className="text-center">
            <h1 className="text-3xl font-bold text-blue-600">JobLovesMe</h1>
            <p className="text-gray-500 mt-2">{isLoginMode ? 'Welcome back' : 'Create an account'}</p>
          </div>
          
          <form onSubmit={handleAuth} className="space-y-4">
            {!isLoginMode && (
              <input
                type="text"
                placeholder="Username"
                className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none"
                value={formData.username}
                onChange={e => setFormData({...formData, username: e.target.value})}
                required
              />
            )}
            <input
              type="email"
              placeholder="Email"
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none"
              value={formData.email}
              onChange={e => setFormData({...formData, email: e.target.value})}
              required
            />
            <input
              type="password"
              placeholder="Password"
              className="w-full px-4 py-3 rounded-xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none"
              value={formData.password}
              onChange={e => setFormData({...formData, password: e.target.value})}
              required
            />
            {error && <p className="text-red-500 text-sm text-center">{error}</p>}
            <button
              disabled={loading}
              className="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 rounded-xl transition-all disabled:opacity-50 flex items-center justify-center gap-2"
            >
              {loading && <Loader2 className="animate-spin h-5 w-5" />}
              {isLoginMode ? 'Login' : 'Sign Up'}
            </button>
          </form>

          <p className="text-center text-gray-500">
            {isLoginMode ? "Don't have an account? " : "Already have an account? "}
            <button onClick={() => setIsLoginMode(!isLoginMode)} className="text-blue-600 font-semibold">
              {isLoginMode ? 'Sign Up' : 'Login'}
            </button>
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center sticky top-0 z-10">
        <h1 className="text-2xl font-bold text-blue-600">JobLovesMe</h1>
        <div className="flex items-center gap-4">
          <span className="text-gray-600 hidden sm:block">Welcome back</span>
          <button onClick={handleLogout} className="p-2 hover:bg-gray-100 rounded-full text-gray-500 transition-colors">
            <LogOut size={20} />
          </button>
        </div>
      </header>

      <main className="max-w-4xl mx-auto p-6 space-y-8">
        {/* Step 1: Resume Upload */}
        <section className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-2 bg-blue-50 text-blue-600 rounded-lg">
              <Upload size={24} />
            </div>
            <h2 className="text-xl font-bold">1. 简历管理</h2>
          </div>
          
          <div className="border-2 border-dashed border-gray-200 rounded-2xl p-8 text-center hover:border-blue-400 transition-colors cursor-pointer relative">
            <input type="file" accept=".pdf" onChange={handleFileUpload} className="absolute inset-0 opacity-0 cursor-pointer" />
            <div className="flex flex-col items-center gap-2">
              {uploading ? <Loader2 className="animate-spin h-10 w-10 text-blue-500" /> : <FileText className="h-10 w-10 text-gray-400" />}
              <p className="font-medium text-gray-600">点击或拖拽上传 PDF 简历</p>
              <p className="text-sm text-gray-400">仅支持 PDF 格式</p>
            </div>
          </div>
          
          {resume && (
            <div className="mt-4 flex items-center gap-2 text-green-600 bg-green-50 px-4 py-2 rounded-lg">
              <CheckCircle size={16} />
              <span className="text-sm font-medium">当前版本: {resume.file_path.split('/').pop()}</span>
            </div>
          )}
        </section>

        {/* Step 2: JD Input & Greeting */}
        <section className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-2 bg-purple-50 text-purple-600 rounded-lg">
              <Send size={24} />
            </div>
            <h2 className="text-xl font-bold">2. 打招呼生成</h2>
          </div>

          <textarea
            placeholder="粘贴 JD (职位描述) 文字到这里..."
            className="w-full h-40 px-4 py-3 rounded-2xl border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none resize-none"
            value={jdText}
            onChange={e => setJdText(e.target.value)}
          />

          <button
            onClick={handleGenerateGreeting}
            disabled={loading || !resume}
            className="mt-4 w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-4 rounded-2xl transition-all disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {loading && <Loader2 className="animate-spin h-5 w-5" />}
            {resume ? '生成打招呼文本' : '请先上传简历'}
          </button>

          {greeting && (
            <div className="mt-8 p-6 bg-blue-50 rounded-2xl relative">
              <h3 className="text-sm font-bold text-blue-600 mb-3 uppercase tracking-wider">推荐打招呼语</h3>
              <p className="text-gray-800 leading-relaxed whitespace-pre-wrap">{greeting}</p>
              <button 
                onClick={() => navigator.clipboard.writeText(greeting)}
                className="mt-4 text-sm text-blue-600 font-semibold hover:underline"
              >
                复制文本
              </button>
            </div>
          )}
        </section>
      </main>
    </div>
  );
}

export default App;
