/**
 *
 * @param props
 * @returns {JSX.Element}
 * @constructor
 * @see https://www.seplus.jp/dokushuzemi/blog/2021/06/quick_start_react_with_graphql.html
 */
function MessageRow(props) {
    return (
        <tr>
            <td>{props.message.id}</td>
            <td>{props.message.message}</td>
            <td>{props.message.user.id}</td>
            <td>{props.message.user.name}</td>
        </tr>
    );
}
export default MessageRow;