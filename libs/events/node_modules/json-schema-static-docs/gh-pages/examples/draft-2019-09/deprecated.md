---
title: Draft 2019-09 - Deprecated Example
description: A schema demonstrating use of deprecated

---
<div class="jssd-deprecated">⚠️ This schema has been marked as deprecated.</div>

# Draft 2019-09 - Deprecated Example

<p>A schema demonstrating use of deprecated</p>

<table>
<tbody>
<tr><th>$id</th><td>https://example.com/deprecated.schema.json</td></tr>
<tr><th>$schema</th><td>https://json-schema.org/draft/2019-09/schema</td></tr>
</tbody>
</table>

## Properties

<table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#firstname">firstName</a></td><td>String</td></tr><tr><td colspan="2"><a href="#lastname">lastName</a></td><td>String</td></tr></tbody></table>


## Example
```
{
    "firstName": "Neil",
    "lastName": "Williams"
}
```

<hr />



## firstName


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">First Name</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">A persons first name</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">No</td>
    </tr>
    
  </tbody>
</table>






## lastName


<table>
  <tbody>
    <tr>
      <th>Deprecated</th>
      <td colspan="2">true</td>
    </tr>
    <tr>
      <th>Title</th>
      <td colspan="2">Last Name</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">A persons last name</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">No</td>
    </tr>
    
  </tbody>
</table>










<hr />

## Schema
```
{
    "$id": "https://example.com/deprecated.schema.json",
    "$schema": "https://json-schema.org/draft/2019-09/schema",
    "title": "Draft 2019-09 - Deprecated Example",
    "description": "A schema demonstrating use of deprecated",
    "examples": [
        {
            "firstName": "Neil",
            "lastName": "Williams"
        }
    ],
    "deprecated": true,
    "type": "object",
    "properties": {
        "firstName": {
            "deprecated": false,
            "type": "string",
            "title": "First Name",
            "description": "A persons first name"
        },
        "lastName": {
            "deprecated": true,
            "type": "string",
            "title": "Last Name",
            "description": "A persons last name"
        }
    }
}
```


