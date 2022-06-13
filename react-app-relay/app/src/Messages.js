import React, {useState, useEffect} from 'react';
import fetchGraphQL from "./fetchGraphQL";
import MessageRow from "./MessageRow";

// MEMO) Relay化していない（打ち切り）
// ref
// https://relay.dev/docs/guided-tour/rendering/queries/
// https://relay.dev/docs/guided-tour/updating-data/graphql-mutations/

/**
 *
 * @returns {JSX.Element}
 * @constructor
 * @see https://www.seplus.jp/dokushuzemi/blog/2021/06/quick_start_react_with_graphql.html
 */
function Messages(props) {
    const [messages, setMessages] = useState([]);
    const [newMessage, setNewMessage] = useState("");
    const rows = [];
    const query = `query getMessages($userID: ID!) {
  messages(userID: $userID) {
    id
    message
    user {
      id
      name
    }
  }
}`;

    // MEMO: query return example
    // {
    //   "data": {
    //     "messages": [
    //       {
    //         "id": "2",
    //         "message": "test message id 2",
    //         "user": {
    //           "id": "1"
    //         }
    //       },
    //       {
    //         "id": "1",
    //         "message": "test message id 1",
    //         "user": {
    //           "id": "1"
    //         }
    //       }
    //     ]
    //   }
    // }

    // 1回だけ
    useEffect(() => {
        fetchGraphQL(query, {userID: props.userId})
            .then(response => {
                setMessages(response.data.messages);
            }).catch(error => {
                console.error(error);
            }
        );
    }, []);

    const onInput = (e) => {
        console.log("onInput  message=" + e.target.value);
        setNewMessage(e.target.value);
    }


    const createMessage = () => {
        console.log("addMessage message=" + newMessage);
        if (!newMessage) {
            return;
        }

        const query = `mutation createMessage($input: NewMessage!) {
  createMessage(input: $input) {
    id
    message
    user {
      id
      name
    }
  }
}`;

        // mutation response example
        // {
        //     "data": {
        //     "createMessage": {
        //         "id": "11",
        //             "message": "test insert",
        //             "user": {
        //             "id": "1",
        //                 "name": "Hoge"
        //         }
        //     }
        // }
        // }

        fetchGraphQL(query, {
            input:
                {
                    "userID": props.userId,
                    "message": newMessage
                }
        })
            .then(response => {
                console.log(response.data);
                console.log(messages);

                // newMessage stateを空にする
                setNewMessage("");

                // messages　のリストを再renderさせるために追加してset
                messages.unshift(response.data.createMessage);
                setMessages(messages);
            }).catch(error => {
            console.error(error);
        });
    }

    messages.map(message => {
        rows.push(
            <MessageRow key={message.id} message={message}/>
        )
    });

    // render する JSXを返す
    return (
        <div>
            <table>
                <thead>
                <tr>
                    <th>id</th>
                    <th>message</th>
                    <th>user id</th>
                    <th>user name</th>
                </tr>
                </thead>
                <tbody>
                {rows}
                </tbody>
            </table>
            <input type="text" value={newMessage} onChange={onInput} placeholder="new message"/>
            <button type="button" onClick={createMessage}>ADD</button>
        </div>
    )

}

export default Messages;