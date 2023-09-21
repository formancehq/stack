---
title: One-of and Any-of
description: Example schema to demonstrate one of and any of

---
# One-of and Any-of

<p>Example schema to demonstrate one of and any of</p>

<table>
<tbody>
<tr><th>$id</th><td>oneof-and-anyof.yml</td></tr>
<tr><th>$schema</th><td>http://json-schema.org/draft-07/schema#</td></tr>
</tbody>
</table>

## Properties

<table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><th rowspan="2">justOne</th><td rowspan="2">One of:</td><td>Object</td></tr><tr><td>Object</td></tr></tbody></table>


## Example
```
{
    "justOne": {
        "propertyA": "With a string value"
    }
}
```
## Example
```
{
    "justOne": {
        "propertyB": 123,
        "propertyC": 456
    }
}
```

<hr />



## justOne


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Just One</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">Property that demonstrates oneOf</td>
    </tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr><tr><th rowspan="2">Type</th><td rowspan="2">One of:</td><td>Object</td></tr><tr><td>Object</td></tr></tr>
    
  </tbody>
</table>



### justOne.0


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">justOne option 0 with a single property</td>
    </tr>
    
    
  </tbody>
</table>



### justOne.0.propertyA


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Property A</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    
  </tbody>
</table>





### justOne.1


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">justOne option 1 with two properties</td>
    </tr>
    
    
  </tbody>
</table>



### justOne.1.propertyB


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Property B</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Integer</td></tr>
    
  </tbody>
</table>




### justOne.1.propertyC


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Property C</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Integer</td></tr>
    
  </tbody>
</table>











## Schema
```
{
    "$id": "oneof-and-anyof.yml",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "One-of and Any-of",
    "description": "Example schema to demonstrate one of and any of",
    "type": "object",
    "examples": [
        {
            "justOne": {
                "propertyA": "With a string value"
            }
        },
        {
            "justOne": {
                "propertyB": 123,
                "propertyC": 456
            }
        }
    ],
    "properties": {
        "justOne": {
            "title": "Just One",
            "description": "Property that demonstrates oneOf",
            "type": "object",
            "oneOf": [
                {
                    "title": "justOne option 0 with a single property",
                    "properties": {
                        "propertyA": {
                            "type": "string",
                            "title": "Property A"
                        }
                    },
                    "required": [
                        "propertyA"
                    ]
                },
                {
                    "title": "justOne option 1 with two properties",
                    "properties": {
                        "propertyB": {
                            "type": "integer",
                            "title": "Property B"
                        },
                        "propertyC": {
                            "type": "integer",
                            "title": "Property C"
                        }
                    },
                    "required": [
                        "propertyB",
                        "propertyC"
                    ]
                }
            ],
            "isRequired": true
        }
    },
    "additionalProperties": false,
    "required": [
        "justOne"
    ]
}
```


