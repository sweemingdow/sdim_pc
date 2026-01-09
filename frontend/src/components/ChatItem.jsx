import TextMsgItem from "./TextMsgItem.jsx";

const ChatItem = ({chat}) => {
    const {item} = chat;
    if (item.msgType === 1) {
        return (<TextMsgItem textChat={chat}/>)
    }
}

export default ChatItem