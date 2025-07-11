import React, { useState } from 'react';
import '../styles/ChatBox.css';

/**
 * ChatBox
 *
 * Simple placeholder chat:
 * - Lets user type messages to an LLM.
 * - Shows conversation history: user messages + fake LLM replies.
 * - In the future, you’ll replace the fake LLM with a real API call.
 */
const ChatBox = ({ cameraId }) => {
  // === Local state ===
  const [messages, setMessages] = useState([]); // [{ sender: 'user'|'llm', text: '...' }]
  const [input, setInput] = useState('');

  // === Handle send ===
  const handleSend = () => {
    if (!input.trim()) return; // Ignore empty

    // Add user message to history
    const userMessage = { sender: 'user', text: input };
    const updatedMessages = [...messages, userMessage];

    console.log(`Sending to LLM for camera: ${cameraId}`); // ✅ Future use!

    // Fake LLM reply
    const llmReply = {
      sender: 'llm',
      text: "🤖 This is a placeholder LLM response!"
    };

    // Update state
    setMessages([...updatedMessages, llmReply]);
    setInput('');
  };

  // === Enter key support ===
  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      handleSend();
    }
  };

  return (
    <div className="chat-box">
        {/* === Optional camera ID info === */}
        <div className="chat-camera-info">
            Chatting with camera: <strong>{cameraId}</strong>
        </div>
      {/* === Chat messages === */}
      <div className="chat-messages">
        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`chat-message ${msg.sender === 'user' ? 'user' : 'llm'}`}
          >
            <strong>{msg.sender === 'user' ? 'You:' : 'LLM:'}</strong> {msg.text}
          </div>
        ))}
      </div>

      {/* === Input area === */}
      <div className="chat-input">
        <input
          type="text"
          placeholder="Type your message..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <button onClick={handleSend}>Send</button>
      </div>
    </div>
  );
};

export default ChatBox;
