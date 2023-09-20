---
title: Person
description: JSON schema example for a person entity

---


# Person

<p>JSON schema example for a person entity</p>

<table>
<tbody>
<tr><th>$id</th><td>person.yml</td></tr>
<tr><th>$schema</th><td>http://json-schema.org/draft-07/schema#</td></tr>
</tbody>
</table>

## Properties

<table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#name">name</a></td><td>Object (of type <a href="./name.html">Name</a>)</td></tr><tr><td colspan="2"><a href="#dateofbirth">dateOfBirth</a></td><td>String</td></tr><tr><td colspan="2"><a href="#address">address</a></td><td>Object</td></tr><tr><td colspan="2"><a href="#friends">friends</a></td><td>Array [<a href="./name.html">Name</a>]</td></tr></tbody></table>


## Example
```
{
    "name": {
        "title": "Mr",
        "firstName": "Seymour",
        "lastName": "Butts"
    },
    "dateOfBirth": "1980-01-01",
    "address": {
        "houseNumber": 41,
        "street": "Some street",
        "city": "Swansea",
        "timeAtAddress": {
            "years": 1,
            "months": 3
        }
    }
}
```
## Example
```
{
    "name": {
        "title": "Mr",
        "firstName": "Jane",
        "lastName": "Smith"
    },
    "dateOfBirth": "1980-01-01",
    "address": {
        "houseNumber": 310,
        "street": "Any street",
        "city": "London"
    },
    "friends": [
        {
            "title": "Mr",
            "firstName": "Seymour",
            "lastName": "Butts"
        },
        {
            "title": "Mrs",
            "firstName": "Marge",
            "lastName": "Simpson"
        }
    ]
}
```

<hr />



## name

  <p>Defined in <a href="./name.html">./name.html</a></p>

<table>
  <tbody>
    <tr>
      <th>$id</th>
      <td colspan="2">name.yml</td>
    </tr>
    <tr>
      <th>Title</th>
      <td colspan="2">Name</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">JSON schema example for a name entity</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Object (of type <a href="./name.html">Name</a>)</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    
  </tbody>
</table>

### Properties
  <table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#nametitle">title</a></td><td>String</td></tr><tr><td colspan="2"><a href="#namefirstname">firstName</a></td><td>String</td></tr><tr><td colspan="2"><a href="#namelastname">lastName</a></td><td>String</td></tr></tbody></table>





## dateOfBirth


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Date of birth</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The date at which a person was born.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Format</th>
      <td colspan="2">date</td>
    </tr><tr>
      <th>Examples</th>
      <td colspan="2"><li>1992-10-23</li></td>
    </tr>
  </tbody>
</table>






## address


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Address</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The address at which a person lives.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Object</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    
  </tbody>
</table>

### Properties
  <table><thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead><tbody><tr><td colspan="2"><a href="#addresshousenumber">houseNumber</a></td><td>String</td></tr><tr><td colspan="2"><a href="#addressstreet">street</a></td><td>String</td></tr><tr><td colspan="2"><a href="#addresscity">city</a></td><td>String</td></tr><tr><td colspan="2"><a href="#addresstimeataddress">timeAtAddress</a></td><td>Object</td></tr></tbody></table>


### address.houseNumber


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">House Number</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The house number at which an address is located.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Min Length</th>
      <td colspan="2">1</td>
    </tr><tr>
      <th>Max Length</th>
      <td colspan="2">10</td>
    </tr>
  </tbody>
</table>




### address.street


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Street</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The street in which an address is located.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Min Length</th>
      <td colspan="2">1</td>
    </tr><tr>
      <th>Max Length</th>
      <td colspan="2">250</td>
    </tr>
  </tbody>
</table>




### address.city


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">City</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The city in which an address is located.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">String</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">Yes</td>
    </tr>
    <tr>
      <th>Min Length</th>
      <td colspan="2">1</td>
    </tr><tr>
      <th>Max Length</th>
      <td colspan="2">250</td>
    </tr>
  </tbody>
</table>




### address.timeAtAddress


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Time at address</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">How long the person has lived at this address.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Object</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">No</td>
    </tr>
    
  </tbody>
