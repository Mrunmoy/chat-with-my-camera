import React, { useEffect, useState } from 'react';
import '../styles/CameraPreview.css';

/**
 * CameraPreview
 *
 * - Loads camera info (from /cameras) to get the static thumbnail.
 * - Polls the backend (/latest?camera_id=...) to get the latest detection snapshot.
 * - If there is no detection yet, shows the static thumbnail.
 * - Polls only when the page is visible to save CPU.
 * - Displays "Last updated" timestamp so user knows when preview last refreshed.
 */
const CameraPreview = ({ cameraId }) => {
  // === Local state ===

  // Holds the camera info (from /cameras) including thumbnail image.
  const [camera, setCamera] = useState(null);

  // Holds the full URL to the latest detection snapshot.
  const [snapshotUrl, setSnapshotUrl] = useState('');

  // Holds when the latest detection was fetched.
  const [lastUpdated, setLastUpdated] = useState(null);

  // === Load the camera config on mount ===
  // This runs once when the component mounts, or when cameraId changes.
  useEffect(() => {
    fetch('http://localhost:8080/cameras')
      .then((res) => res.json())
      .then((data) => {
        // Find the camera that matches the ID from the URL.
        const cam = data.find((c) => c.id === cameraId);
        setCamera(cam);
      })
      .catch((err) =>
        console.error('Error fetching camera config:', err)
      );
  }, [cameraId]);

  // === Poll the backend for the latest detection snapshot ===
  // Runs when the component mounts and cameraId changes.
  // Uses visibility API to poll only when the tab is visible.
  useEffect(() => {
    let intervalId; // ID for setInterval so we can clear it.

    // Function to fetch the latest detection for this camera.
    const fetchLatest = () => {
      fetch(`http://localhost:8080/latest?camera_id=${cameraId}`)
        .then((res) => {
          if (!res.ok) {
            throw new Error('No detection yet');
          }
          return res.json();
        })
        .then((data) => {
          if (data.snapshot_file) {
            // Your backend returns "./snapshots/filename.jpg"
            // So split off the filename and build the snapshot URL.
            const fileName = data.snapshot_file.split('/').pop();
            setSnapshotUrl(
              `http://localhost:8080/snapshot?file=${fileName}`
            );
          } else {
            // If no detection, fallback to static thumbnail.
            setSnapshotUrl('');
          }
          // Always update last updated timestamp.
          setLastUpdated(new Date());
        })
        .catch(() => {
          // If there’s an error or no detection, fallback to static.
          setSnapshotUrl('');
          setLastUpdated(new Date());
        });
    };

    // Start polling every 5 seconds.
    const startPolling = () => {
      fetchLatest(); // Run immediately.
      intervalId = setInterval(fetchLatest, 5000);
    };

    // Stop polling.
    const stopPolling = () => {
      clearInterval(intervalId);
    };

    // This function runs whenever the tab visibility changes.
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        startPolling(); // Tab is visible → start polling.
      } else {
        stopPolling();  // Tab hidden → stop polling to save CPU.
      }
    };

    // Add the event listener for tab visibility changes.
    document.addEventListener('visibilitychange', handleVisibilityChange);

    // Start polling immediately if the tab is visible.
    if (document.visibilityState === 'visible') {
      startPolling();
    }

    // Clean up: clear the interval and remove event listener when unmounted.
    return () => {
      stopPolling();
      document.removeEventListener(
        'visibilitychange',
        handleVisibilityChange
      );
    };
  }, [cameraId]); // Runs again if the camera ID changes.

  // === Render ===

  // If the camera config hasn’t loaded yet, show a loading message.
  if (!camera) {
    return <div>Loading camera preview...</div>;
  }

  return (
    <div className="camera-preview">
      {/* Show latest snapshot if available; fallback to static thumbnail */}
      <img
        src={
          snapshotUrl ||
          `http://localhost:8080/static/${camera.thumbnail}`
        }
        alt={`Preview for ${camera.id}`}
      />

      {/* Show the last updated timestamp if we have one */}
      <div className="last-updated">
        {lastUpdated && (
          <p>Last updated: {lastUpdated.toLocaleTimeString()}</p>
        )}
      </div>
    </div>
  );
};

export default CameraPreview;
