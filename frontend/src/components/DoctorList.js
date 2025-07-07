import React, { useEffect, useState } from 'react';
import { fetchDoctors } from '../api';

export default function DoctorList({ onSelect }) {
  const [doctors, setDoctors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  useEffect(() => {
    fetchDoctors()
      .then(setDoctors)
      .catch(() => setError('Failed to load doctors'))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading doctors...</div>;
  if (error) return <div style={{ color: 'red' }}>{error}</div>;

  const filteredDoctors = doctors.filter(doc =>
    doc.name.toLowerCase().includes(search.toLowerCase()) ||
    doc.email.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div style={{ maxWidth: 400, margin: '2rem auto' }}>
      <h2>Select a Doctor</h2>
      <input
        type="text"
        placeholder="Search doctor by name or email"
        value={search}
        onChange={e => setSearch(e.target.value)}
        style={{ width: '100%', marginBottom: 12, padding: 6 }}
      />
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {filteredDoctors.map(doc => (
          <li key={doc.id} style={{ marginBottom: 8 }}>
            <button onClick={() => onSelect(doc)} style={{ width: '100%' }}>
              {doc.name} ({doc.email})
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
} 