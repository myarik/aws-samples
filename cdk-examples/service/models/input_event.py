"""
Input event model
"""

from pydantic import BaseModel, HttpUrl


class ServiceUrlInput(BaseModel):
    """
    Order schema
    """

    url: HttpUrl


class JinjaTemplateInput(BaseModel):
    """
    Jinja template input schema
    """

    title: str
    content: str
