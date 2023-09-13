"""Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT."""

from __future__ import annotations
import dataclasses
from ..shared import stagesenddestinationaccount as shared_stagesenddestinationaccount
from ..shared import stagesenddestinationpayment as shared_stagesenddestinationpayment
from ..shared import stagesenddestinationwallet as shared_stagesenddestinationwallet
from dataclasses_json import Undefined, dataclass_json
from sdk import utils
from typing import Optional


@dataclass_json(undefined=Undefined.EXCLUDE)
@dataclasses.dataclass
class StageSendDestination:
    
    account: Optional[shared_stagesenddestinationaccount.StageSendDestinationAccount] = dataclasses.field(default=None, metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('account'), 'exclude': lambda f: f is None }})
    payment: Optional[shared_stagesenddestinationpayment.StageSendDestinationPayment] = dataclasses.field(default=None, metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('payment'), 'exclude': lambda f: f is None }})
    wallet: Optional[shared_stagesenddestinationwallet.StageSendDestinationWallet] = dataclasses.field(default=None, metadata={'dataclasses_json': { 'letter_case': utils.get_field_name('wallet'), 'exclude': lambda f: f is None }})
    