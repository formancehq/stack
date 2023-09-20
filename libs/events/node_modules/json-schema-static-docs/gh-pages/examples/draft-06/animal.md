---
title: Draft 06 - User Example
description: A schema demonstrating some draft 06 features

---


# Draft 06 - User Example

<p>A schema demonstrating some draft 06 features</p>

<table>
<tbody>
<tr><th>$id</th><td>draft-06-animal.yml</td></tr>
<tr><th>$schema</th><td>http://json-schema.org/draft-06/schema#</td></tr>
</tbody>
</table>

## Properties

<table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#animaltype">animalType</a></td><td>String</td></tr><tr><td colspan="2"><a href="#canfly">canFly</a></td><td>Boolean</td></tr></tbody></table>


## Example
```
{
    "animalType": "Bear",
    "canFly": false
}
```
## Example
```
{
    "animalType": "Bat",
    "canFly": true
}
```

<hr />



## animalType


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Animal Type</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    
  </tbody>
</table>






## canFly


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Can Fly</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Boolean</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    
  </tbody>
</table>










<hr />

## Schema
```
{
    "$id": "draft-06-animal.yml",
    "$schema": "http://json-schema.org/draft-06/schema#",
    "title": "Draft 06 - User Example",
    "description": "A schema demonstrating some draft 06 features",
    "examples": [
        {
            "animalType": "Bear",
            "canFly": false
        },
        {
            "animalType": "Bat",
            "canFly": true
        }
    ],
    "required": [
        "animalType",
        "canFly"
    ],
    "type": "object",
    "properties": {
        "animalType": {
            "type": "string",
            "title": "Animal Type"
        },
        "canFly": {
            "type": "boolean",
            "title": "Can Fly"
        }
    }
}
```


