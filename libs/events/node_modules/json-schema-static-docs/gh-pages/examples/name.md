---
title: Name
description: JSON schema example for a name entity

---


# Name

<p>JSON schema example for a name entity</p>

<table>
<tbody>
<tr><th>$id</th><td>name.yml</td></tr>
<tr><th>$schema</th><td>http://json-schema.org/draft-07/schema#</td></tr>
</tbody>
</table>

## Properties

<table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#title">title</a></td><td>String</td></tr><tr><td colspan="2"><a href="#firstname">firstName</a></td><td>String</td></tr><tr><td colspan="2"><a href="#lastname">lastName</a></td><td>String</td></tr></tbody></table>


## Example
```
{
    "title": "Mr",
    "firstName": "Seymour",
    "lastName": "Butts"
}
```

<hr />



## title


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Title</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The title of a name entity</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Default</th>
      <td colspan="2">Mr</td>
    </tr>
    <tr>
      <th>Enum</th>
      <td colspan="2"><ul><li>Mr</li><li>Mrs</li><li>Miss</li></ul></td>
    </tr>
  </tbody>
</table>






## firstName


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">First Name</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The first name of a name entity</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Min Length</th>
      <td colspan="2">3</td>
    </tr><tr>
      <th>Max Length</th>
      <td colspan="2">100</td>
    </tr><tr>
      <th>Examples</th>
      <td colspan="2"><li>Tom</li><li>Dick</li><li>Harry</li></td>
    </tr>
  </tbody>
</table>






## lastName


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Last Name</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The last name of a name entity</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Min Length</th>
      <td colspan="2">3</td>
    </tr><tr>
      <th>Max Length</th>
      <td colspan="2">100</td>
    </tr><tr>
      <th>Examples</th>
      <td colspan="2"><li>Smith</li><li>Jones</li></td>
    </tr>
  </tbody>
</table>










<hr />

## Schema
```
{
    "$id": "name.yml",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "Name",
    "description": "JSON schema example for a name entity",
    "type": "object",
    "examples": [
        {
            "title": "Mr",
            "firstName": "Seymour",
            "lastName": "Butts"
        }
    ],
    "properties": {
        "title": {
            "title": "Title",
            "description": "The title of a name entity",
            "type": "string",
            "default": "Mr",
            "enum": [
                "Mr",
                "Mrs",
                "Miss"
            ]
        },
        "firstName": {
            "title": "First Name",
            "description": "The first name of a name entity",
            "type": "string",
            "minLength": 3,
            "maxLength": 100,
            "examples": [
                "Tom",
                "Dick",
                "Harry"
            ]
        },
        "lastName": {
            "title": "Last Name",
            "description": "The last name of a name entity",
            "type": "string",
            "minLength": 3,
            "maxLength": 100,
            "examples": [
                "Smith",
                "Jones"
            ]
        }
    },
    "additionalProperties": false,
    "required": [
        "title",
        "firstName",
        "lastName"
    ]
}
```


