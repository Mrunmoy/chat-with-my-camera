import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import Dashboard from './components/Dashboard';
import CameraPage from './components/CameraPage'; // We’ll create this next!

/**
 * App component serves as the main entry point for the Camera Monitoring Dashboard application.
 * It renders the main title and the Dashboard component.
 *
 * @returns {JSX.Element} The JSX code to render the App component.
 */
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/camera/:cameraId" element={<CameraPage />} />
      </Routes>
    </Router>
  );
}

export default App;
