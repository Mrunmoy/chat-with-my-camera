import '../styles/CameraCard.css';
import COLORS from '../constants/colours';
import { Link } from 'react-router-dom';

// === CameraCard Component ===
// Displays a colored block with camera name and a label below.

function CameraCard({ name, label, color }) {
  // Use the passed color name or fallback to a default
  const blockStyle = {
    backgroundColor: COLORS[color] || COLORS.pastelBlue, // fallback
  };

  return (
    <Link to={`/camera/${name}`} className="camera-card">
      <div className="camera-card__block" style={blockStyle}>
        {name}
      </div>
      <div className="camera-card__label">
        {label}
      </div>
    </Link>
  );
}

export default CameraCard;
