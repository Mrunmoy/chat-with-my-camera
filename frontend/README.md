# Frontend

![React](https://img.shields.io/badge/React-18-61DAFB?logo=react&logoColor=white)

## Overview

This folder contains the **React frontend** for `chat-with-my-camera`. It shows:

- Multi-camera **Dashboard** (auto grid layout).
- **Camera Detail Page** → preview, timeline, and chat.
- Pastel camera cards with random colors.
- Calls backend `/cameras`, `/timeline`, and `/chat`.

---

## How to Run

```bash
cd frontend
npm install
npm run dev
```

Runs with Vite by default on port 5173.

---

## Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── Dashboard.jsx  # Grid of camera cards
│   │   ├── CameraCard.jsx # Individual camera box
│   │   ├── CameraPage.jsx # Layout for single camera view
│   │   ├── CameraPreview.jsx # Shows latest detection image
│   │   ├── HistoryBox.jsx # Timeline picker + Chrono timeline
│   │   ├── ChatBox.jsx    # Chat input & messages
│   ├── styles/            # All CSS modules
│   ├── App.jsx            # Main Router
│   ├── main.jsx           # Entry point
```

---

## How It Works

- Dashboard fetches `/cameras` → maps to `<CameraCard>` → each card links to `/camera/:cameraId`.

- CameraPage uses **react-split** for resizable sections:

- Preview → `<CameraPreview>`
- History → `<HistoryBox>`
- Chat → `<ChatBox>`

- History uses **react-datepicker** for date range and **react-chrono** for event timeline.

- ChatBox sends `{ camera_id, message }` to `/chat` and shows back LLM answer.

---

## Example Flow

- Open Dashboard → see all cameras.
- Click a card → open CameraPage.
- Preview shows latest detection or fallback thumbnail.
- Timeline auto-queries `/timeline?camera_id=...` with date range.
- Chat sends plain text → backend extracts object → LLM answers with real detection info.

---

## Tech Stack

- React + Vite
- react-router-dom (routing)
- react-split-pane (resizable layout)
- react-datepicker, react-chrono
- Custom CSS (styles/)

