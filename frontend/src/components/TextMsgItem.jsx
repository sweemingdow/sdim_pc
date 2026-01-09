const TextMsgItem = ({textChat}) => {
    const {item, pw, ph} = textChat
    const {isSelf} = item
    const maxWidth = pw * 0.65
    if (isSelf) {
        return (<div style={{
            alignSelf: "flex-end",
            boxSizing: "border-box",
            marginRight: 18,
            display: "flex",
            flexDirection: "row-reverse",
            marginTop: 8,
            marginBottom: 8,
            minHeight: 34,
            maxWidth: maxWidth,
            justifyContent: "reverse",
        }}>
            <div style={{
                flexShrink: 0,
                width: 34,
                height: 34,
                backgroundColor: "red",
                marginLeft: 10,
            }}>
            </div>

            <div style={{
                width: 0,
                height: 0,
                borderTop: '6px solid transparent',
                borderBottom: '6px solid transparent',
                borderLeft: '6px solid #8ce99a',
                marginTop: 11,
            }}></div>


            <div style={{
                boxSizing: "border-box",
                display: "flex",
                borderRadius: 5,
                backgroundColor: "#8ce99a",
                paddingTop: 8,
                paddingBottom: 8,
                paddingLeft: 10,
                paddingRight: 10,
                minHeight: 34,
                textAlign: "start",
                border: 'none',
                fontSize: 12,
                color: "#161616",
            }}>{item.content.text}</div>
        </div>);
    }

    return (<div style={{
        alignSelf: "flex-start",
        boxSizing: "border-box",
        marginRight: 18,
        display: "flex",
        flexDirection: "row",
        marginTop: 8,
        marginBottom: 8,
        minHeight: 34,
        maxWidth: maxWidth,
    }}>
        <div style={{
            flexShrink: 0,
            width: 34,
            height: 34,
            backgroundColor: "red",
            marginLeft: 18,
        }}>
        </div>

        <div style={{
            width: 0,
            height: 0,
            borderTop: '6px solid transparent',
            borderBottom: '6px solid transparent',
            borderRight: '6px solid white',
            marginLeft: 10,
            marginTop: 11,
        }}></div>


        <div style={{
            boxSizing: "border-box",
            display: "flex",
            alignItems: "center",
            paddingTop: 8,
            paddingBottom: 8,
            paddingLeft: 10,
            paddingRight: 10,
            borderRadius: 5,
            backgroundColor: "white",
            minHeight: 34,
            textAlign: "start",
            border: 'none',
            fontSize: 12,
            color: "#161616",
        }}>{item.content.text}</div>
    </div>);
}

export default TextMsgItem