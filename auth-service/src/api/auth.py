from typing import Annotated
import jwt
from fastapi import APIRouter, Body
from schemas.users import UserSchemaAdd

# from api.dependencies import users_service
# from schemas.users import UserSchemaAdd
# from services.users import UsersService

router = APIRouter(
    prefix="/auth",
    tags=["Auth"],
)

@router.post("/add_user")
async def add_user(user: UserSchemaAdd = Body(...)):
    print("user", user)
    encoded_jwt = jwt.encode(user.dict(), "secret_key", algorithm='HS256')
    return {"access_token": encoded_jwt}