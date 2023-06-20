"""Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT."""

from __future__ import annotations
import dataclasses
import requests as requests_http
from ..shared import error as shared_error
from ..shared import listrunsresponse as shared_listrunsresponse
from typing import Optional


@dataclasses.dataclass
class ListInstancesRequest:
    
    running: Optional[bool] = dataclasses.field(default=None, metadata={'query_param': { 'field_name': 'running', 'style': 'form', 'explode': True }})
    r"""Filter running instances"""
    workflow_id: Optional[str] = dataclasses.field(default=None, metadata={'query_param': { 'field_name': 'workflowID', 'style': 'form', 'explode': True }})
    r"""A workflow id"""
    

@dataclasses.dataclass
class ListInstancesResponse:
    
    content_type: str = dataclasses.field()
    status_code: int = dataclasses.field()
    error: Optional[shared_error.Error] = dataclasses.field(default=None)
    r"""General error"""
    list_runs_response: Optional[shared_listrunsresponse.ListRunsResponse] = dataclasses.field(default=None)
    r"""List of workflow instances"""
    raw_response: Optional[requests_http.Response] = dataclasses.field(default=None)
    