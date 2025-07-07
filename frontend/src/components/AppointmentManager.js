import React, { useEffect, useState } from 'react';
import { fetchAppointments, cancelAppointment, updateAppointment, fetchDoctors, fetchSlots } from '../api';

export default function AppointmentManager({ appointmentsChanged, onCancel }) {
  const [appointments, setAppointments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [cancelingId, setCancelingId] = useState(null);
  const [editingId, setEditingId] = useState(null);
  const [doctors, setDoctors] = useState([]);
  const [slots, setSlots] = useState([]);
  const [selectedDoctor, setSelectedDoctor] = useState(null);
  const [selectedSlot, setSelectedSlot] = useState(null);
  const [rescheduleError, setRescheduleError] = useState('');
  const [rescheduling, setRescheduling] = useState(false);

  const load = () => {
    setLoading(true);
    fetchAppointments()
      .then(setAppointments)
      .catch(() => setError('Failed to load appointments'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

  useEffect(() => {
    load();
  }, [appointmentsChanged]); // refetch when appointmentsChanged changes

  const handleCancel = async (id) => {
    setCancelingId(id);
    setError('');
    try {
      await cancelAppointment(id);
      load();
      if (onCancel) onCancel();
    } catch (err) {
      setError(err.response?.data?.error || 'Cancel failed');
    } finally {
      setCancelingId(null);
    }
  };

  const startEdit = async (app) => {
    setEditingId(app.id);
    setSelectedDoctor({ id: app.doctor_id });
    setSelectedSlot(app.slot_id);
    setRescheduleError('');
    setRescheduling(false);
    // Load doctors and slots for the current doctor
    const docs = await fetchDoctors();
    setDoctors(docs);
    const slots = await fetchSlots(app.doctor_id);
    setSlots(slots.filter(s => !s.is_booked || s.id === app.slot_id));
  };

  const handleDoctorChange = async (e) => {
    const docId = Number(e.target.value);
    setSelectedDoctor(doctors.find(d => d.id === docId));
    setSelectedSlot(null);
    const slots = await fetchSlots(docId);
    setSlots(slots.filter(s => !s.is_booked));
  };

  const handleSlotChange = (e) => {
    setSelectedSlot(Number(e.target.value));
  };

  const handleReschedule = async (app) => {
    setRescheduling(true);
    setRescheduleError('');
    try {
      await updateAppointment(app.id, { slot_id: selectedSlot, doctor_id: selectedDoctor.id });
      setEditingId(null);
      load();
      if (onCancel) onCancel();
    } catch (err) {
      setRescheduleError(err.response?.data?.error || 'Reschedule failed');
    } finally {
      setRescheduling(false);
    }
  };

  if (loading) return <div>Loading appointments...</div>;
  if (error) return <div style={{ color: 'red' }}>{error}</div>;

  return (
    <div style={{ maxWidth: 500, margin: '2rem auto' }}>
      <h2>My Appointments</h2>
      {appointments.length === 0 ? (
        <div>No appointments found.</div>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {appointments.map(app => (
            <li key={app.id} style={{ marginBottom: 8, border: '1px solid #ccc', padding: 8 }}>
              <div>Doctor ID: {app.doctor_id}</div>
              <div>Slot ID: {app.slot_id}</div>
              <div>Status: {app.status}</div>
              {editingId === app.id ? (
                <div style={{ marginTop: 8, border: '1px solid #eee', padding: 8, borderRadius: 4 }}>
                  <div>
                    <label>Doctor: </label>
                    <select value={selectedDoctor?.id || ''} onChange={handleDoctorChange}>
                      <option value='' disabled>Select doctor</option>
                      {doctors.map(doc => (
                        <option key={doc.id} value={doc.id}>{doc.name}</option>
                      ))}
                    </select>
                  </div>
                  <div style={{ marginTop: 8 }}>
                    <label>Slot: </label>
                    <select value={selectedSlot || ''} onChange={handleSlotChange}>
                      <option value='' disabled>Select slot</option>
                      {slots.map(slot => (
                        <option key={slot.id} value={slot.id}>
                          {new Date(slot.start_time).toLocaleString()} - {new Date(slot.end_time).toLocaleString()}
                        </option>
                      ))}
                    </select>
                  </div>
                  <button onClick={() => handleReschedule(app)} disabled={rescheduling || !selectedSlot || !selectedDoctor} style={{ marginTop: 8 }}>
                    {rescheduling ? 'Rescheduling...' : 'Save'}
                  </button>
                  <button onClick={() => setEditingId(null)} style={{ marginLeft: 8, marginTop: 8 }}>Cancel</button>
                  {rescheduleError && <div style={{ color: 'red', marginTop: 8 }}>{rescheduleError}</div>}
                </div>
              ) : (
                <>
                  <button
                    onClick={() => startEdit(app)}
                    style={{ marginTop: 8, marginRight: 8 }}
                  >
                    Edit/Reschedule
                  </button>
                  <button
                    onClick={() => handleCancel(app.id)}
                    disabled={cancelingId === app.id}
                    style={{ marginTop: 8 }}
                  >
                    {cancelingId === app.id ? 'Canceling...' : 'Cancel'}
                  </button>
                </>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
} 