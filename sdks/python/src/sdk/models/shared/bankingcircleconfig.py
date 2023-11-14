"""Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT."""

from __future__ import annotations
import dataclasses
from dataclasses_json import Undefined, dataclass_json
from sdk import utils
from typing import Optional


@dataclass_json(undefined=Undefined.EXCLUDE)
@dataclasses.dataclass
class BankingCircleConfig:
    
    authorization_endpoint: str = dataclasses.field(metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('authorizationEndpoint') }})
    endpoint: str = dataclasses.field(metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('endpoint') }})
    name: str = dataclasses.field(metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('name') }})
    password: str = dataclasses.field(metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('password') }})
    user_certificate: str = dataclasses.field(metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('userCertificate') }})
    user_certificate_key: str = dataclasses.field(metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('userCertificateKey') }})
    username: str = dataclasses.field(metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('username') }})
    polling_period: Optional[str] = dataclasses.field(default=None, metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('pollingPeriod'), 'exclude': lambda f: f is None }})
    r"""The frequency at which the connector will try to fetch new BalanceTransaction objects from Banking Circle API."""
    