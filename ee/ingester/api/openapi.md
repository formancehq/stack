<!-- Generator: Widdershins v4.0.1 -->

<h1 id="formance-simple-ingester-service-api">Formance Simple ingester Service API v0.1.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="http://localhost:8080">http://localhost:8080</a>

<h1 id="formance-simple-ingester-service-api-connectors">Connectors</h1>

## List connectors

<a id="opIdlistConnectors"></a>

`GET /connectors`

> Example responses

> 200 Response

```json
{
  "cursor": {
    "pageSize": 15,
    "hasMore": false,
    "previous": "YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=",
    "next": "",
    "data": [
      {
        "id": "string",
        "createdAt": "2019-08-24T14:15:22Z"
      }
    ]
  }
}
```

<h3 id="list-connectors-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Connectors list|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="list-connectors-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» cursor|any|false|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|[Cursor](#schemacursor)|false|none|none|
|»»» pageSize|integer(int64)|true|none|none|
|»»» hasMore|boolean|true|none|none|
|»»» previous|string|false|none|none|
|»»» next|string|false|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» data|[allOf]|false|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»»»» *anonymous*|object|false|none|none|
|»»»»» driver|string|true|none|none|
|»»»»» config|object|true|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»»»» *anonymous*|object|false|none|none|
|»»»»» id|string|true|none|none|
|»»»»» createdAt|string(date-time)|true|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## Create connector

<a id="opIdcreateConnector"></a>

`POST /connectors`

> Body parameter

```json
{
  "driver": "string",
  "config": {}
}
```

<h3 id="create-connector-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[ConnectorConfiguration](#schemaconnectorconfiguration)|false|none|

> Example responses

> 201 Response

```json
{
  "data": {
    "id": "string",
    "createdAt": "2019-08-24T14:15:22Z"
  }
}
```

<h3 id="create-connector-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created pipeline|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="create-connector-responseschema">Response Schema</h3>

Status Code **201**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» data|any|true|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» driver|string|true|none|none|
|»»» config|object|true|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» id|string|true|none|none|
|»»» createdAt|string(date-time)|true|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## Get connector state

<a id="opIdgetConnectorState"></a>

`GET /connectors/{connectorID}`

<h3 id="get-connector-state-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|connectorID|path|string|true|The connector id|

> Example responses

> 200 Response

```json
{
  "data": {
    "id": "string",
    "createdAt": "2019-08-24T14:15:22Z"
  }
}
```

<h3 id="get-connector-state-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Connector information|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="get-connector-state-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» data|any|true|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» driver|string|true|none|none|
|»»» config|object|true|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» id|string|true|none|none|
|»»» createdAt|string(date-time)|true|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## Delete connector

<a id="opIddeleteConnector"></a>

`DELETE /connectors/{connectorID}`

<h3 id="delete-connector-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|connectorID|path|string|true|The connector id|

> Example responses

> default Response

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}
```

<h3 id="delete-connector-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Connector deleted|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="formance-simple-ingester-service-api-pipelines">Pipelines</h1>

## List pipelines

<a id="opIdlistPipelines"></a>

`GET /pipelines`

> Example responses

> 200 Response

```json
{
  "cursor": {
    "pageSize": 15,
    "hasMore": false,
    "previous": "YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=",
    "next": "",
    "data": [
      {
        "id": "string",
        "state": {
          "label": "STOP",
          "cursor": "string",
          "previousState": {}
        },
        "createdAt": "2019-08-24T14:15:22Z"
      }
    ]
  }
}
```

<h3 id="list-pipelines-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Pipelines list|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="list-pipelines-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» cursor|any|false|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|[Cursor](#schemacursor)|false|none|none|
|»»» pageSize|integer(int64)|true|none|none|
|»»» hasMore|boolean|true|none|none|
|»»» previous|string|false|none|none|
|»»» next|string|false|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» data|[allOf]|false|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»»»» *anonymous*|object|false|none|none|
|»»»»» module|string|true|none|none|
|»»»»» connectorID|string|true|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»»»» *anonymous*|object|false|none|none|
|»»»»» id|string|true|none|none|
|»»»»» state|[State](#schemastate)|true|none|none|
|»»»»»» label|string|true|none|none|
|»»»»»» cursor|string|false|none|none|
|»»»»»» previousState|[State](#schemastate)|false|none|none|
|»»»»» createdAt|string(date-time)|true|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|label|STOP|
|label|PAUSE|
|label|INIT|
|label|READY|

<aside class="success">
This operation does not require authentication
</aside>

## Create pipeline

<a id="opIdcreatePipeline"></a>

`POST /pipelines`

> Body parameter

```json
{
  "module": "string",
  "connectorID": "string"
}
```

<h3 id="create-pipeline-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[CreatePipelineRequest](#schemacreatepipelinerequest)|false|none|

> Example responses

> 201 Response

```json
{
  "data": {
    "id": "string",
    "state": {
      "label": "STOP",
      "cursor": "string",
      "previousState": {}
    },
    "createdAt": "2019-08-24T14:15:22Z"
  }
}
```

<h3 id="create-pipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created ipeline|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="create-pipeline-responseschema">Response Schema</h3>

Status Code **201**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» data|any|true|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» module|string|true|none|none|
|»»» connectorID|string|true|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» id|string|true|none|none|
|»»» state|[State](#schemastate)|true|none|none|
|»»»» label|string|true|none|none|
|»»»» cursor|string|false|none|none|
|»»»» previousState|[State](#schemastate)|false|none|none|
|»»» createdAt|string(date-time)|true|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|label|STOP|
|label|PAUSE|
|label|INIT|
|label|READY|

<aside class="success">
This operation does not require authentication
</aside>

## Get pipeline state

<a id="opIdgetPipelineState"></a>

`GET /pipelines/{pipelineID}`

<h3 id="get-pipeline-state-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pipelineID|path|string|true|The pipeline id|

> Example responses

> 200 Response

```json
{
  "data": {
    "id": "string",
    "state": {
      "label": "STOP",
      "cursor": "string",
      "previousState": {}
    },
    "createdAt": "2019-08-24T14:15:22Z"
  }
}
```

<h3 id="get-pipeline-state-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Pipeline information|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="get-pipeline-state-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» data|any|true|none|none|

*allOf*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» module|string|true|none|none|
|»»» connectorID|string|true|none|none|

*and*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|»» *anonymous*|object|false|none|none|
|»»» id|string|true|none|none|
|»»» state|[State](#schemastate)|true|none|none|
|»»»» label|string|true|none|none|
|»»»» cursor|string|false|none|none|
|»»»» previousState|[State](#schemastate)|false|none|none|
|»»» createdAt|string(date-time)|true|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|label|STOP|
|label|PAUSE|
|label|INIT|
|label|READY|

<aside class="success">
This operation does not require authentication
</aside>

## Delete pipeline

<a id="opIddeletePipeline"></a>

`DELETE /pipelines/{pipelineID}`

<h3 id="delete-pipeline-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pipelineID|path|string|true|The pipeline id|

> Example responses

> default Response

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}
```

