"""
Input event model
"""
from pydantic import BaseModel


class Order(BaseModel):
    """
    Order schema
    """
    id: int
    quantity: int
    description: str


class JinjaTemplateInput(BaseModel):
    """
    Jinja template input schema
    """
    title: str
    content: str