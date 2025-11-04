from pydantic import BaseModel, Field


class UserSchema(BaseModel):
    id: int = Field(...)
    name: str  = Field(...)
    password: str = Field(...)


    class Config:
        from_attributes = True


class UserSchemaAdd(BaseModel):
    name: str
    password: str