<h3 id="delete-pipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Pipeline deleted|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## Reset pipeline

<a id="opIdresetPipeline"></a>

`POST /pipelines/{pipelineID}/reset`

<h3 id="reset-pipeline-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pipelineID|path|string|true|The pipeline id|

> Example responses

> default Response

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}
```

<h3 id="reset-pipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline reset|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## Pause pipeline

<a id="opIdpausePipeline"></a>

`POST /pipelines/{pipelineID}/pause`

<h3 id="pause-pipeline-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pipelineID|path|string|true|The pipeline id|

> Example responses

> default Response

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}
```

<h3 id="pause-pipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline pause|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## Resume pipeline

<a id="opIdresumePipeline"></a>

`POST /pipelines/{pipelineID}/resume`

<h3 id="resume-pipeline-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pipelineID|path|string|true|The pipeline id|

> Example responses

> default Response

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}
```

<h3 id="resume-pipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline resumed|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## Start pipeline

<a id="opIdstartPipeline"></a>

`POST /pipelines/{pipelineID}/start`

<h3 id="start-pipeline-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pipelineID|path|string|true|The pipeline id|

> Example responses

> default Response

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}
```

<h3 id="start-pipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline started|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## Stop pipeline

<a id="opIdstopPipeline"></a>

`POST /pipelines/{pipelineID}/stop`

<h3 id="stop-pipeline-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pipelineID|path|string|true|The pipeline id|

> Example responses

