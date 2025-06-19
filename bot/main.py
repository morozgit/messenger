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


@app.post("/ai_bot_chat")
async def aiBotChat(msg: Message):
    user_text = msg.text

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json"
    }

    payload = {
        "model": "opengvlab/internvl3-14b:free", 
        "messages": [
            {"role": "system", "content": "Ты дружелюбный помощник в чате. Отвечай кратко, но информативно, до 5-6 предложений."},
            {"role": "user", "content": user_text}
        ],
        "max_tokens": 512
    }

    try:
        async with httpx.AsyncClient() as client:
            resp = await client.post(OPENROUTER_URL, json=payload, headers=headers)
            resp.raise_for_status()
            data = resp.json()
            return {"reply": data["choices"][0]["message"]["content"].strip()}
    except Exception as e:
        return {"reply": f"Ошибка: {str(e)}"}



@app.post("/my_ai_bot_chat")
async def myAiBotChat(msg: Message):
    # user_text = msg.text
    return {"reply": 'Бот в разработке'}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8085)