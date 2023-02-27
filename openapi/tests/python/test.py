import Formance

configuration = Formance.Configuration(
    host="http://localhost:3068"
)

with Formance.ApiClient(configuration) as api_client:
    print("TODO: Actually just SDK import + compile")
