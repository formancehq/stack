

# Transaction


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**timestamp** | **OffsetDateTime** |  |  |
|**postings** | [**List&lt;Posting&gt;**](Posting.md) |  |  |
|**reference** | **String** |  |  [optional] |
|**metadata** | **Map&lt;String, String&gt;** |  |  |
|**txid** | **UUID** |  |  |
|**preCommitVolumes** | **Map&lt;String, Map&lt;String, Volume&gt;&gt;** |  |  [optional] |
|**postCommitVolumes** | **Map&lt;String, Map&lt;String, Volume&gt;&gt;** |  |  [optional] |



