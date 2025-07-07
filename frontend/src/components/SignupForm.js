import React, { useState } from 'react';
import { register } from '../api';

export default function SignupForm({ onSignup, onSwitchToLogin }) {
  const [form, setForm] = useState({ name: '', email: '', password: '', role: 'patient' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = e => setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async e => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      await register(form.name, form.email, form.password, form.role);
      onSignup();
    } catch (err) {
      setError(err.response?.data?.error || 'Signup failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ maxWidth: 320, margin: '2rem auto' }}>
      <h2>Sign Up</h2>
      <input name="name" placeholder="Name" value={form.name} onChange={handleChange} required style={{ width: '100%', marginBottom: 8 }} />
      <input name="email" type="email" placeholder="Email" value={form.email} onChange={handleChange} required style={{ width: '100%', marginBottom: 8 }} />
      <input name="password" type="password" placeholder="Password" value={form.password} onChange={handleChange} required style={{ width: '100%', marginBottom: 8 }} />
      <select name="role" value={form.role} onChange={handleChange} style={{ width: '100%', marginBottom: 8 }}>
        <option value="patient">Patient</option>
        <option value="doctor">Doctor</option>
      </select>
      <button type="submit" disabled={loading} style={{ width: '100%' }}>
        {loading ? 'Signing up...' : 'Sign Up'}
      </button>
      {error && <div style={{ color: 'red', marginTop: 8 }}>{error}</div>}
      <div style={{ marginTop: 8 }}>
        Already have an account?{' '}
        <button type="button" onClick={onSwitchToLogin} style={{ color: 'blue', background: 'none', border: 'none', cursor: 'pointer' }}>
          Login
        </button>
      </div>
    </form>
  );
}