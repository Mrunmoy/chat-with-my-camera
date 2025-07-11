import Dashboard from './components/Dashboard';

/**
 * App component serves as the main entry point for the Camera Monitoring Dashboard application.
 * It renders the main title and the Dashboard component.
 *
 * @returns {JSX.Element} The JSX code to render the App component.
 */
function App() {
  return (
    <div>
      <h1>Camera Monitoring Dashboard</h1>
      <Dashboard />
    </div>
  );
}

export default App;
