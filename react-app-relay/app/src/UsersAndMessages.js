import React from 'react';
import Users from "./Users";
import graphql from "babel-plugin-relay/macro";
import {loadQuery, usePreloadedQuery} from "react-relay/hooks";
import RelayEnvironment from "./RelayEnvironment";

// MEMO) relay化
// ref https://relay.dev/docs/getting-started/step-by-step-guide/#step-5-fetching-a-query-with-relay

// Define a query
// query name : UsersAndMessages(filename) + "RepositoryNameQuery"
const RepositoryNameQuery = graphql`
    query UsersAndMessagesRepositoryNameQuery {
        users {
            id
            name
        }
    }`;


// Immediately load the query as our app starts. For a real app, we'd move this
// into our routing configuration, preloading data as we transition to new routes.
const preloadedQuery = loadQuery(RelayEnvironment, RepositoryNameQuery, {
    /* query variables */
});

/**
 *
 * @returns {JSX.Element}
 * @constructor
 * @see https://www.seplus.jp/dokushuzemi/blog/2021/06/quick_start_react_with_graphql.html
 */
function UsersAndMessages() {
    const data = usePreloadedQuery(RepositoryNameQuery, preloadedQuery);

    const userRows = [];
    data.users.forEach(user => {
        userRows.push(
            <Users user={user} key={user.id}/>
        )
    });
    return (
        <div>
        {userRows}
        </div>
    );
}

export default UsersAndMessages;