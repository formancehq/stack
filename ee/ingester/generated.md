---
title: Formance Simple ingester Service API v0.1.0
language_tabs:
  - http: HTTP
language_clients:
  - http: ""
toc_footers: []
includes: []
search: false
highlight_theme: darkula
headingLevel: 2

---

<!-- Generator: Widdershins v4.0.1 -->

<h1 id="formance-simple-ingester-service-api">Formance Simple ingester Service API v0.1.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="http://localhost:8080">http://localhost:8080</a>

<h1 id="formance-simple-ingester-service-api-connectors">Connectors</h1>

## listConnectors

<a id="opIdlistConnectors"></a>

> Code samples

```http
GET http://localhost:8080/connectors HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`GET /connectors`

*List connectors*

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

<h3 id="listconnectors-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Connectors list|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="listconnectors-responseschema">Response Schema</h3>

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

## createConnector

<a id="opIdcreateConnector"></a>

> Code samples

```http
POST http://localhost:8080/connectors HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Accept: application/json

```

`POST /connectors`

*Create connector*

> Body parameter

```json
{
  "driver": "string",
  "config": {}
}
```

<h3 id="createconnector-parameters">Parameters</h3>

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

<h3 id="createconnector-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created pipeline|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="createconnector-responseschema">Response Schema</h3>

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

## getConnectorState

<a id="opIdgetConnectorState"></a>

> Code samples

```http
GET http://localhost:8080/connectors/{connectorID} HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`GET /connectors/{connectorID}`

*Get connector state*

<h3 id="getconnectorstate-parameters">Parameters</h3>

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

<h3 id="getconnectorstate-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Connector information|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="getconnectorstate-responseschema">Response Schema</h3>

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

## deleteConnector

<a id="opIddeleteConnector"></a>

> Code samples

```http
DELETE http://localhost:8080/connectors/{connectorID} HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`DELETE /connectors/{connectorID}`

*Delete connector*

<h3 id="deleteconnector-parameters">Parameters</h3>

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

<h3 id="deleteconnector-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Connector deleted|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="formance-simple-ingester-service-api-pipelines">Pipelines</h1>

## listPipelines

<a id="opIdlistPipelines"></a>

> Code samples

```http
GET http://localhost:8080/pipelines HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`GET /pipelines`

*List pipelines*

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

<h3 id="listpipelines-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Pipelines list|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="listpipelines-responseschema">Response Schema</h3>

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

## createPipeline

<a id="opIdcreatePipeline"></a>

> Code samples

```http
POST http://localhost:8080/pipelines HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Accept: application/json

```

`POST /pipelines`

*Create pipeline*

> Body parameter

```json
{
  "module": "string",
  "connectorID": "string"
}
```

<h3 id="createpipeline-parameters">Parameters</h3>

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

<h3 id="createpipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created ipeline|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="createpipeline-responseschema">Response Schema</h3>

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

## getPipelineState

<a id="opIdgetPipelineState"></a>

> Code samples

```http
GET http://localhost:8080/pipelines/{pipelineID} HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`GET /pipelines/{pipelineID}`

*Get pipeline state*

<h3 id="getpipelinestate-parameters">Parameters</h3>

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

<h3 id="getpipelinestate-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Pipeline information|Inline|
|default|Default|General error|[Error](#schemaerror)|

<h3 id="getpipelinestate-responseschema">Response Schema</h3>

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

## deletePipeline

<a id="opIddeletePipeline"></a>

> Code samples

```http
DELETE http://localhost:8080/pipelines/{pipelineID} HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`DELETE /pipelines/{pipelineID}`

*Delete pipeline*

<h3 id="deletepipeline-parameters">Parameters</h3>

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

<h3 id="deletepipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Pipeline deleted|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## resetPipeline

<a id="opIdresetPipeline"></a>

> Code samples

```http
POST http://localhost:8080/pipelines/{pipelineID}/reset HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`POST /pipelines/{pipelineID}/reset`

*Reset pipeline*

<h3 id="resetpipeline-parameters">Parameters</h3>

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

<h3 id="resetpipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline reset|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## pausePipeline

<a id="opIdpausePipeline"></a>

> Code samples

```http
POST http://localhost:8080/pipelines/{pipelineID}/pause HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`POST /pipelines/{pipelineID}/pause`

*Pause pipeline*

<h3 id="pausepipeline-parameters">Parameters</h3>

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

<h3 id="pausepipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline pause|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## resumePipeline

<a id="opIdresumePipeline"></a>

> Code samples

```http
POST http://localhost:8080/pipelines/{pipelineID}/resume HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`POST /pipelines/{pipelineID}/resume`

*Resume pipeline*

<h3 id="resumepipeline-parameters">Parameters</h3>

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

<h3 id="resumepipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline resumed|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## startPipeline

<a id="opIdstartPipeline"></a>

> Code samples

```http
POST http://localhost:8080/pipelines/{pipelineID}/start HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`POST /pipelines/{pipelineID}/start`

*Start pipeline*

<h3 id="startpipeline-parameters">Parameters</h3>

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

<h3 id="startpipeline-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|202|[Accepted](https://tools.ietf.org/html/rfc7231#section-6.3.3)|Pipeline started|None|
|default|Default|General error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## stopPipeline

<a id="opIdstopPipeline"></a>

> Code samples

```http
POST http://localhost:8080/pipelines/{pipelineID}/stop HTTP/1.1
Host: localhost:8080
Accept: application/json

```

`POST /pipelines/{pipelineID}/stop`

*Stop pipeline*

<h3 id="stoppipeline-parameters">Parameters</h3>

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

<h3 id="stoppipeline-responses">Responses</h3>

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

