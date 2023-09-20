---
title: Enum Documentation
description: JSON schema example demonstrating documentation of enum values using the custom meta:enum keyword. This must be enabled using the enableMetaEnum config option.

---


# Enum Documentation

<p>JSON schema example demonstrating documentation of enum values using the custom meta:enum keyword. This must be enabled using the enableMetaEnum config option.</p>

<table>
<tbody>
<tr><th>$id</th><td>enum-documentation.yml</td></tr>
<tr><th>$schema</th><td>http://json-schema.org/draft-07/schema#</td></tr>
</tbody>
</table>

## Properties

<table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#status">status</a></td><td>String</td></tr></tbody></table>


## Example
```
{
    "status": "Active"
}
```

<hr />



## status


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Status</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The status of something</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Enum</th>
      <td colspan="2"><dl><dt>Active</dt><dd>The thing is currently active and in use</dd><dt>Suspended</dt><dd>The thing is currently suspended and may later become Active or Terminated</dd><dt>Deleted</dt><dd>The thing has been permanently terminated</dd></dl></td>
    </tr>
  </tbody>
</table>










<hr />

## Schema
```
{
    "$id": "enum-documentation.yml",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "Enum Documentation",
    "description": "JSON schema example demonstrating documentation of enum values using the custom meta:enum keyword. This must be enabled using the enableMetaEnum config option.",
    "type": "object",
    "examples": [
        {
            "status": "Active"
        }
    ],
    "properties": {
        "status": {
            "title": "Status",
            "description": "The status of something",
            "type": "string",
            "enum": [
                "Active",
                "Suspended",
                "Terminated"
            ],
            "meta:enum": {
                "Active": "The thing is currently active and in use",
                "Suspended": "The thing is currently suspended and may later become Active or Terminated",
                "Deleted": "The thing has been permanently terminated"
            }
        }
    },
    "additionalProperties": false,
    "required": [
        "status"
    ]
}
```


