"""Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT."""

from __future__ import annotations
import dataclasses
import requests as requests_http
from ..shared import accountscursorresponse as shared_accountscursorresponse
from ..shared import errorresponse as shared_errorresponse
from datetime import datetime
from typing import Any, Optional



@dataclasses.dataclass
class ListAccountsRequest:
    ledger: str = dataclasses.field(metadata={'path_param': { 'field_name': 'ledger', 'style': 'simple', 'explode': False }})
    r"""Name of the ledger."""
    cursor: Optional[str] = dataclasses.field(default=None, metadata={'query_param': { 'field_name': 'cursor', 'style': 'form', 'explode': True }})
    r"""Parameter used in pagination requests. Maximum page size is set to 15.
    Set to the value of next for the next page of results.
    Set to the value of previous for the previous page of results.
    No other parameters can be set when this parameter is set.
    """
    expand: Optional[str] = dataclasses.field(default=None, metadata={'query_param': { 'field_name': 'expand', 'style': 'form', 'explode': True }})
    page_size: Optional[int] = dataclasses.field(default=None, metadata={'query_param': { 'field_name': 'pageSize', 'style': 'form', 'explode': True }})
    r"""The maximum number of results to return per page."""
    pit: Optional[datetime] = dataclasses.field(default=None, metadata={'query_param': { 'field_name': 'pit', 'style': 'form', 'explode': True }})
    request_body: Optional[dict[str, Any]] = dataclasses.field(default=None, metadata={'request': { 'media_type': 'application/json' }})
    




@dataclasses.dataclass
class ListAccountsResponse:
    content_type: str = dataclasses.field()
    status_code: int = dataclasses.field()
    accounts_cursor_response: Optional[shared_accountscursorresponse.AccountsCursorResponse] = dataclasses.field(default=None)
    r"""OK"""
    error_response: Optional[shared_errorresponse.ErrorResponse] = dataclasses.field(default=None)
    r"""Error"""
    raw_response: Optional[requests_http.Response] = dataclasses.field(default=None)
    

