from fastapi import FastAPI
from pydantic import BaseModel
import os
import httpx
from dotenv import load_dotenv
import uvicorn

load_dotenv()

OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")
OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"

app = FastAPI()


class Message(BaseModel):
    text: str


@app.post("/chat")
async def chat(msg: Message):
    user_text = msg.text

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json"
    }

    payload = {
        "model": "moonshotai/kimi-dev-72b:free", 
        "messages": [
            {"role": "system", "content": "Ты дружелюбный помощник в чате."},
            {"role": "user", "content": user_text}
        ],
        "max_tokens": 100
    }

    try:
        async with httpx.AsyncClient() as client:
            resp = await client.post(OPENROUTER_URL, json=payload, headers=headers)
            resp.raise_for_status()
            data = resp.json()
            return {"reply": data["choices"][0]["message"]["content"].strip()}
    except Exception as e:
        return {"reply": f"Ошибка: {str(e)}"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8085)