import React, { useEffect, useState } from 'react';
import '../styles/CameraPreview.css';

/**
 * CameraPreview
 *
 * Shows the latest camera thumbnail for a given camera ID.
 * - Calls the backend API `/api/cameras` to get camera data.
 * - Picks the matching camera based on the `cameraId` prop.
 * - Displays the thumbnail image, resized to max 640x480.
 */
const CameraPreview = ({ cameraId }) => {
  // === Local state to store the camera object ===
  const [camera, setCamera] = useState(null);

  /**
   * === useEffect ===
   *
   * Runs ONCE when the component mounts (or whenever cameraId changes).
   * - Calls the API endpoint `/api/cameras` to get all cameras.
   * - Finds the camera whose ID matches `cameraId`.
   * - Updates local state with that camera object.
   */
  useEffect(() => {
    fetch('/api/cameras')
      .then((res) => res.json())
      .then((data) => {
        // Find the camera that matches the given ID
        const cam = data.find((c) => c.id === cameraId);
        setCamera(cam); // Save it to state
      })
      .catch((err) => console.error(err)); // Log any errors
  }, [cameraId]); // 👈 Dependency: runs again if cameraId changes

  // If the camera is still loading, show placeholder
  if (!camera) {
    return <div>Loading camera preview...</div>;
  }

  // Once loaded, display the thumbnail image
  return (
    <div className="camera-preview">
      <img
        src={`http://localhost:8080${camera.thumbnail}`}
        alt={`Thumbnail for ${camera.id}`}
      />
    </div>
  );
};

export default CameraPreview;
