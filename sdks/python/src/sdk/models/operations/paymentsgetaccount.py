"""Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT."""

from __future__ import annotations
import dataclasses
import requests as requests_http
from ..shared import paymentsaccountresponse as shared_paymentsaccountresponse
from typing import Optional



@dataclasses.dataclass
class PaymentsgetAccountRequest:
    account_id: str = dataclasses.field(metadata={'path_param': { 'field_name': 'accountId', 'style': 'simple', 'explode': False }})
    r"""The account ID."""
    




@dataclasses.dataclass
class PaymentsgetAccountResponse:
    content_type: str = dataclasses.field()
    status_code: int = dataclasses.field()
    payments_account_response: Optional[shared_paymentsaccountresponse.PaymentsAccountResponse] = dataclasses.field(default=None)
    r"""OK"""
    raw_response: Optional[requests_http.Response] = dataclasses.field(default=None)
    

