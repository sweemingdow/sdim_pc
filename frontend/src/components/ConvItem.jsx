import "./ConvItem.css"
import noDisturbIcon from '../assets/images/no_disturb.png'
import {prettyTime} from '../utils/time_format.js'
/*
* type (
	ConvItem struct {
		ConvId       string           `json:"convId,omitempty"`
		ConvType     preinld.ConvType `json:"convType,omitempty"`
		Icon         string           `json:"icon,omitempty"`
		Title        string           `json:"title,omitempty"`
		RelationId   string           `json:"relationId,omitempty"`
		Remark       string           `json:"remark,omitempty"`
		PinTop       bool             `json:"pinTop,omitempty"`
		NoDisturb    bool             `json:"noDisturb,omitempty"`
		MsgSeq       int64            `json:"msgSeq,omitempty"`
		LastMsg      preinld.Msg      `json:"lastMsg,omitempty"`
		BrowseMsgSeq int64            `json:"browseMsgSeq,omitempty"`
		UnreadCount  int64            `json:"unreadCount,omitempty"`
		Cts          int64            `json:"cts,omitempty"`
		Uts          int64            `json:"uts,omitempty"`
	}
)
* */
const ConvItem = ({convItem, onClick}) => {
    const displayLastMsg = () => {
        // 单聊
        if (convItem.convType === 1) {
            const msg = convItem.lastMsg
            if (msg) {
                if (msg.content) {
                    // 文字消息
                    if (msg.content.type === 1) {
                        return msg.content.content.text
                    }
                }

            }
        } else if (convItem.convType === 2) {
            console.log(`convItem:\n${JSON.stringify(convItem)}`);

            const msg = convItem.lastMsg
            if (msg) {
                if (msg.content) {
                    // 文字消息
                    if (msg.content.type === 1000) { // 命令消息
                        if (msg.content.content.subCmd === 1001) {// 邀请进群提示
                            if (msg.content.content.inviteContent) {
                                return msg.content.content.inviteContent.inviteHint
                            }
                        }
                    }
                }
            }
        } else {
        }


        return ""
    }

    const unreadCountStyle = (cnt) => {
        let width
        let display = "flex"
        let text = cnt + ''
        if (cnt <= 0) {
            display = "none"
        } else if (cnt >= 1 && cnt < 10) {
            width = 16
        } else if (cnt >= 10 && cnt < 100) {
            width = 22
        } else {
            width = 28
            text = '99+'
        }

        return {
            style: {
                ...{
                    display,
                    height: 16,
                    position: "absolute",
                    backgroundColor: "#FA5151",
                    right: -8,
                    top: 8,
                    borderRadius: 16,
                    justifyContent: 'center',
                    alignItems: "center",
                    fontSize: 10
                }, width
            }, text,
        }
    }

    const {style, text} = unreadCountStyle(convItem.unreadCount || 0)

    return (<div className={`conv-item ${convItem.selected ? 'selected' : 'unselected'}`}
                 onClick={onClick}
                 style={{
                     width: '100%', height: 64, boxSizing: "border-box", paddingLeft: 12, paddingRight: 12,
                 }}>

        <div id="ci-row" style={{
            width: '100%', height: '100%', display: "flex", flexDirection: "row"
        }}>

            <div style={{
                alignSelf: "center",
                width: 36,
                height: '100%',
                position: "relative",
                display: "flex",
                justifyContent: "center",
                alignItems: "center"
            }}>
                <img id="ci-row-img"
                     src={convItem.icon}
                     style={{
                         width: 36, height: 36, flexShrink: 0,
                     }}>
                </img>

                <div style={style}>
                    {text}
                </div>

            </div>

            <div id="ci-row-r" style={{
                flex: 1,
                boxSizing: "border-box", // backgroundColor:"red",
                width: '100%', // height: '100%',
                height: 40,
                marginTop: 12,
                display: "flex",
                flexDirection: "column",
                marginLeft: 12,
                justifyContent: "space-between",
                minWidth: 0,
            }}>

                <div id="ci-row-r-t" style={{
                    display: "flex",
                    flexDirection: "row",
                    fontSize: 14,
                    justifyContent: "space-between",
                    alignItems: "flex-end",
                    minWidth: 0,
                }}>
                    <div style={{
                        minWidth: 0,
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                        whiteSpace: "nowrap",
                        color: '#333',
                        fontWeight: "bold"
                    }}>
                        {convItem.title}
                    </div>

                    <div style={{
                        flexShrink: 0, marginLeft: 8, fontSize: 11, color: '#999'
                    }}>
                        {prettyTime(convItem.uts)}
                    </div>
                </div>

                <div id="ci-row-r-b" style={{
                    display: "flex",
                    flexDirection: "row",
                    fontSize: 14,
                    justifyContent: "space-between",
                    alignItems: "flex-end",
                    minWidth: 0,
                }}>
                    <div style={{
                        minWidth: 0,
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                        whiteSpace: "nowrap",
                        fontSize: 12,
                        color: '#999'
                    }}>
                        {displayLastMsg()}
                    </div>

                    <img src={noDisturbIcon}
                         style={{
                             display: convItem.noDisturb ? 'inline-block' : 'none',
                             width: 14,
                             height: 14,
                             flexShrink: 0,
                             marginLeft: 8
                         }}/>

                    {/*<div style={{*/}
                    {/*    flexShrink: 0, marginLeft: 8, fontSize: 12, color: '#999'*/}
                    {/*}}>*/}
                    {/*    {convItem.noDisturb ? 'mute' : ''}*/}
                    {/*</div>*/}
                </div>


            </div>
        </div>
    </div>)
}

export default ConvItem