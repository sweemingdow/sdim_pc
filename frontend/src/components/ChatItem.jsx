import P2pTextMsgItem from "./P2pTextMsgItem.jsx";

const ChatItem = ({chat}) => {
    const {msg, pw, ph} = chat;
    if (msg.msgType === 1) {
        return (<P2pTextMsgItem
            msg={msg}
            pw={pw}
            ph={ph}/>)
    }
}

export default ChatItem