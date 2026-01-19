import '../css/msg_item_time.css'

const CmdHintMsgItem = ({msg}) => {
    // console.log(`msg=${JSON.stringify(msg)}`)

    const content = msg.content.content

    const extractText = () => {
        if (content.subCmd === 1001) {
            return content.inviteContent.inviteHint
        }

        return ""
    }

    const timeEleStyle = () => {
        return {display: msg.timeShow ? 'block' : 'none'}
    }

    return (
        <div className="msg-item-outer">
            <div
                style={timeEleStyle()}
                className="msg-item-time">
                {msg.timeText}
            </div>
            <div style={{
                boxSizing: "border-box",
                marginRight: 18,
                display: "flex",
                marginBottom: 8,
                minHeight: 34,
                justifyContent: "center",
                alignItems: "center"
            }}>
                <div style={{color: "#9E9E9E", fontSize: 12}}>{extractText()}</div>
            </div>

        </div>
    )
}

export default CmdHintMsgItem