</table>



### address.timeAtAddress.years


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Years</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The number of years lived at this address.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Number</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">No</td>
    </tr>
    <tr>
      <th>Minimum</th>
      <td colspan="2">1</td>
    </tr><tr>
      <th>Minimum</th>
      <td colspan="2">1</td>
    </tr><tr>
      <th>Maximum</th>
      <td colspan="2">100</td>
    </tr>
  </tbody>
</table>




### address.timeAtAddress.months


<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Months</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">The number of months lived at this address.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Integer</td></tr>
    <tr>
      <th>Required</th>
      <td colspan="2">No</td>
    </tr>
    <tr>
      <th>Minimum</th>
      <td colspan="2">1</td>
    </tr><tr>
      <th>Minimum</th>
      <td colspan="2">1</td>
    </tr><tr>
      <th>Maximum</th>
      <td colspan="2">12</td>
    </tr>
  </tbody>
</table>








## friends

  <p>Defined in <a href="./name.html">./name.html</a></p>

<table>
  <tbody>
    <tr>
      <th>Title</th>
      <td colspan="2">Friends</td>
    </tr>
    <tr>
      <th>Description</th>
      <td colspan="2">An array containing the names of a person&#x27;s friends.</td>
    </tr>
    <tr><th>Type</th><td colspan="2">Array [<a href="./name.html">Name</a>]</td></tr>
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
    "$id": "person.yml",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "Person",
    "description": "JSON schema example for a person entity",
    "type": "object",
    "examples": [
        {
            "name": {
                "title": "Mr",
                "firstName": "Seymour",
                "lastName": "Butts"
            },
            "dateOfBirth": "1980-01-01",
            "address": {
                "houseNumber": 41,
                "street": "Some street",
                "city": "Swansea",
                "timeAtAddress": {
                    "years": 1,
                    "months": 3
                }
            }
        },
        {
            "name": {
                "title": "Mr",
                "firstName": "Jane",
                "lastName": "Smith"
            },
            "dateOfBirth": "1980-01-01",
            "address": {
                "houseNumber": 310,
                "street": "Any street",
                "city": "London"
            },
            "friends": [
                {
                    "title": "Mr",
                    "firstName": "Seymour",
                    "lastName": "Butts"
                },
                {
                    "title": "Mrs",
                    "firstName": "Marge",
                    "lastName": "Simpson"
                }
            ]
        }
    ],
    "properties": {
        "name": {
            "$ref": "./name.yml"
        },
        "dateOfBirth": {
            "title": "Date of birth",
            "description": "The date at which a person was born.",
            "type": "string",
            "format": "date",
            "examples": [
                "1992-10-23"
            ]
        },
        "address": {
            "title": "Address",
            "description": "The address at which a person lives.",
            "type": "object",
            "properties": {
                "houseNumber": {
                    "title": "House Number",
                    "description": "The house number at which an address is located.",
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 10
                },
                "street": {
                    "title": "Street",
                    "description": "The street in which an address is located.",
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 250
                },
                "city": {
                    "title": "City",
                    "description": "The city in which an address is located.",
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 250
                },
                "timeAtAddress": {
                    "title": "Time at address",
                    "description": "How long the person has lived at this address.",
                    "type": "object",
                    "properties": {
                        "years": {
                            "title": "Years",
                            "description": "The number of years lived at this address.",
                            "type": "number",
                            "minimum": 1,
                            "maximum": 100
                        },
                        "months": {
                            "title": "Months",
                            "description": "The number of months lived at this address.",
                            "type": "integer",
                            "minimum": 1,
                            "maximum": 12
                        }
                    }
                }
            },
            "required": [
                "houseNumber",
                "street",
                "city"
            ],
            "additionalProperties": false
        },
        "friends": {
            "title": "Friends",
            "description": "An array containing the names of a person's friends.",
            "type": "array",
            "items": {
                "$ref": "./name.yml"
            }
        }
    },
    "additionalProperties": false,
    "required": [
        "name",
        "dateOfBirth",
        "address"
    ]
}
```


