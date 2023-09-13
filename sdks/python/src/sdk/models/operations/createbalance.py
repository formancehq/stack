"""Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT."""

from __future__ import annotations
import dataclasses
import requests as requests_http
from ..shared import createbalancerequest as shared_createbalancerequest
from ..shared import createbalanceresponse as shared_createbalanceresponse
from ..shared import walletserrorresponse as shared_walletserrorresponse
from typing import Optional


@dataclasses.dataclass
class CreateBalanceRequest:
    
    id: str = dataclasses.field(metadata={'path_param': { 'field_name': 'id', 'style': 'simple', 'explode': False }})
    create_balance_request: Optional[shared_createbalancerequest.CreateBalanceRequest] = dataclasses.field(default=None, metadata={'request': { 'media_type': 'application/json' }})
    

@dataclasses.dataclass
class CreateBalanceResponse:
    
    content_type: str = dataclasses.field()
    status_code: int = dataclasses.field()
    create_balance_response: Optional[shared_createbalanceresponse.CreateBalanceResponse] = dataclasses.field(default=None)
    r"""Created balance"""
    raw_response: Optional[requests_http.Response] = dataclasses.field(default=None)
    wallets_error_response: Optional[shared_walletserrorresponse.WalletsErrorResponse] = dataclasses.field(default=None)
    r"""Error"""
    