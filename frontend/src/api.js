import axios from 'axios';

const API_BASES = {
  auth: 'http://localhost:8080',
  user: 'http://localhost:8081',
  schedule: 'http://localhost:8082',
  appointment: 'http://localhost:8083',
};

export function setToken(token) {
  localStorage.setItem('jwt', token);
}

export function getToken() {
  return localStorage.getItem('jwt');
}

export function removeToken() {
  localStorage.removeItem('jwt');
}

export async function login(email, password) {
  const res = await axios.post(`${API_BASES.auth}/login`, { email, password });
  setToken(res.data.token);
  return res.data;
}

export async function fetchDoctors() {
  const res = await axios.get(`${API_BASES.user}/doctors`);
  return res.data;
}

export async function fetchSlots(doctorId) {
  const token = getToken();
  const user = token ? JSON.parse(atob(token.split('.')[1])) : null;
  if (user && user.role === 'doctor' && Number(user.user_id) === Number(doctorId)) {
    const res = await axios.get(`${API_BASES.schedule}/slots?doctor_id=${doctorId}`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    return res.data;
  }
  const res = await axios.get(`${API_BASES.schedule}/slots?doctor_id=${doctorId}`);
  return res.data;
}

export async function bookAppointment(doctorId, slotId) {
  const token = getToken();
  const res = await axios.post(
    `${API_BASES.appointment}/appointments`,
    { doctor_id: doctorId, slot_id: slotId },
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return res.data;
}

export async function fetchAppointments() {
  const token = getToken();
  const res = await axios.get(`${API_BASES.appointment}/appointments`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data;
}

export async function cancelAppointment(appointmentId) {
  const token = getToken();
  const res = await axios.delete(`${API_BASES.appointment}/appointments/${appointmentId}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data;
}

export async function updateAppointment(appointmentId, data) {
  const token = getToken();
  const res = await axios.put(
    `${API_BASES.appointment}/appointments/${appointmentId}`,
    data,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return res.data;
}

export async function register(name, email, password, role) {
  const res = await axios.post('http://localhost:8080/register', {
    name,
    email,
    password,
    role,
  });
  return res.data;
}

export async function createSlot(start_time, end_time) {
  const token = getToken();
  const res = await axios.post(
    `${API_BASES.schedule}/slots`,
    { start_time, end_time },
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return res.data;
}

export async function deleteSlot(slotId) {
  const token = getToken();
  const res = await axios.delete(
    `${API_BASES.schedule}/slots/${slotId}`,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return res.data;
}

// Decode JWT to get user info (role, id, etc.)
export function getUserFromToken() {
  const token = getToken();
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload;
  } catch {
    return null;
  }
} 