import '../styles/CameraCard.css';
import COLORS from '../constants/colours';

// === CameraCard Component ===
// Displays a colored block with camera name and a label below.

function CameraCard({ name, label, color }) {
  // Use the passed color name or fallback to a default
  const blockStyle = {
    backgroundColor: COLORS[color] || COLORS.pastelBlue, // fallback
  };

  return (
    <div className="camera-card">
      {/* === Colorful Block === */}
      <div className="camera-card__block" style={blockStyle}>
        {name}
      </div>

      {/* === Label Text === */}
      <div className="camera-card__label">
        {label}
      </div>
    </div>
  );
}

export default CameraCard;
