/**
 * @generated SignedSource<<1f22da2f8d7df7cbd4a3d3d3a527fadf>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type MessagesRepositoryNameQuery$variables = {|
  userID: string,
|};
export type MessagesRepositoryNameQuery$data = {|
  +messages: ?$ReadOnlyArray<?{|
    +id: string,
    +message: string,
    +user: {|
      +id: string,
      +name: string,
    |},
  |}>,
|};
export type MessagesRepositoryNameQuery = {|
  variables: MessagesRepositoryNameQuery$variables,
  response: MessagesRepositoryNameQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "userID"
  }
],
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v2 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "userID",
        "variableName": "userID"
      }
    ],
    "concreteType": "Message",
    "kind": "LinkedField",
    "name": "messages",
    "plural": true,
    "selections": [
      (v1/*: any*/),
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "message",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "User",
        "kind": "LinkedField",
        "name": "user",
        "plural": false,
        "selections": [
          (v1/*: any*/),
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "name",
            "storageKey": null
          }
        ],
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "MessagesRepositoryNameQuery",
    "selections": (v2/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MessagesRepositoryNameQuery",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "962595ec1bedaa552f95bac2c5bdf198",
    "id": null,
    "metadata": {},
    "name": "MessagesRepositoryNameQuery",
    "operationKind": "query",
    "text": "query MessagesRepositoryNameQuery(\n  $userID: ID!\n) {\n  messages(userID: $userID) {\n    id\n    message\n    user {\n      id\n      name\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "04286c797b8350c6b3dba6296e34f35c";

module.exports = ((node/*: any*/)/*: Query<
  MessagesRepositoryNameQuery$variables,
  MessagesRepositoryNameQuery$data,
>*/);
