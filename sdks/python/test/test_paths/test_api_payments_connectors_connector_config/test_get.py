# coding: utf-8

"""


    Generated by: https://openapi-generator.tech
"""

import unittest
from unittest.mock import patch

import urllib3

import FormanceHQ
from FormanceHQ.paths.api_payments_connectors_connector_config import get  # noqa: E501
from FormanceHQ import configuration, schemas, api_client

from .. import ApiTestMixin


class TestApiPaymentsConnectorsConnectorConfig(ApiTestMixin, unittest.TestCase):
    """
    ApiPaymentsConnectorsConnectorConfig unit test stubs
        Read the config of a connector  # noqa: E501
    """
    _configuration = configuration.Configuration()

    def setUp(self):
        used_api_client = api_client.ApiClient(configuration=self._configuration)
        self.api = get.ApiForget(api_client=used_api_client)  # noqa: E501

    def tearDown(self):
        pass

    response_status = 200




if __name__ == '__main__':
    unittest.main()
