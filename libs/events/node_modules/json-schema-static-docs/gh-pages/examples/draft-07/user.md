---
title: Draft 07 - User Example
description: A schema demonstrating some draft 07 features

---


# Draft 07 - User Example

<p>A schema demonstrating some draft 07 features</p>

<table>
<tbody>
<tr><th>$id</th><td>draft-07-user.yml</td></tr>
<tr><th>$schema</th><td>http://json-schema.org/draft-07/schema#</td></tr>
</tbody>
</table>

## Properties

<table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#username">username</a></td><td>String=seymour_butz</td></tr><tr><td colspan="2"><a href="#password">password</a></td><td>String</td></tr></tbody></table>


## Example
```
{
    "username": "seymour_butz",
    "password": "M0zT4v3rn"
}
```

<hr />



## username


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Username</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">This is a description</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Read Only</th>
      <td colspan="2">true</td>
    </tr>
    <tr>
      <th>Const</th>
      <td colspan="2">seymour_butz</td>
    </tr>
  </tbody>
</table>






## password


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Password</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">A write only password property</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">No</td>
    </tr>
    <tr>
      <th>Write Only</th>
      <td colspan="2">true</td>
    </tr>
    
  </tbody>
</table>










<hr />

## Schema
```
{
    "$id": "draft-07-user.yml",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "Draft 07 - User Example",
    "description": "A schema demonstrating some draft 07 features",
    "examples": [
        {
            "username": "seymour_butz",
            "password": "M0zT4v3rn"
        }
    ],
    "required": [
        "username"
    ],
    "type": "object",
    "properties": {
        "username": {
            "type": "string",
            "title": "Username",
            "const": "seymour_butz",
            "readOnly": true,
            "description": "This is a description",
            "$comment": "This is a comment"
        },
        "password": {
            "type": "string",
            "title": "Password",
            "writeOnly": true,
            "description": "A write only password property"
        }
    }
}
```


