import P2pTextMsgItem from "./P2pTextMsgItem.jsx";
import CmdHintMsgItem from "./CmdHintMsgItem.jsx";

const ChatItem = ({chat}) => {
    const {msg, pw, ph} = chat;
    if (msg.msgType === 1) {
        return (<P2pTextMsgItem
            msg={msg}
            pw={pw}
            ph={ph}/>)
    }else if (msg.msgType === 1000) {
        return (
            <CmdHintMsgItem msg={msg}/>
        )
    }
}

export default ChatItem