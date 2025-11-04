import uvicorn
from fastapi import FastAPI

app = FastAPI(
    title="Упрощенный аналог Jira/Asana"
)

@app.post("/add_user")
async def add_user(body: dict):
    return {"access_token": body}

if __name__ == "__main__":
    uvicorn.run(app="main:app", reload=True)
