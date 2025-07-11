import '../styles/Dashboard.css';

import CameraCard from './CameraCard';

function Dashboard() {
  return (
    <div className="dashboard">
      {/* Temporary dummy cards for now */}
      <CameraCard name="Cam 1" label="Front Door" color="pastelRed" />
      <CameraCard name="Cam 2" label="Backyard" color="pastelOrange" />
      <CameraCard name="Cam 3" label="Garage" color="pastelYellow" />
      <CameraCard name="Cam 4" label="Living Room" color="pastelGreen" />
      <CameraCard name="Cam 5" label="Office" color="pastelBlue" />
      <CameraCard name="Cam 6" label="Side Gate" color="pastelPurple" />
      <CameraCard name="Cam 7" label="Basement" color="pastelPink" />
      <CameraCard name="Cam 8" label="Guest Room" color="pastelRed" />
      <CameraCard name="Cam 9" label="Basement Window" color="pastelOrange" />
      <CameraCard name="Cam 10" label="Garden" color="pastelYellow" />
      {/* Add more cards as needed */}
    </div>
  );
}

export default Dashboard;
