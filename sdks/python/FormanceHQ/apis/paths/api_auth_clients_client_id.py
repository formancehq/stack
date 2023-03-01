from FormanceHQ.paths.api_auth_clients_client_id.get import ApiForget
from FormanceHQ.paths.api_auth_clients_client_id.put import ApiForput
from FormanceHQ.paths.api_auth_clients_client_id.delete import ApiFordelete


class ApiAuthClientsClientId(
    ApiForget,
    ApiForput,
    ApiFordelete,
):
    pass
