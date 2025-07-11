import React, { useState, useEffect } from 'react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css'; // Required CSS for react-datepicker
import { Chrono } from 'react-chrono';
import '../styles/HistoryBox.css';

/**
 * HistoryBox
 *
 * - Displays two date/time pickers: Start and End.
 * - Fetches event history for the selected time range.
 * - Shows events in a scrollable vertical timeline (React Chrono).
 */
const HistoryBox = ({ cameraId }) => {
  // === State: Start and End date/time ===
  const [startDate, setStartDate] = useState(
    new Date(Date.now() - 24 * 60 * 60 * 1000) // Default: 24 hours ago
  );
  const [endDate, setEndDate] = useState(new Date()); // Default: now

  // === State: Event items for the timeline ===
  const [events, setEvents] = useState([]);

  /**
   * useEffect:
   * Whenever startDate or endDate changes,
   * call your backend to get events for that camera in the selected time range.
   *
   * 👉 NOTE: Your backend expects:
   * GET /timeline?camera_id=garage_webcam&start_time=1720000000&end_time=1720999999
   * Timestamps must be EPOCH SECONDS (not ISO).
   */
  useEffect(() => {
    // Convert JS Date objects to epoch seconds
    const startEpoch = Math.floor(startDate.getTime() / 1000);
    const endEpoch = Math.floor(endDate.getTime() / 1000);

    // Build the correct URL with query params
    const url = `/timeline?camera_id=${cameraId}&start_time=${startEpoch}&end_time=${endEpoch}`;

    // Call your backend API
    fetch(url)
      .then((res) => res.json())
      .then((data) => {
        // ✅ Transform each event for React Chrono format:
        // title: time, cardTitle: summary, cardSubtitle/cardDetailedText: extra info
        const chronoItems = data.map((event) => ({
          title: new Date(event.timestamp * 1000).toLocaleString(), // Convert epoch to local time
          cardTitle: event.type || 'Detection Event',
          cardSubtitle: event.details || '',
          cardDetailedText: event.description || '',
        }));
        setEvents(chronoItems);
      })
      .catch((err) => console.error('Error fetching timeline:', err));
  }, [cameraId, startDate, endDate]); // Dependencies: refetch when these change

  return (
    <div className="history-box">
      {/* === Filter bar with Start and End pickers === */}
      <div className="history-filters">
        <div>
          <label>Start:</label>
          <DatePicker
            selected={startDate}
            onChange={(date) => setStartDate(date)}
            showTimeSelect           // Enables time picker
            timeFormat="HH:mm"       // 24h format
            dateFormat="Pp"          // Localized date & time
          />
        </div>

        <div>
          <label>End:</label>
          <DatePicker
            selected={endDate}
            onChange={(date) => setEndDate(date)}
            showTimeSelect
            timeFormat="HH:mm"
            dateFormat="Pp"
          />
        </div>
      </div>

      {/* === Chrono timeline === */}
      <div className="history-timeline">
        {events.length > 0 ? (
          <Chrono
            items={events}
            mode="VERTICAL"      // Vertical timeline layout
            scrollable           // Allow vertical scroll
            hideControls         // Hide left/right controls for vertical mode
            theme={{
              primary: '#6fba1c',
              secondary: '#1e1e1e',
              cardBgColor: '#222'
            }}
          />
        ) : (
          <div>No events for selected range.</div>
        )}
      </div>
    </div>
  );
};

export default HistoryBox;
