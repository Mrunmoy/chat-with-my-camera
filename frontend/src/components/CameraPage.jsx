import React from 'react';
import Split from 'react-split';
import '../styles/CameraPage.css';

/**
 * CameraPage
 *
 * Page layout for a single camera view.
 * - preview: shows the latest camera photo or live feed
 * - history: shows detection/event history
 * - chat: optional chat box at the bottom
 *
 * Uses plain CSS Flexbox for fixed side-by-side layout.
 */
const CameraPage = ({ preview, history, chat }) => {
  return (
    <Split
      direction="vertical"
      sizes={[80, 20]} // % for main vs chat area
      minSize={100}
      gutterSize={10}
      className="camera-layout"
    >
      <div className="main-content">
        <div className="preview">{preview}</div>
        <div className="history">{history}</div>
      </div>
      <div className="chat">{chat}</div>
    </Split>
  );
};

export default CameraPage;