> default Response

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}
```

<h3 id="stop-pipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline stopped|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_CreatePipelineRequest">CreatePipelineRequest</h2>
<!-- backwards compatibility -->
<a id="schemacreatepipelinerequest"></a>
<a id="schema_CreatePipelineRequest"></a>
<a id="tocScreatepipelinerequest"></a>
<a id="tocscreatepipelinerequest"></a>

```json
{
  "module": "string",
  "connectorID": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|module|string|true|none|none|
|connectorID|string|true|none|none|

<h2 id="tocS_CreateConnectorRequest">CreateConnectorRequest</h2>
<!-- backwards compatibility -->
<a id="schemacreateconnectorrequest"></a>
<a id="schema_CreateConnectorRequest"></a>
<a id="tocScreateconnectorrequest"></a>
<a id="tocscreateconnectorrequest"></a>

```json
{
  "driver": "string",
  "config": {}
}

```

### Properties

*None*

<h2 id="tocS_PipelineConfiguration">PipelineConfiguration</h2>
<!-- backwards compatibility -->
<a id="schemapipelineconfiguration"></a>
<a id="schema_PipelineConfiguration"></a>
<a id="tocSpipelineconfiguration"></a>
<a id="tocspipelineconfiguration"></a>

```json
{
  "module": "string",
  "connectorID": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|module|string|true|none|none|
|connectorID|string|true|none|none|

<h2 id="tocS_ConnectorConfiguration">ConnectorConfiguration</h2>
<!-- backwards compatibility -->
<a id="schemaconnectorconfiguration"></a>
<a id="schema_ConnectorConfiguration"></a>
<a id="tocSconnectorconfiguration"></a>
<a id="tocsconnectorconfiguration"></a>

```json
{
  "driver": "string",
  "config": {}
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|driver|string|true|none|none|
|config|object|true|none|none|

<h2 id="tocS_Connector">Connector</h2>
<!-- backwards compatibility -->
<a id="schemaconnector"></a>
<a id="schema_Connector"></a>
<a id="tocSconnector"></a>
<a id="tocsconnector"></a>

```json
{
  "id": "string",
  "createdAt": "2019-08-24T14:15:22Z"
}

```

### Properties

allOf

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[ConnectorConfiguration](#schemaconnectorconfiguration)|false|none|none|

and

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|object|false|none|none|
|» id|string|true|none|none|
|» createdAt|string(date-time)|true|none|none|

<h2 id="tocS_Pipeline">Pipeline</h2>
<!-- backwards compatibility -->
<a id="schemapipeline"></a>
<a id="schema_Pipeline"></a>
<a id="tocSpipeline"></a>
<a id="tocspipeline"></a>

```json
{
  "id": "string",
  "state": {
    "label": "STOP",
    "cursor": "string",
    "previousState": {}
  },
  "createdAt": "2019-08-24T14:15:22Z"
}

```

### Properties

allOf

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[PipelineConfiguration](#schemapipelineconfiguration)|false|none|none|

and

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|object|false|none|none|
|» id|string|true|none|none|
|» state|[State](#schemastate)|true|none|none|
|» createdAt|string(date-time)|true|none|none|

<h2 id="tocS_State">State</h2>
<!-- backwards compatibility -->
<a id="schemastate"></a>
<a id="schema_State"></a>
<a id="tocSstate"></a>
<a id="tocsstate"></a>

```json
{
  "label": "STOP",
  "cursor": "string",
  "previousState": {
    "label": "STOP",
    "cursor": "string",
    "previousState": {}
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|label|string|true|none|none|
|cursor|string|false|none|none|
|previousState|[State](#schemastate)|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|label|STOP|
|label|PAUSE|
|label|INIT|
|label|READY|

<h2 id="tocS_Error">Error</h2>
<!-- backwards compatibility -->
<a id="schemaerror"></a>
<a id="schema_Error"></a>
<a id="tocSerror"></a>
<a id="tocserror"></a>

```json
{
  "errorCode": "string",
  "errorMessage": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|errorCode|string|true|none|none|
|errorMessage|string|true|none|none|

<h2 id="tocS_Cursor">Cursor</h2>
<!-- backwards compatibility -->
<a id="schemacursor"></a>
<a id="schema_Cursor"></a>
<a id="tocScursor"></a>
<a id="tocscursor"></a>

```json
{
  "pageSize": 15,
  "hasMore": false,
  "previous": "YXVsdCBhbmQgYSBtYXhpbXVtIG1heF9yZXN1bHRzLol=",
  "next": ""
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pageSize|integer(int64)|true|none|none|
|hasMore|boolean|true|none|none|
|previous|string|false|none|none|
|next|string|false|none|none|

