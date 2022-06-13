/**
 * @generated SignedSource<<453350febedc0b9af2489993845cec6e>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type UsersAndMessagesRepositoryNameQuery$variables = {||};
export type UsersAndMessagesRepositoryNameQuery$data = {|
  +users: $ReadOnlyArray<{|
    +id: string,
    +name: string,
  |}>,
|};
export type UsersAndMessagesRepositoryNameQuery = {|
  variables: UsersAndMessagesRepositoryNameQuery$variables,
  response: UsersAndMessagesRepositoryNameQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": null,
    "args": null,
    "concreteType": "User",
    "kind": "LinkedField",
    "name": "users",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "id",
        "storageKey": null
      },
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
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "UsersAndMessagesRepositoryNameQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "UsersAndMessagesRepositoryNameQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "ca3bbd560299b7c3ffd50eb4049548e3",
    "id": null,
    "metadata": {},
    "name": "UsersAndMessagesRepositoryNameQuery",
    "operationKind": "query",
    "text": "query UsersAndMessagesRepositoryNameQuery {\n  users {\n    id\n    name\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "9df8c94b649cdc5fdcf0d5a470626359";

module.exports = ((node/*: any*/)/*: Query<
  UsersAndMessagesRepositoryNameQuery$variables,
  UsersAndMessagesRepositoryNameQuery$data,
>*/);
