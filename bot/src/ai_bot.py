from fastapi import APIRouter
from src.models import Message
from dotenv import load_dotenv
import os
import httpx

ai_bot_router = APIRouter()

load_dotenv()

OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")
OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"


@ai_bot_router.post("/ai_bot_chat")
async def aiBotChat(msg: Message):
    user_text = msg.text

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json"
    }

    payload = {
        "model": "meta-llama/llama-3.2-11b-vision-instruct:free", 
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
