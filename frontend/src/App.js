import React, { useState } from 'react';
import LoginForm from './components/LoginForm';
import DoctorList from './components/DoctorList';
import SlotSelector from './components/SlotSelector';
import AppointmentManager from './components/AppointmentManager';
import { getToken, removeToken, getUserFromToken } from './api';
import SignupForm from './components/SignupForm';

export default function App() {
  const [loggedIn, setLoggedIn] = useState(!!getToken());
  const [selectedDoctor, setSelectedDoctor] = useState(null);
  const [showSignup, setShowSignup] = useState(false);
  const [appointmentsChanged, setAppointmentsChanged] = useState(false);
  const [slotsChanged, setSlotsChanged] = useState(false);

  const handleLogout = () => {
    removeToken();
    setLoggedIn(false);
    setSelectedDoctor(null);
  };

  const handleAppointmentChanged = () => {
    setAppointmentsChanged(prev => !prev);
    setSlotsChanged(prev => !prev);
  };

  const user = getUserFromToken();
  const role = user?.role;

  if (!loggedIn) {
    return showSignup
      ? <SignupForm onSignup={() => setShowSignup(false)} onSwitchToLogin={() => setShowSignup(false)} />
      : <LoginForm onLogin={() => setLoggedIn(true)} onSwitchToSignup={() => setShowSignup(true)} />;
  }

  return (
    <div>
      <button onClick={handleLogout} style={{ float: 'right', margin: 16 }}>Logout</button>
      <h1 style={{ textAlign: 'center' }}>Doctor Appointment Booking</h1>
      {role === 'doctor' ? (
        <>
          <SlotSelector doctor={{ id: user.user_id, name: user.name }} onBook={handleAppointmentChanged} slotsChanged={slotsChanged} />
        </>
      ) : (
        <>
          <DoctorList onSelect={doc => {
            setSelectedDoctor(null);
            setTimeout(() => setSelectedDoctor(doc), 0);
          }} />
          <SlotSelector doctor={selectedDoctor} onBook={handleAppointmentChanged} slotsChanged={slotsChanged} />
          <AppointmentManager appointmentsChanged={appointmentsChanged} onCancel={handleAppointmentChanged} />
        </>
      )}
    </div>
  );
}
