# LLM


![LLM](https://img.shields.io/badge/LLM-Ollama-orange?logo=OpenAI&logoColor=white)

## Overview

This folder documents the **Local LLM integration** for `chat-with-my-camera`.

We use **Ollama** (or similar) running locally. The backend uses a **2-step prompt chain**:

- Extract relevant objects from the user’s question.
- Query the timeline for matches and generate a final natural-language answer.

Everything stays local — no cloud API needed.

---

## How It Works

### Extraction

- `/chat` endpoint calls Ollama to extract the main labels/objects from the question.
- Example prompt:
  > *"User question: When did you last see a car? Extract the main object(s) as JSON array."*

### Timeline Lookup

- The backend queries SQLite for the last detection matching that object.
- If found, a context string is built: timestamps, labels, etc.
- If not found, a fallback answer is used.

### Final Prompt

- The context is inserted into a final user-friendly prompt:
  ```text
  Camera: garage_webcam
  \nDetection context:\n- Last detection: 2025-07-11 Labels: ["car"]\n\nUser question: When did you last see a car?
  ```
- Ollama returns the final chat response.

### JSON Structure

**Request:**

```json
{
  "camera_id": "garage_webcam",
  "message": "When did you last see a car?"
}
```

**Response:**

```json
{
  "answer": "I last saw a car at 2pm yesterday."
}
```

---

## How to Run Ollama

```bash
ollama serve
ollama pull llama3
```

Runs on `localhost:11434`. Your backend must be able to POST to it.

---

## Tips

Test your prompt locally with curl:

```bash
curl http://localhost:11434/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3",
    "messages": [
      {"role": "system", "content": "You are a helpful camera assistant."},
      {"role": "user", "content": "When did you last see a car?"}
    ]
  }'
```

Keep system messages clear and consistent. Use JSON arrays for extraction for easy parsing.
