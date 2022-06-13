import React from "react";
import Messages from "./Messages";

function Users(props) {
    return (
        <div>
            <h1>Messages User: {props.user.id}</h1>
            <p>id: {props.user.id}</p>
            <p>name: {props.user.name}</p>
            <Messages key={props.user.id} userId={props.user.id} />
        </div>
    );
}
export default Users;