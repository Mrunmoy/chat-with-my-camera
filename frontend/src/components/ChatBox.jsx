import React, { useState } from 'react';
import '../styles/ChatBox.css';

/**
 * ChatBox
 *
 * - Lets the user type messages to your backend LLM.
 * - Sends { camera_id, message } to `/chat`.
 * - Shows the full conversation history: user & bot.
 */
const ChatBox = ({ cameraId }) => {
  // === Local state ===
  const [messages, setMessages] = useState([]); // [{ sender: 'user'|'llm', text: '...' }]
  const [input, setInput] = useState('');

  // === Send user input to backend ===
  const handleSend = async () => {
    if (!input.trim()) return; // Ignore empty

    // Add user message to history
    const userMessage = { sender: 'user', text: input };
    const updatedMessages = [...messages, userMessage];
    setMessages(updatedMessages);

    console.log(`Sending to /chat for camera: ${cameraId}`);

    try {
      // ✅ Send POST to your backend
      const res = await fetch('http://localhost:8080/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          camera_id: cameraId,
          message: input
        }),
      });

      const data = await res.json();

      // Add LLM response to chat history
      const llmReply = {
        sender: 'llm',
        text: data.answer || '🤖 No response from LLM!'
      };

      setMessages([...updatedMessages, llmReply]);
    } catch (err) {
      console.error('Error calling LLM:', err);

      // Add fallback error message
      const errorReply = {
        sender: 'llm',
        text: "🤖 Oops! Couldn't reach the backend."
      };
      setMessages([...updatedMessages, errorReply]);
    }

    // Clear input box
    setInput('');
  };

  // === Enter key ===
  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      handleSend();
    }
  };

  return (
    <div className="chat-box">
      {/* === Optional camera ID banner === */}
      <div className="chat-camera-info">
        Chatting with camera: <strong>{cameraId}</strong>
      </div>

      {/* === Chat conversation === */}
      <div className="chat-messages">
        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`chat-message ${msg.sender}`}
          >
            <strong>{msg.sender === 'user' ? 'You:' : 'LLM:'}</strong> {msg.text}
          </div>
        ))}
      </div>

      {/* === Input area === */}
      <div className="chat-input">
        <input
          type="text"
          placeholder="Type your question..."
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
