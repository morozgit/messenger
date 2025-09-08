from fastapi import FastAPI

import uvicorn
from src.ai_bot import ai_bot_router 
from src.my_ai_bot import my_ai_bot_router 


app = FastAPI()

app.include_router(ai_bot_router)
app.include_router(my_ai_bot_router)


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8085)