import React, { useEffect, useState } from 'react';
import { fetchSlots, bookAppointment, createSlot, getUserFromToken, deleteSlot } from '../api';

export default function SlotSelector({ doctor, onBook, slotsChanged }) {
  const [slots, setSlots] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [booking, setBooking] = useState(false);
  const [bookedSlotId, setBookedSlotId] = useState(null);
  const [failedSlotIds, setFailedSlotIds] = useState([]);
  const [slotInputs, setSlotInputs] = useState([
    { start: '', end: '' }
  ]);
  const [creating, setCreating] = useState(false);
  const [deletingSlotId, setDeletingSlotId] = useState(null);

  const refreshSlots = () => {
    setLoading(true);
    fetchSlots(doctor.id)
      .then(setSlots)
      .catch(() => setError('Failed to load slots'))
      .finally(() => {
        setLoading(false);
        setBookedSlotId(null);
      });
  };

  useEffect(() => {
    if (!doctor) return;
    refreshSlots();
  }, [doctor, slotsChanged]);

  // Reset failedSlotIds and bookedSlotId when doctor or slotsChanged changes
  useEffect(() => {
    setFailedSlotIds([]);
    setBookedSlotId(null);
  }, [doctor, slotsChanged]);

  const handleSlotInputChange = (idx, field, value) => {
    setSlotInputs(inputs => inputs.map((input, i) => i === idx ? { ...input, [field]: value } : input));
  };

  const handleAddSlotInput = () => {
    setSlotInputs(inputs => [...inputs, { start: '', end: '' }]);
  };

  const handleRemoveSlotInput = (idx) => {
    setSlotInputs(inputs => inputs.length > 1 ? inputs.filter((_, i) => i !== idx) : inputs);
  };

  const handleCreateSlots = async (e) => {
    e.preventDefault();
    setCreating(true);
    setError('');
    try {
      for (const input of slotInputs) {
        if (!input.start || !input.end) continue;
        await createSlot(new Date(input.start).toISOString(), new Date(input.end).toISOString());
      }
      setSlotInputs([{ start: '', end: '' }]);
      refreshSlots();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create slot(s)');
    } finally {
      setCreating(false);
    }
  };

  const handleBook = async (slot) => {
    setBooking(true);
    setError('');
    try {
      await bookAppointment(doctor.id, slot.id);
      setBookedSlotId(slot.id);
      setSlots(slots => slots.filter(s => s.id !== slot.id)); // Remove slot locally
      onBook();
      refreshSlots(); // Optionally, still refresh from backend
    } catch (err) {
      const msg = err.response?.data?.error || 'Booking failed';
      setError(msg);
      if (msg === 'Slot already booked') {
        setFailedSlotIds(ids => [...ids, slot.id]);
      }
    } finally {
      setBooking(false);
    }
  };

  const handleDeleteSlot = async (slotId) => {
    setDeletingSlotId(slotId);
    setError('');
    try {
      await deleteSlot(slotId);
      refreshSlots();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to delete slot');
    } finally {
      setDeletingSlotId(null);
    }
  };

  if (!doctor) return null;
  if (loading) return <div>Loading slots...</div>;
  if (error) return <div style={{ color: 'red' }}>{error}</div>;

  const user = getUserFromToken();
  const userRole = user?.role;

  if (userRole === 'doctor' && Number(user?.user_id) === Number(doctor.id)) {
    // Show all slots, including booked
  }

  const visibleSlots = userRole === 'doctor' && Number(user?.user_id) === Number(doctor.id)
    ? slots
    : slots.filter(slot => !slot.is_booked);

  return (
    <div style={{ maxWidth: 400, margin: '2rem auto' }}>
      <h2>Slots for {doctor.name}</h2>
      {/* Doctor slot creation form: only show if logged-in user is a doctor and is viewing their own profile */}
      {userRole === 'doctor' && Number(user?.user_id) === Number(doctor.id) && (
        <form onSubmit={handleCreateSlots} style={{ marginBottom: 24, border: '1px solid #ccc', padding: 16, borderRadius: 8 }}>
          <h3>Create Slots</h3>
          {slotInputs.map((input, idx) => (
            <div key={idx} style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
              <input
                type="datetime-local"
                value={input.start}
                onChange={e => handleSlotInputChange(idx, 'start', e.target.value)}
                required
                style={{ marginRight: 8 }}
              />
              <span>to</span>
              <input
                type="datetime-local"
                value={input.end}
                onChange={e => handleSlotInputChange(idx, 'end', e.target.value)}
                required
                style={{ marginLeft: 8, marginRight: 8 }}
              />
              <button type="button" onClick={() => handleRemoveSlotInput(idx)} disabled={slotInputs.length === 1} style={{ marginLeft: 4 }}>
                &times;
              </button>
            </div>
          ))}
          <button type="button" onClick={handleAddSlotInput} style={{ marginRight: 8 }}>Add Another Slot</button>
          <button type="submit" disabled={creating}>Create Slot(s)</button>
        </form>
      )}
      {/* Doctor's own slot list: show all slots with status */}
      {userRole === 'doctor' && Number(user?.user_id) === Number(doctor.id) && (
        <div>
          <h3>Your Slots</h3>
          {visibleSlots.length === 0 ? (
            <div>No slots created yet.</div>
          ) : (
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {visibleSlots.map(slot => (
                <li key={slot.id} style={{ marginBottom: 8, padding: 8, border: '1px solid #eee', borderRadius: 4, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <span>{new Date(slot.start_time).toLocaleString()} - {new Date(slot.end_time).toLocaleString()}</span>
                  <span style={{ marginLeft: 12, fontWeight: 'bold', color: slot.is_booked ? 'red' : 'green' }}>
                    {slot.is_booked ? 'Booked' : 'Available'}
                  </span>
                  <button
                    onClick={() => handleDeleteSlot(slot.id)}
                    disabled={deletingSlotId === slot.id || slot.is_booked}
                    style={{ marginLeft: 16 }}
                    title={slot.is_booked ? 'Cannot delete a booked slot' : 'Delete slot'}
                  >
                    {deletingSlotId === slot.id ? 'Deleting...' : 'Delete'}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
      {/* Slot booking UI: only show if logged-in user is a patient */}
      {userRole === 'patient' && (visibleSlots.length === 0 ? (
        <div>No available slots.</div>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {visibleSlots.map(slot => (
            <li key={slot.id} style={{ marginBottom: 8 }}>
              <button
                onClick={() => handleBook(slot)}
                disabled={booking || bookedSlotId === slot.id || failedSlotIds.includes(slot.id)}
                style={{
                  width: '100%',
                  background: failedSlotIds.includes(slot.id) ? '#ccc' : undefined,
                  color: failedSlotIds.includes(slot.id) ? '#888' : undefined,
                  cursor: failedSlotIds.includes(slot.id) ? 'not-allowed' : undefined
                }}
              >
                {new Date(slot.start_time).toLocaleString()} - {new Date(slot.end_time).toLocaleString()}
                {bookedSlotId === slot.id ? ' (Booked!)' : ''}
                {failedSlotIds.includes(slot.id) ? ' (Unavailable)' : ''}
              </button>
            </li>
          ))}
        </ul>
      ))}
      {error && <div style={{ color: 'red', marginTop: 8 }}>{error}</div>}
    </div>
  );
} 