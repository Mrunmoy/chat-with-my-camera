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
 * - Uses a header banner with Back button + camera name.
 * - Uses react-split:
 *    - Outer vertical split: main-content on top, chat box below.
 *    - Inner horizontal split: Preview left, History right.
 */
const CameraPage = () => {
  // === Get the cameraId from the URL ===
  const { cameraId } = useParams();

  return (
    <div className="camera-page">
      {/* === Header Banner === */}
      <div className="camera-header">
        <Link to="/" className="back-button">⬅ Back to Dashboard</Link>
        <div className="camera-name">Camera: {cameraId}</div>
        <div className="filler"></div>
      </div>

      {/* === Outer Split: Vertical === */}
      <Split
        direction="vertical"
        sizes={[80, 20]}   // %: Main content 80%, Chat 20%
        minSize={100}
        gutterSize={10}
        className="camera-layout"
      >
        {/* === Inner Split: Horizontal === */}
        <Split
          direction="horizontal"
          sizes={[60, 40]} // %: Preview 60%, History 40% (adjust as you like)
          minSize={200}
          gutterSize={10}
          className="main-content"
        >
          <div className="preview">
            <CameraPreview cameraId={cameraId} />
          </div>
          <div className="history">
            <HistoryBox cameraId={cameraId} />
          </div>
        </Split>

        {/* === Chat box below === */}
        <div className="chat">
          <ChatBox cameraId={cameraId} />
        </div>
      </Split>
    </div>
  );
};

export default CameraPage;
