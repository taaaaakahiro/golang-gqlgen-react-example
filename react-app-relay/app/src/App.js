import React from 'react';
import './App.css';
import {
  RelayEnvironmentProvider
} from 'react-relay/hooks';
import RelayEnvironment from './RelayEnvironment';
import UsersAndMessages from "./UsersAndMessages";

const { Suspense } = React;

// MEMO) Relay化
// ref https://relay.dev/docs/getting-started/step-by-step-guide/#step-5-fetching-a-query-with-relay

// The above component needs to know how to access the Relay environment, and we
// need to specify a fallback in case it suspends:
// - <RelayEnvironmentProvider> tells child components how to talk to the current
//   Relay Environment instance
// - <Suspense> specifies a fallback in case a child suspends.
function AppRoot(props) {
  return (
      <RelayEnvironmentProvider environment={RelayEnvironment}>
        <Suspense fallback={'Loading...'}>
          <UsersAndMessages />
        </Suspense>
      </RelayEnvironmentProvider>
  );
}

export default AppRoot;