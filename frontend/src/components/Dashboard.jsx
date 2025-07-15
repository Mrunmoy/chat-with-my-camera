import React, { useEffect, useState } from 'react';
import '../styles/Dashboard.css';

import CameraCard from './CameraCard';

/**
 * Dashboard
 *
 * - Fetches all cameras from your API.
 * - Maps each camera to your reusable <CameraCard />.
 */
function Dashboard() {
  // === State: all cameras ===
  const [cameras, setCameras] = useState([]);

  // === Fetch cameras once when Dashboard mounts ===
  useEffect(() => {
    fetch('http://localhost:8080/cameras')
      .then((res) => res.json())
      .then((data) => {
        console.log('Fetched cameras:', data);  // ✅ Add this!
        setCameras(data);
      })
      .catch((err) => console.error('Error fetching cameras:', err));
  }, []);

  return (
    <div className="dashboard">
      {cameras.length > 0 ? (
        cameras.map((camera) => (
          <CameraCard
            key={camera.id}
            name={camera.id}
            label={camera.type}   // Or use camera.label if you have it!
            thumbnail={camera.thumbnail}
            color={getRandomPastel()}
          />
        ))
      ) : (
        <div>Loading cameras...</div>
      )}
    </div>
  );
}

const getRandomPastel = () => {
  const pastelColors = [
    'pastelRed', 'pastelOrange', 'pastelYellow',
    'pastelGreen', 'pastelBlue', 'pastelPurple', 'pastelPink'
  ];
  return pastelColors[Math.floor(Math.random() * pastelColors.length)];
};


export default Dashboard;
