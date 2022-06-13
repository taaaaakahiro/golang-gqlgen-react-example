/**
 * @generated SignedSource<<807fc9883b4897f12b13c6f95e259328>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type AppRepositoryNameQuery$variables = {||};
export type AppRepositoryNameQuery$data = {|
  +users: $ReadOnlyArray<{|
    +id: string,
    +name: string,
  |}>,
|};
export type AppRepositoryNameQuery = {|
  variables: AppRepositoryNameQuery$variables,
  response: AppRepositoryNameQuery$data,
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
    "name": "AppRepositoryNameQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "AppRepositoryNameQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "82b21ff4335428c42fb8a5fccee9c535",
    "id": null,
    "metadata": {},
    "name": "AppRepositoryNameQuery",
    "operationKind": "query",
    "text": "query AppRepositoryNameQuery {\n  users {\n    id\n    name\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "00e44202d9245e3bf58c37fa73260a1b";

module.exports = ((node/*: any*/)/*: Query<
  AppRepositoryNameQuery$variables,
  AppRepositoryNameQuery$data,
>*/);
