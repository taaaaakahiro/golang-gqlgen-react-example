import React, {useState, useEffect}  from 'react';
import fetchGraphQL from "./fetchGraphQL";
import UserMessages from "./UserMessages";

/**
 *
 * @returns {JSX.Element}
 * @constructor
 * @see https://www.seplus.jp/dokushuzemi/blog/2021/06/quick_start_react_with_graphql.html
 */
function UsersAndMessages() {
    const [users, setUsers] = useState([]);
    const query = `query getUsers {
  users {
    id
    name
  }
}`;

    // MEMO: query return example
    // {
    //   "data": {
    //     "users": [
    //       {
    //         "id": "2",
    //         "name": "Fuga"
    //       },
    //       {
    //         "id": "1",
    //         "name": "Hoge"
    //       }
    //     ]
    //   }
    // }

    useEffect(() => {
        fetchGraphQL(query, {})
            .then(response => {
                setUsers(response.data.users);
            }).catch(error => {
            console.error(error);
        });
    }, []);
    const userRows = [];
    users.forEach(user => {
        userRows.push(
            <UserMessages user={user} key={user.id}/>
        )
    });

    return (
        <div>
        {userRows}
        </div>
    );
}

export default UsersAndMessages;