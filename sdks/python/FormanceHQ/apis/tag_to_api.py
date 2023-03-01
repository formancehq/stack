import typing_extensions

from FormanceHQ.apis.tags import TagValues
from FormanceHQ.apis.tags.accounts_api import AccountsApi
from FormanceHQ.apis.tags.balances_api import BalancesApi
from FormanceHQ.apis.tags.clients_api import ClientsApi
from FormanceHQ.apis.tags.ledger_api import LedgerApi
from FormanceHQ.apis.tags.logs_api import LogsApi
from FormanceHQ.apis.tags.mapping_api import MappingApi
from FormanceHQ.apis.tags.orchestration_api import OrchestrationApi
from FormanceHQ.apis.tags.payments_api import PaymentsApi
from FormanceHQ.apis.tags.scopes_api import ScopesApi
from FormanceHQ.apis.tags.script_api import ScriptApi
from FormanceHQ.apis.tags.search_api import SearchApi
from FormanceHQ.apis.tags.server_api import ServerApi
from FormanceHQ.apis.tags.stats_api import StatsApi
from FormanceHQ.apis.tags.transactions_api import TransactionsApi
from FormanceHQ.apis.tags.users_api import UsersApi
from FormanceHQ.apis.tags.wallets_api import WalletsApi
from FormanceHQ.apis.tags.webhooks_api import WebhooksApi
from FormanceHQ.apis.tags.default_api import DefaultApi

TagToApi = typing_extensions.TypedDict(
    'TagToApi',
    {
        TagValues.ACCOUNTS: AccountsApi,
        TagValues.BALANCES: BalancesApi,
        TagValues.CLIENTS: ClientsApi,
        TagValues.LEDGER: LedgerApi,
        TagValues.LOGS: LogsApi,
        TagValues.MAPPING: MappingApi,
        TagValues.ORCHESTRATION: OrchestrationApi,
        TagValues.PAYMENTS: PaymentsApi,
        TagValues.SCOPES: ScopesApi,
        TagValues.SCRIPT: ScriptApi,
        TagValues.SEARCH: SearchApi,
        TagValues.SERVER: ServerApi,
        TagValues.STATS: StatsApi,
        TagValues.TRANSACTIONS: TransactionsApi,
        TagValues.USERS: UsersApi,
        TagValues.WALLETS: WalletsApi,
        TagValues.WEBHOOKS: WebhooksApi,
        TagValues.DEFAULT: DefaultApi,
    }
)

tag_to_api = TagToApi(
    {
        TagValues.ACCOUNTS: AccountsApi,
        TagValues.BALANCES: BalancesApi,
        TagValues.CLIENTS: ClientsApi,
        TagValues.LEDGER: LedgerApi,
        TagValues.LOGS: LogsApi,
        TagValues.MAPPING: MappingApi,
        TagValues.ORCHESTRATION: OrchestrationApi,
        TagValues.PAYMENTS: PaymentsApi,
        TagValues.SCOPES: ScopesApi,
        TagValues.SCRIPT: ScriptApi,
        TagValues.SEARCH: SearchApi,
        TagValues.SERVER: ServerApi,
        TagValues.STATS: StatsApi,
        TagValues.TRANSACTIONS: TransactionsApi,
        TagValues.USERS: UsersApi,
        TagValues.WALLETS: WalletsApi,
        TagValues.WEBHOOKS: WebhooksApi,
        TagValues.DEFAULT: DefaultApi,
    }
)
