import React from 'react';
import { useParams, Link } from 'react-router-dom';
import Split from 'react-split';
import '../styles/CameraPage.css';

import CameraPreview from './CameraPreview';
import HistoryBox from './HistoryBox';
import ChatBox from './ChatBox';

/**
 * CameraPage
 *
 * Full page layout for viewing a single camera.
 * - Uses a header banner with Back button + camera name
 * - Uses react-split to create a resizable vertical split:
 *    - Main content (preview + history) on top
 *    - Chat box at the bottom
 */
const CameraPage = () => {
  // === Get the cameraId from the URL ===
  const { cameraId } = useParams();

  return (
    <div className="camera-page">
      {/* === Header Banner === */}
      <div className="camera-header">
        {/* Back to Dashboard link */}
        <Link to="/" className="back-button">⬅ Back to Dashboard</Link>

        {/* Camera name pulled from URL param */}
        <div className="camera-header">
        <div className="camera-name">Camera: {cameraId}</div>
        <div className="filler"></div> {/* Takes up remaining space to help centering */}
      </div>
      </div>

      {/* === Main Layout === */}
      <Split
        direction="vertical" // Split top/bottom
        sizes={[80, 20]}      // % for main vs chat area
        minSize={100}         // Minimum size for chat
        gutterSize={10}       // Size of draggable gutter
        className="camera-layout"
      >
        {/* === Top section: preview + history === */}
        <div className="main-content">
          <div className="preview">
            <CameraPreview cameraId={cameraId} />
          </div>
          <div className="history">
            <HistoryBox cameraId={cameraId} />
          </div>
        </div>

        {/* === Bottom section: chat box === */}
        <div className="chat">
          <ChatBox cameraId={cameraId} />
        </div>
      </Split>
    </div>
  );
};

export default CameraPage;
