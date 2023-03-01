# coding: utf-8

"""
    Formance Stack API

    Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions />   # noqa: E501

    The version of the OpenAPI document: v1.0.20230301
    Contact: support@formance.com
    Generated by: https://openapi-generator.tech
"""

from FormanceHQ.paths.api_payments_connectors_stripe_transfers.post import ConnectorsStripeTransfer
from FormanceHQ.paths.api_payments_connectors_connector_transfers.post import ConnectorsTransfer
from FormanceHQ.paths.api_payments_connectors_connector_tasks_task_id.get import GetConnectorTask
from FormanceHQ.paths.api_payments_payments_payment_id.get import GetPayment
from FormanceHQ.paths.api_payments_connectors_connector.post import InstallConnector
from FormanceHQ.paths.api_payments_connectors.get import ListAllConnectors
from FormanceHQ.paths.api_payments_connectors_configs.get import ListConfigsAvailableConnectors
from FormanceHQ.paths.api_payments_connectors_connector_tasks.get import ListConnectorTasks
from FormanceHQ.paths.api_payments_connectors_connector_transfers.get import ListConnectorsTransfers
from FormanceHQ.paths.api_payments_payments.get import ListPayments
from FormanceHQ.paths.api_payments_accounts.get import PaymentslistAccounts
from FormanceHQ.paths.api_payments_connectors_connector_config.get import ReadConnectorConfig
from FormanceHQ.paths.api_payments_connectors_connector_reset.post import ResetConnector
from FormanceHQ.paths.api_payments_connectors_connector.delete import UninstallConnector
from FormanceHQ.paths.api_payments_payments_payment_id_metadata.patch import UpdateMetadata


class PaymentsApi(
    ConnectorsStripeTransfer,
    ConnectorsTransfer,
    GetConnectorTask,
    GetPayment,
    InstallConnector,
    ListAllConnectors,
    ListConfigsAvailableConnectors,
    ListConnectorTasks,
    ListConnectorsTransfers,
    ListPayments,
    PaymentslistAccounts,
    ReadConnectorConfig,
    ResetConnector,
    UninstallConnector,
    UpdateMetadata,
):
    """NOTE: This class is auto generated by OpenAPI Generator
    Ref: https://openapi-generator.tech

    Do not edit the class manually.
    """
    pass
