import FormanceHQ

configuration = FormanceHQ.Configuration(
    host="http://localhost:3068"
)

with FormanceHQ.ApiClient(configuration) as api_client:
    print("TODO: Actually just SDK import + compile")
