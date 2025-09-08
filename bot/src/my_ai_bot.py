from src.models import Message
from fastapi import APIRouter

my_ai_bot_router = APIRouter()

@my_ai_bot_router.post("/my_ai_bot_chat")
async def myAiBotChat(msg: Message):
    # user_text = msg.text
    return {"reply": 'Бот в разработке'}