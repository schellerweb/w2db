w2db a Go Backend Example for the [w2ui](https://w2ui.com/) JS Library
==========

## Description
### Demonstrates some backend functions based on the w2ui library.
[w2ui](https://w2ui.com/) is a great JavaScript UI Library to create powerfull web pages. This example demostrate the server backend for w2grid and w2form and uses a SQLite database to CRUD (Create, Read, Update & Delete) records. It only describes the basic functions of the two modules. A detailed documentation of the modules can be found on the developer's website. 

![img001](./doc/img-01.png)

The w2grid module has some important functions such as sorting, searching or editing cells directly. It also supports infinite scrolling, which means records are dynamically loaded as you scroll. This example uses a SQLite database file w2db.db with 200,000 records in a single table. Please note that all data sets are fictitious.

<details><summary>Structure of Table: 'customers'</summary>
<p>
<br/>
<h5>
<table>
<tr><td>ID</td><td>Name</td><td>Type</td><td>Description</td></tr>
<tr><td>1</td><td>recid</td><td>integer</td><td>Autoincrement record id</td></tr>
<tr><td>2</td><td>usr</td><td>nvarchar(64)</td><td>Email address as username</td></tr>
<tr><td>3</td><td>pwd</td><td>nvarchar(32)</td><td>Password</td></tr>
<tr><td>4</td><td>title</td><td>nvarchar(12)</td><td>Title</td></tr>
<tr><td>5</td><td>fname</td><td>nvarchar(32)</td><td>First Name</td></tr>
<tr><td>6</td><td>lname</td><td>nvarchar(32)</td><td>Last Name</td></tr>
<tr><td>7</td><td>company</td><td>nvarchar(48)</td><td>Company Name</td></tr>
<tr><td>8</td><td>street</td><td>nvarchar(48)</td><td>Street address</td></tr>
<tr><td>9</td><td>city</td><td>nvarchar(48)</td><td>City address</td></tr>
<tr><td>10</td><td>state</td><td>nvarchar(32)</td><td>State address</td></tr>
<tr><td>11</td><td>zip</td><td>nvarchar(8)</td><td>ZIP Code</td></tr>
<tr><td>12</td><td>country</td><td>char(2)</td><td>Alpha-2 Country Code</td></tr>
<tr><td>12</td><td>phone</td><td>nvarchar(20)</td><td>Phone number</td></tr>

</table>
</h5>
</p>
</details><br/>

The communication between client and server works in JSON format. Please note that not all parameters are sent with every request. Parameters such as search, sorting, etc. are only transmitted if they are required.

<details><summary>Client/Server comunication</summary>

Client request:
```json
request: {
    "limit": 100,
    "offset": 0,
    "searchLogic": "OR",
    "search": [
        {"field": "fname", "type": "text", "operator": "begins", "value": "sche"},
        {"field": "lname", "type": "text", "operator": "begins", "value": "sche"}
    ],
    "sort": [
        {"field": "usr", "direction": "asc"}
    ]
}
```
If successful, the client expects the following data from the server.
```json
{
    "status": "success",
    "total": 200000,
    "records": [
        {"recid": 1, "usr": "richard.bowen@bolagonline.com", "pwd": "si0yoh7Eey3", "title": "Mr.", ... "phone": "+1-401-996-3972"},
        {"recid": 2, "usr": "april.wang@farmouthub.com", "pwd": "roch0Eich", "title": "Mrs.", ... "phone": "+1-773-569-0935"},
        {"recid": 3, "usr": "luis.bailey@pernabucana.com", "pwd": "Ahl3giegah", "title": "Mr.", ... "phone": "+1-509-657-4729"},
        {... till recid=100}
    ]
}
```
If there is an error, the server sends an error message.
```json
{
    "status": "error",
    "message": "error-message"
}
```
To delete a record (toolbar-button 'Delete'), the following request will be sent to the server.
```json
request: {
    "action": "delete",
    "recid": [565]
}
```
To save changes made with inline editing (toolbar-button 'Save'), the following request will be send to the server. 
```json
request: {
    "action": "save",
    "changes": [
        {"recid": 1, "usr": "richard.bowen@extech.com", "pwd": "si@122gTxMO"},
        {"recid": 2, "pwd": "EiGuudeW!3"},
    ]
}
```
The server response of successful delete or save can be: 
```json
{
    "status": "success"
}
```
and on error:
```json
{
    "status": "error",
    "message": "error-message"
}
```
The w2form module provides another way to edit or add records. This is used for the Add New and Edit toolbar buttons. The communication between client and server is shown below.</br>
Request a record:
```json
{
    "cmd": "get",
    "name": "form",
    "recid": [10]
}
```
Response a record:
```json
{
    "status": "success",
    "record": {
        "recid": 10,
        "usr": "william.nolan@ajscats.com",
        "pwd": "aiv0Aing9ee",
        "title": "Mr.",
        "fname": "William",
        "lname": "Nolan",
        "company": "Deco Refreshments, Inc.",
        "street": "3031 Monroe Avenue",
        "city": "Tampa",
        "state": "Florida",
        "zip": "33610",
        "country": "US",
        "phone": "+1-941-803-1575"
    }
}
```
Request save this changed record.
```json
request: {
    "cmd": "save",
    "recid": 10,
    "name": "form",
    "record":{
        "recid": 10,
        "usr": "william.nolan@decorefreshments.com",
        "pwd": "aiv0Aing9ee",
        "title": "Mr.",
        "fname": "William",
        "lname": "Nolan",
        "company": "Deco Refreshments, Inc.",
        "street": "3031 Monroe Avenue",
        "city": "Tampa",
        "state": "Florida",
        "zip": "33610",
        "country": "US",
        "phone": "+1-941-803-1575"
    }
}
```
If there is an error, the response is.
```json
{
    "status": "error",
    "message": "error-message"
}
```
</details>
</br>

[some screen shots](./doc/screens.md)
<hr/>

## Requirements
* [JavaScript library w2ui](https://w2ui.com/)
* [JavaScript library jQuery](https://jquery.com/)
<hr/>

[Versions history](./doc/version.md)

## License
MIT License

Copyright (c) 2023 Holger Scheller

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Note also the licensing of the used JavaScript librarys [w2ui](https://w2ui.com/) and [jQuery](https://jquery.com/)